package service

import (
	"context"
	"encoding/json"
	"math"

	"github.com/yehezkiel1086/vortex/internal/core/domain"
	"github.com/yehezkiel1086/vortex/internal/core/port"
)

type RiskScoringService struct {
	indRepo port.IndicatorRepository
}

func NewRiskScoringService(indRepo port.IndicatorRepository) port.RiskScoringService {
	return &RiskScoringService{
		indRepo: indRepo,
	}
}

func (s *RiskScoringService) CalculateRisk(
	ctx context.Context,
	indicator *domain.Indicator,
	observations []*domain.Observation,
	enrichments []*domain.Enrichment,
	relationships []*domain.Relationship,
) (*domain.RiskScore, error) {
	if indicator == nil {
		return nil, domain.ErrInvalidInput
	}

	breakdown := domain.RiskBreakdown{}

	// 1. Reputation (0 - 30)
	for _, enr := range enrichments {
		if enr.Provider == domain.ProviderVirusTotal && len(enr.Data) > 0 {
			var vt domain.ThreatIntelData
			if err := json.Unmarshal(enr.Data, &vt); err == nil {
				if vt.MaliciousVotes > 0 {
					score := float64(vt.MaliciousVotes) * 0.6
					if score > 30.0 {
						score = 30.0
					}
					if score > breakdown.Reputation {
						breakdown.Reputation = score
					}
				}
				if vt.Reputation < 0 {
					score := math.Abs(float64(vt.Reputation)) * 0.75
					if score > 30.0 {
						score = 30.0
					}
					if score > breakdown.Reputation {
						breakdown.Reputation = score
					}
				}
			}
		}
	}

	// 2. Attack Severity (0 - 25)
	maxSeverityScore := 0.0
	totalConf := 0.0
	for _, obs := range observations {
		score := 0.0
		switch obs.Severity {
		case domain.SeverityCritical:
			score = 25.0
		case domain.SeverityHigh:
			score = 20.0
		case domain.SeverityMedium:
			score = 15.0
		case domain.SeverityLow:
			score = 10.0
		case domain.SeverityInformational:
			score = 5.0
		}
		if score > maxSeverityScore {
			maxSeverityScore = score
		}
		totalConf += obs.Confidence
	}
	breakdown.Severity = maxSeverityScore

	// 3. Frequency (0 - 20)
	freqScore := float64(len(observations)) * 5.0
	if freqScore > 20.0 {
		freqScore = 20.0
	}
	breakdown.Frequency = freqScore

	// 4. Confidence (0 - 15)
	avgConfidence := indicator.Confidence
	if len(observations) > 0 {
		avgConfidence = totalConf / float64(len(observations))
	}
	if avgConfidence > 1.0 {
		avgConfidence = avgConfidence / 100.0
	}
	if avgConfidence == 0.0 {
		avgConfidence = 0.5
	}
	breakdown.Confidence = avgConfidence * 15.0

	// 5. Correlation (0 - 10)
	corrScore := float64(len(relationships)) * 3.5
	if corrScore > 10.0 {
		corrScore = 10.0
	}
	breakdown.Correlation = corrScore

	// Total Risk Calculation
	total := breakdown.Reputation + breakdown.Severity + breakdown.Frequency + breakdown.Confidence + breakdown.Correlation
	if total > 100.0 {
		total = 100.0
	}
	total = math.Round(total*10) / 10 // round to 1 decimal place

	riskScore := &domain.RiskScore{
		TotalScore: total,
		Level:      domain.CalculateRiskLevel(total),
		Confidence: math.Round(avgConfidence*100) / 100,
		Breakdown:  breakdown,
	}

	// Persist updated risk score and confidence on indicator
	if s.indRepo != nil {
		_, _ = s.indRepo.UpdateRiskScore(ctx, indicator.ID, total, avgConfidence)
	}

	return riskScore, nil
}
