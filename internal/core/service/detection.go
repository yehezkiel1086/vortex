package service

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/yehezkiel1086/vortex/internal/core/domain"
	"github.com/yehezkiel1086/vortex/internal/core/port"
)

var (
	sqliPattern = regexp.MustCompile(`(?i)(\bUNION\b.*\bSELECT\b|'\s*OR\s*['0-9]|\bSLEEP\(\d+\)|--|\bINFORMATION_SCHEMA\b)`)
	xssPattern  = regexp.MustCompile(`(?i)(<script.*?>|javascript:|onerror\s*=|onload\s*=|alert\()`)
	pathTrav    = regexp.MustCompile(`(?i)(\.\./\.\.|/etc/passwd|win\.ini|\.\.\\\.\.)`)
)

type DetectionService struct {
	obsRepo port.ObservationRepository
	indRepo port.IndicatorRepository
}

func NewDetectionService(obsRepo port.ObservationRepository, indRepo port.IndicatorRepository) port.DetectionService {
	return &DetectionService{
		obsRepo: obsRepo,
		indRepo: indRepo,
	}
}

func (s *DetectionService) Detect(ctx context.Context, event *domain.Event) ([]*domain.Observation, error) {
	if event == nil {
		return nil, domain.ErrInvalidInput
	}

	var observations []*domain.Observation
	eventTime := event.Timestamp
	if eventTime.IsZero() {
		eventTime = time.Now().UTC()
	}

	// Helper to find or resolve indicator ID for the event's source IP
	var indicatorID uuid.UUID
	if event.SourceIP != "" && s.indRepo != nil {
		if ind, err := s.indRepo.GetByTypeValue(ctx, domain.IndicatorTypeIP, event.SourceIP); err == nil && ind != nil {
			indicatorID = ind.ID
		}
	}
	if indicatorID == uuid.Nil {
		indicatorID = uuid.New()
	}

	createObs := func(attackType, techniqueID string, severity domain.Severity, confidence float64) *domain.Observation {
		return &domain.Observation{
			ID:          uuid.New(),
			IndicatorID: indicatorID,
			EventID:     event.ID,
			AttackType:  attackType,
			TechniqueID: techniqueID,
			Timestamp:   eventTime,
			Severity:    severity,
			Confidence:  confidence,
			Source:      event.Source,
			CreatedAt:   eventTime,
		}
	}

	payload := strings.ToLower(string(event.RawPayload)) + " " + strings.ToLower(event.URL)
	attackType := strings.ToLower(event.AttackType)

	// Rule 1: SSH Brute Force
	if attackType == string(domain.AttackTypeSSHBruteForce) || (event.DestinationPort == 22 && strings.Contains(payload, "failed")) {
		observations = append(observations, createObs(string(domain.AttackTypeSSHBruteForce), "T1110", domain.SeverityHigh, 0.90))
	}

	// Rule 2: Port Scanning
	if attackType == string(domain.AttackTypePortScan) || strings.Contains(attackType, "scan") {
		observations = append(observations, createObs(string(domain.AttackTypePortScan), "T1046", domain.SeverityMedium, 0.80))
	}

	// Rule 3: SQL Injection (SQLi)
	if attackType == string(domain.AttackTypeSQLInjection) || sqliPattern.MatchString(payload) {
		observations = append(observations, createObs(string(domain.AttackTypeSQLInjection), "T1190", domain.SeverityHigh, 0.92))
	}

	// Rule 4: Cross-Site Scripting (XSS)
	if attackType == string(domain.AttackTypeXSS) || xssPattern.MatchString(payload) {
		observations = append(observations, createObs(string(domain.AttackTypeXSS), "T1059.007", domain.SeverityMedium, 0.85))
	}

	// Rule 5: Path Traversal
	if attackType == string(domain.AttackTypePathTraversal) || pathTrav.MatchString(payload) {
		observations = append(observations, createObs(string(domain.AttackTypePathTraversal), "T1083", domain.SeverityMedium, 0.85))
	}

	// Rule 6: Malware Download
	if attackType == string(domain.AttackTypeMalwareDownload) || (event.FileHash != "" && event.URL != "") {
		observations = append(observations, createObs(string(domain.AttackTypeMalwareDownload), "T1105", domain.SeverityHigh, 0.95))
	}

	// Fallback if no specific rule matched but event has attack_type
	if len(observations) == 0 && event.AttackType != "" {
		observations = append(observations, createObs(event.AttackType, "T1000", event.Severity, event.Confidence))
	}

	// Persist observations if repo available
	if s.obsRepo != nil {
		for i, obs := range observations {
			saved, err := s.obsRepo.Create(ctx, obs)
			if err == nil {
				observations[i] = saved
			}
		}
	}

	return observations, nil
}
