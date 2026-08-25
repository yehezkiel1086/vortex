package service

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/yehezkiel1086/vortex/internal/core/domain"
)

func TestExtractorService(t *testing.T) {
	extractor := NewExtractorService(true)
	ctx := context.Background()

	t.Run("extract from direct fields", func(t *testing.T) {
		event := &domain.Event{
			Timestamp:       time.Now(),
			Source:          "honeypot",
			SourceIP:        "185.10.20.30",
			DestinationIP:   "10.0.0.1", // Private IP, should be filtered if filterPrivateIPs is true
			Domain:          "evil-c2.com",
			URL:             "https://malware-drop.com/payload.exe",
			FileHash:        "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			Severity:        domain.SeverityHigh,
			Confidence:      0.9,
		}

		indicators, err := extractor.ExtractIndicators(ctx, event)
		if err != nil {
			t.Fatalf("failed to extract indicators: %v", err)
		}

		typeMap := make(map[domain.IndicatorType]int)
		for _, ind := range indicators {
			typeMap[ind.Type]++
		}

		if typeMap[domain.IndicatorTypeIP] != 1 {
			t.Errorf("expected 1 public IP indicator, got %d", typeMap[domain.IndicatorTypeIP])
		}
		// Expect domains: evil-c2.com and malware-drop.com extracted from URL
		if typeMap[domain.IndicatorTypeDomain] != 2 {
			t.Errorf("expected 2 domain indicators, got %d", typeMap[domain.IndicatorTypeDomain])
		}
		if typeMap[domain.IndicatorTypeURL] != 1 {
			t.Errorf("expected 1 URL indicator, got %d", typeMap[domain.IndicatorTypeURL])
		}
		if typeMap[domain.IndicatorTypeSHA256] != 1 {
			t.Errorf("expected 1 SHA256 indicator, got %d", typeMap[domain.IndicatorTypeSHA256])
		}
	})

	t.Run("extract from raw payload", func(t *testing.T) {
		rawJSON := map[string]interface{}{
			"log": "Attack detected from 198.51.100.45 attempting to download http://bad-site.org/sh.bin with hash 5d41402abc4b2a76b9719d911017c592",
		}
		rawBytes, _ := json.Marshal(rawJSON)

		event := &domain.Event{
			Timestamp:  time.Now(),
			Source:     "suricata",
			RawPayload: rawBytes,
			Severity:   domain.SeverityMedium,
		}

		indicators, err := extractor.ExtractIndicators(ctx, event)
		if err != nil {
			t.Fatalf("failed to extract: %v", err)
		}

		foundIP := false
		foundMD5 := false
		foundURL := false

		for _, ind := range indicators {
			if ind.Type == domain.IndicatorTypeIP && ind.Value == "198.51.100.45" {
				foundIP = true
			}
			if ind.Type == domain.IndicatorTypeMD5 && ind.Value == "5d41402abc4b2a76b9719d911017c592" {
				foundMD5 = true
			}
			if ind.Type == domain.IndicatorTypeURL && ind.Value == "http://bad-site.org/sh.bin" {
				foundURL = true
			}
		}

		if !foundIP || !foundMD5 || !foundURL {
			t.Errorf("failed to extract embedded indicators: IP=%v, MD5=%v, URL=%v", foundIP, foundMD5, foundURL)
		}
	})
}
