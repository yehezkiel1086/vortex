package virustotal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yehezkiel1086/vortex/internal/core/domain"
	"github.com/yehezkiel1086/vortex/internal/core/port"
)

type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

type vtResponse struct {
	Data struct {
		Attributes struct {
			Reputation        int `json:"reputation"`
			LastAnalysisStats struct {
				Harmless   int `json:"harmless"`
				Malicious  int `json:"malicious"`
				Suspicious int `json:"suspicious"`
				Undetected int `json:"undetected"`
			} `json:"last_analysis_stats"`
			Tags                        []string `json:"tags"`
			PopularThreatClassification struct {
				SuggestedThreatLabel string `json:"suggested_threat_label"`
			} `json:"popular_threat_classification"`
		} `json:"attributes"`
	} `json:"data"`
}

func New(apiKey string, baseURL string, timeout time.Duration) port.ThreatIntelClient {
	if baseURL == "" {
		baseURL = "https://www.virustotal.com/api/v3"
	}
	if timeout == 0 {
		timeout = 5 * time.Second
	}
	return &Client{
		apiKey:     apiKey,
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: timeout},
	}
}

func (c *Client) LookupHash(ctx context.Context, hash string) (*domain.ThreatIntelData, error) {
	if c.apiKey == "" {
		// If no API key configured, return empty intelligence gracefully
		return &domain.ThreatIntelData{
			Reputation: 0,
		}, nil
	}

	url := fmt.Sprintf("%s/files/%s", c.baseURL, hash)
	return c.queryVT(ctx, url)
}

func (c *Client) LookupIP(ctx context.Context, ip string) (*domain.ThreatIntelData, error) {
	if c.apiKey == "" {
		return &domain.ThreatIntelData{
			Reputation: 0,
		}, nil
	}

	url := fmt.Sprintf("%s/ip_addresses/%s", c.baseURL, ip)
	return c.queryVT(ctx, url)
}

func (c *Client) queryVT(ctx context.Context, url string) (*domain.ThreatIntelData, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create virustotal request: %w", err)
	}

	req.Header.Set("x-apikey", c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("virustotal request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return &domain.ThreatIntelData{Reputation: 0}, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("virustotal api returned status %d", resp.StatusCode)
	}

	var res vtResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("failed to decode virustotal response: %w", err)
	}

	attr := res.Data.Attributes
	return &domain.ThreatIntelData{
		MaliciousVotes: attr.LastAnalysisStats.Malicious,
		HarmlessVotes:  attr.LastAnalysisStats.Harmless,
		Reputation:     attr.Reputation,
		Tags:           attr.Tags,
		MalwareFamily:  attr.PopularThreatClassification.SuggestedThreatLabel,
		LastAnalysis:   fmt.Sprintf("malicious: %d, harmless: %d", attr.LastAnalysisStats.Malicious, attr.LastAnalysisStats.Harmless),
	}, nil
}
