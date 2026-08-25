package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yehezkiel1086/vortex/internal/core/domain"
	"github.com/yehezkiel1086/vortex/internal/core/port"
)

const CacheTTL = 24 * time.Hour

type EnrichmentService struct {
	enrichmentRepo port.EnrichmentRepository
	cacheRepo      port.CacheRepository
	geoIPClient    port.GeoIPClient
	tiClient       port.ThreatIntelClient
}

func NewEnrichmentService(
	enrichmentRepo port.EnrichmentRepository,
	cacheRepo port.CacheRepository,
	geoIPClient port.GeoIPClient,
	tiClient port.ThreatIntelClient,
) port.EnrichmentService {
	return &EnrichmentService{
		enrichmentRepo: enrichmentRepo,
		cacheRepo:      cacheRepo,
		geoIPClient:    geoIPClient,
		tiClient:       tiClient,
	}
}

func (s *EnrichmentService) EnrichIndicator(ctx context.Context, indicator *domain.Indicator) ([]*domain.Enrichment, error) {
	if indicator == nil {
		return nil, domain.ErrInvalidInput
	}

	var results []*domain.Enrichment

	switch indicator.Type {
	case domain.IndicatorTypeIP:
		// 1. GeoIP enrichment
		if s.geoIPClient != nil {
			if enr := s.fetchOrGetCached(ctx, indicator, domain.ProviderGeoIP, func() (interface{}, error) {
				return s.geoIPClient.LookupIP(ctx, indicator.Value)
			}); enr != nil {
				results = append(results, enr)
			}
		}

		// 2. Threat Intel (VirusTotal)
		if s.tiClient != nil {
			if enr := s.fetchOrGetCached(ctx, indicator, domain.ProviderVirusTotal, func() (interface{}, error) {
				return s.tiClient.LookupIP(ctx, indicator.Value)
			}); enr != nil {
				results = append(results, enr)
			}
		}

	case domain.IndicatorTypeSHA256, domain.IndicatorTypeSHA1, domain.IndicatorTypeMD5:
		if s.tiClient != nil {
			if enr := s.fetchOrGetCached(ctx, indicator, domain.ProviderVirusTotal, func() (interface{}, error) {
				return s.tiClient.LookupHash(ctx, indicator.Value)
			}); enr != nil {
				results = append(results, enr)
			}
		}
	}

	return results, nil
}

func (s *EnrichmentService) fetchOrGetCached(
	ctx context.Context,
	indicator *domain.Indicator,
	provider domain.Provider,
	fetcher func() (interface{}, error),
) *domain.Enrichment {
	cacheKey := fmt.Sprintf("enrichment:%s:%s:%s", indicator.Type, indicator.Value, provider)

	// 1. Check Redis cache
	if s.cacheRepo != nil {
		if cachedBytes, err := s.cacheRepo.Get(ctx, cacheKey); err == nil && len(cachedBytes) > 0 {
			var cached domain.Enrichment
			if err := json.Unmarshal(cachedBytes, &cached); err == nil {
				return &cached
			}
		}
	}

	// 2. Fetch fresh data from provider
	data, err := fetcher()
	if err != nil {
		return nil
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	now := time.Now().UTC()
	expiresAt := now.Add(CacheTTL)

	enrichment := &domain.Enrichment{
		ID:          uuid.New(),
		IndicatorID: indicator.ID,
		Provider:    provider,
		Data:        json.RawMessage(dataBytes),
		FetchedAt:   now,
		ExpiresAt:   &expiresAt,
	}

	// 3. Persist in database
	if s.enrichmentRepo != nil {
		saved, err := s.enrichmentRepo.Upsert(ctx, enrichment)
		if err == nil {
			enrichment = saved
		}
	}

	// 4. Store in Redis cache
	if s.cacheRepo != nil {
		if encBytes, err := json.Marshal(enrichment); err == nil {
			_ = s.cacheRepo.Set(ctx, cacheKey, encBytes, CacheTTL)
		}
	}

	return enrichment
}
