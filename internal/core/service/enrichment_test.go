package service

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/yehezkiel1086/vortex/internal/core/domain"
)

type mockCacheRepo struct {
	store map[string][]byte
}

func (m *mockCacheRepo) Get(ctx context.Context, key string) ([]byte, error) {
	val, exists := m.store[key]
	if !exists {
		return nil, domain.ErrNotFound
	}
	return val, nil
}

func (m *mockCacheRepo) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	m.store[key] = value
	return nil
}

func (m *mockCacheRepo) Delete(ctx context.Context, key string) error {
	delete(m.store, key)
	return nil
}

type mockEnrichmentRepo struct {
	saved []*domain.Enrichment
}

func (m *mockEnrichmentRepo) Upsert(ctx context.Context, enrichment *domain.Enrichment) (*domain.Enrichment, error) {
	m.saved = append(m.saved, enrichment)
	return enrichment, nil
}

func (m *mockEnrichmentRepo) GetByIndicatorID(ctx context.Context, indicatorID uuid.UUID) ([]*domain.Enrichment, error) {
	return m.saved, nil
}

func (m *mockEnrichmentRepo) GetByProvider(ctx context.Context, indicatorID uuid.UUID, provider domain.Provider) (*domain.Enrichment, error) {
	for _, e := range m.saved {
		if e.IndicatorID == indicatorID && e.Provider == provider {
			return e, nil
		}
	}
	return nil, domain.ErrNotFound
}

type mockGeoIPClient struct {
	called bool
}

func (m *mockGeoIPClient) LookupIP(ctx context.Context, ip string) (*domain.GeoIPData, error) {
	m.called = true
	return &domain.GeoIPData{
		Country:     "United States",
		CountryCode: "US",
		City:        "Chicago",
		ASN:         "AS15169",
	}, nil
}

type mockTIClient struct {
	called bool
}

func (m *mockTIClient) LookupHash(ctx context.Context, hash string) (*domain.ThreatIntelData, error) {
	m.called = true
	return &domain.ThreatIntelData{
		MaliciousVotes: 40,
		Reputation:     -20,
		MalwareFamily:  "agent_tesla",
	}, nil
}

func (m *mockTIClient) LookupIP(ctx context.Context, ip string) (*domain.ThreatIntelData, error) {
	m.called = true
	return &domain.ThreatIntelData{
		MaliciousVotes: 10,
		Reputation:     -5,
	}, nil
}

func TestEnrichmentService(t *testing.T) {
	ctx := context.Background()

	t.Run("fresh enrichment saves to DB and cache", func(t *testing.T) {
		cache := &mockCacheRepo{store: make(map[string][]byte)}
		repo := &mockEnrichmentRepo{}
		geo := &mockGeoIPClient{}
		ti := &mockTIClient{}

		svc := NewEnrichmentService(repo, cache, geo, ti)

		indicator := &domain.Indicator{
			ID:    uuid.New(),
			Type:  domain.IndicatorTypeIP,
			Value: "185.10.20.30",
		}

		enrichments, err := svc.EnrichIndicator(ctx, indicator)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(enrichments) != 2 {
			t.Fatalf("expected 2 enrichments (geoip + virustotal), got %d", len(enrichments))
		}

		if !geo.called || !ti.called {
			t.Errorf("expected external clients to be invoked")
		}

		// Verify cache was populated
		cacheKey := "enrichment:ip:185.10.20.30:geoip"
		if _, exists := cache.store[cacheKey]; !exists {
			t.Errorf("expected cache key %s to be set", cacheKey)
		}
	})

	t.Run("cached enrichment avoids external calls", func(t *testing.T) {
		cachedEnrichment := &domain.Enrichment{
			ID:          uuid.New(),
			IndicatorID: uuid.New(),
			Provider:    domain.ProviderGeoIP,
			Data:        json.RawMessage(`{"country":"Germany"}`),
			FetchedAt:   time.Now(),
		}
		cachedBytes, _ := json.Marshal(cachedEnrichment)

		cacheKey := "enrichment:ip:8.8.8.8:geoip"
		cache := &mockCacheRepo{store: map[string][]byte{cacheKey: cachedBytes}}
		repo := &mockEnrichmentRepo{}
		geo := &mockGeoIPClient{} // should not be called for geoip

		svc := NewEnrichmentService(repo, cache, geo, nil)

		indicator := &domain.Indicator{
			ID:    uuid.New(),
			Type:  domain.IndicatorTypeIP,
			Value: "8.8.8.8",
		}

		enrichments, err := svc.EnrichIndicator(ctx, indicator)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(enrichments) != 1 {
			t.Fatalf("expected 1 cached enrichment, got %d", len(enrichments))
		}

		if geo.called {
			t.Errorf("expected geoip client NOT to be called when cached")
		}
	})
}
