package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Provider string

const (
	ProviderVirusTotal Provider = "virustotal"
	ProviderGeoIP      Provider = "geoip"
	ProviderAbuseIPDB  Provider = "abuseipdb"
	ProviderAlienVault Provider = "alienvault"
	ProviderShodan     Provider = "shodan"
)

type Enrichment struct {
	ID          uuid.UUID       `json:"id"`
	IndicatorID uuid.UUID       `json:"indicator_id"`
	Provider    Provider        `json:"provider"`
	Data        json.RawMessage `json:"data"`
	FetchedAt   time.Time       `json:"fetched_at"`
	ExpiresAt   *time.Time      `json:"expires_at,omitempty"`
}

type GeoIPData struct {
	Country     string  `json:"country,omitempty"`
	CountryCode string  `json:"country_code,omitempty"`
	City        string  `json:"city,omitempty"`
	Latitude    float64 `json:"latitude,omitempty"`
	Longitude   float64 `json:"longitude,omitempty"`
	ASN         string  `json:"asn,omitempty"`
	Org         string  `json:"org,omitempty"`
	ISP         string  `json:"isp,omitempty"`
}

type ThreatIntelData struct {
	MaliciousVotes int      `json:"malicious_votes"`
	HarmlessVotes  int      `json:"harmless_votes"`
	Reputation     int      `json:"reputation"`
	Tags           []string `json:"tags,omitempty"`
	MalwareFamily  string   `json:"malware_family,omitempty"`
	LastAnalysis   string   `json:"last_analysis,omitempty"`
}
