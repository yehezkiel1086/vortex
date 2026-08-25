package port

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/yehezkiel1086/vortex/internal/core/domain"
)

type EventRepository interface {
	Create(ctx context.Context, event *domain.Event) (*domain.Event, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Event, error)
	List(ctx context.Context, limit, offset int32) ([]*domain.Event, error)
	Count(ctx context.Context) (int64, error)
}

type IndicatorRepository interface {
	Upsert(ctx context.Context, indicator *domain.Indicator) (*domain.Indicator, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Indicator, error)
	GetByTypeValue(ctx context.Context, indType domain.IndicatorType, value string) (*domain.Indicator, error)
	List(ctx context.Context, limit, offset int32) ([]*domain.Indicator, error)
	ListByType(ctx context.Context, indType domain.IndicatorType, limit, offset int32) ([]*domain.Indicator, error)
	UpdateRiskScore(ctx context.Context, id uuid.UUID, riskScore, confidence float64) (*domain.Indicator, error)
	Count(ctx context.Context) (int64, error)
}

type ObservationRepository interface {
	Create(ctx context.Context, observation *domain.Observation) (*domain.Observation, error)
	ListByIndicatorID(ctx context.Context, indicatorID uuid.UUID) ([]*domain.Observation, error)
	ListByEventID(ctx context.Context, eventID uuid.UUID) ([]*domain.Observation, error)
}

type EnrichmentRepository interface {
	Upsert(ctx context.Context, enrichment *domain.Enrichment) (*domain.Enrichment, error)
	GetByIndicatorID(ctx context.Context, indicatorID uuid.UUID) ([]*domain.Enrichment, error)
	GetByProvider(ctx context.Context, indicatorID uuid.UUID, provider domain.Provider) (*domain.Enrichment, error)
}

type RelationshipRepository interface {
	Upsert(ctx context.Context, relationship *domain.Relationship) (*domain.Relationship, error)
	GetByIndicatorID(ctx context.Context, indicatorID uuid.UUID) ([]*domain.Relationship, error)
}

type AlertRepository interface {
	Create(ctx context.Context, alert *domain.Alert) (*domain.Alert, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Alert, error)
	List(ctx context.Context, limit, offset int32) ([]*domain.Alert, error)
	ListByStatus(ctx context.Context, status domain.AlertStatus, limit, offset int32) ([]*domain.Alert, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.AlertStatus) (*domain.Alert, error)
	CountByStatus(ctx context.Context) (map[string]int64, error)
}

type CacheRepository interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

type EventPublisher interface {
	PublishEvent(ctx context.Context, event *domain.Event) error
}

type GeoIPClient interface {
	LookupIP(ctx context.Context, ip string) (*domain.GeoIPData, error)
}

type ThreatIntelClient interface {
	LookupHash(ctx context.Context, hash string) (*domain.ThreatIntelData, error)
	LookupIP(ctx context.Context, ip string) (*domain.ThreatIntelData, error)
}
