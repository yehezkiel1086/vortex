package geoip

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGeoIPClient(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"status": "success",
			"country": "United States",
			"countryCode": "US",
			"city": "Ashburn",
			"lat": 39.0438,
			"lon": -77.4874,
			"isp": "Amazon.com, Inc.",
			"org": "AWS EC2",
			"as": "AS16509 Amazon.com, Inc."
		}`))
	}))
	defer mockServer.Close()

	client := New(mockServer.URL, 2*time.Second)
	ctx := context.Background()

	data, err := client.LookupIP(ctx, "54.239.28.85")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if data.Country != "United States" {
		t.Errorf("expected country United States, got %s", data.Country)
	}
	if data.CountryCode != "US" {
		t.Errorf("expected countryCode US, got %s", data.CountryCode)
	}
	if data.City != "Ashburn" {
		t.Errorf("expected city Ashburn, got %s", data.City)
	}
	if data.ASN != "AS16509 Amazon.com, Inc." {
		t.Errorf("expected ASN AS16509, got %s", data.ASN)
	}
}
