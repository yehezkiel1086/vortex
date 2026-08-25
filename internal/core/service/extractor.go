package service

import (
	"context"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/yehezkiel1086/vortex/internal/core/domain"
	"github.com/yehezkiel1086/vortex/internal/core/port"
	"github.com/yehezkiel1086/vortex/internal/core/util"
)

var (
	rawIPv4Regex   = regexp.MustCompile(`\b(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\b`)
	rawSHA256Regex = regexp.MustCompile(`\b[a-fA-F0-9]{64}\b`)
	rawMD5Regex    = regexp.MustCompile(`\b[a-fA-F0-9]{32}\b`)
	rawURLRegex    = regexp.MustCompile(`https?://[^\s"'<>]+`)
)

type ExtractorService struct {
	filterPrivateIPs bool
}

func NewExtractorService(filterPrivateIPs bool) port.ExtractorService {
	return &ExtractorService{
		filterPrivateIPs: filterPrivateIPs,
	}
}

func (s *ExtractorService) ExtractIndicators(ctx context.Context, event *domain.Event) ([]*domain.Indicator, error) {
	if event == nil {
		return nil, domain.ErrInvalidInput
	}

	extracted := make(map[string]*domain.Indicator)
	eventTime := event.Timestamp
	if eventTime.IsZero() {
		eventTime = time.Now().UTC()
	}

	addIndicator := func(indType domain.IndicatorType, value string) {
		value = strings.TrimSpace(value)
		if value == "" {
			return
		}

		if indType == domain.IndicatorTypeIP {
			if !util.IsValidIP(value) {
				return
			}
			if s.filterPrivateIPs && util.IsPrivateOrLoopbackIP(value) {
				return
			}
		}

		key := string(indType) + ":" + strings.ToLower(value)
		if _, exists := extracted[key]; !exists {
			extracted[key] = &domain.Indicator{
				ID:         uuid.New(),
				Type:       indType,
				Value:      value,
				FirstSeen:  eventTime,
				LastSeen:   eventTime,
				RiskScore:  0,
				Confidence: event.Confidence,
				Status:     domain.IndicatorStatusActive,
				CreatedAt:  eventTime,
				UpdatedAt:  eventTime,
			}
		}
	}

	// 1. Direct fields
	if event.SourceIP != "" {
		addIndicator(domain.IndicatorTypeIP, event.SourceIP)
	}
	if event.DestinationIP != "" {
		addIndicator(domain.IndicatorTypeIP, event.DestinationIP)
	}
	if event.Domain != "" && util.IsValidDomain(event.Domain) {
		addIndicator(domain.IndicatorTypeDomain, strings.ToLower(event.Domain))
	}
	if event.URL != "" && util.IsValidURL(event.URL) {
		addIndicator(domain.IndicatorTypeURL, event.URL)
		if parsed, err := url.Parse(event.URL); err == nil {
			hostname := parsed.Hostname()
			if util.IsValidIP(hostname) {
				addIndicator(domain.IndicatorTypeIP, hostname)
			} else if util.IsValidDomain(hostname) {
				addIndicator(domain.IndicatorTypeDomain, strings.ToLower(hostname))
			}
		}
	}
	if event.FileHash != "" {
		switch {
		case util.IsValidSHA256(event.FileHash):
			addIndicator(domain.IndicatorTypeSHA256, strings.ToLower(event.FileHash))
		case util.IsValidSHA1(event.FileHash):
			addIndicator(domain.IndicatorTypeSHA1, strings.ToLower(event.FileHash))
		case util.IsValidMD5(event.FileHash):
			addIndicator(domain.IndicatorTypeMD5, strings.ToLower(event.FileHash))
		}
	}

	// 2. Scan RawPayload if available
	if len(event.RawPayload) > 0 {
		payloadStr := string(event.RawPayload)

		// Extract embedded IPs
		for _, ipMatch := range rawIPv4Regex.FindAllString(payloadStr, -1) {
			addIndicator(domain.IndicatorTypeIP, ipMatch)
		}

		// Extract embedded SHA256
		for _, shaMatch := range rawSHA256Regex.FindAllString(payloadStr, -1) {
			addIndicator(domain.IndicatorTypeSHA256, strings.ToLower(shaMatch))
		}

		// Extract embedded MD5
		for _, md5Match := range rawMD5Regex.FindAllString(payloadStr, -1) {
			addIndicator(domain.IndicatorTypeMD5, strings.ToLower(md5Match))
		}

		// Extract embedded URLs
		for _, urlMatch := range rawURLRegex.FindAllString(payloadStr, -1) {
			if util.IsValidURL(urlMatch) {
				addIndicator(domain.IndicatorTypeURL, urlMatch)
			}
		}
	}

	result := make([]*domain.Indicator, 0, len(extracted))
	for _, ind := range extracted {
		result = append(result, ind)
	}

	return result, nil
}
