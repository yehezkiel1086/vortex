package geoip

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
	httpClient *http.Client
	baseURL    string
}

type ipAPIResponse struct {
	Status      string  `json:"status"`
	Message     string  `json:"message"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	City        string  `json:"city"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	ISP         string  `json:"isp"`
	Org         string  `json:"org"`
	AS          string  `json:"as"`
}

func New(baseURL string, timeout time.Duration) port.GeoIPClient {
	if baseURL == "" {
		baseURL = "http://ip-api.com/json"
	}
	if timeout == 0 {
		timeout = 5 * time.Second
	}
	return &Client{
		httpClient: &http.Client{Timeout: timeout},
		baseURL:    baseURL,
	}
}

func (c *Client) LookupIP(ctx context.Context, ip string) (*domain.GeoIPData, error) {
	url := fmt.Sprintf("%s/%s?fields=status,message,country,countryCode,city,lat,lon,isp,org,as", c.baseURL, ip)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create geoip request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("geoip request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("geoip provider returned status %d", resp.StatusCode)
	}

	var res ipAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("failed to decode geoip response: %w", err)
	}

	if res.Status == "fail" {
		return nil, fmt.Errorf("geoip lookup failed: %s", res.Message)
	}

	return &domain.GeoIPData{
		Country:     res.Country,
		CountryCode: res.CountryCode,
		City:        res.City,
		Latitude:    res.Lat,
		Longitude:   res.Lon,
		ASN:         res.AS,
		Org:         res.Org,
		ISP:         res.ISP,
	}, nil
}
