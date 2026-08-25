package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/yehezkiel1086/vortex/internal/core/domain"
)

type IngestionService interface {
	IngestEvent(ctx context.Context, event *domain.Event) (*domain.Event, error)
}

type ExtractorService interface {
	ExtractIndicators(ctx context.Context, event *domain.Event) ([]*domain.Indicator, error)
}

type DetectionService interface {
	Detect(ctx context.Context, event *domain.Event) ([]*domain.Observation, error)
}

type EnrichmentService interface {
	EnrichIndicator(ctx context.Context, indicator *domain.Indicator) ([]*domain.Enrichment, error)
}

type CorrelationService interface {
	Correlate(ctx context.Context, indicator *domain.Indicator, event *domain.Event) ([]*domain.Relationship, error)
}

type RiskScoringService interface {
	CalculateRisk(ctx context.Context, indicator *domain.Indicator, observations []*domain.Observation, enrichments []*domain.Enrichment, relationships []*domain.Relationship) (*domain.RiskScore, error)
}

type AlertService interface {
	EvaluateAlert(ctx context.Context, indicator *domain.Indicator, event *domain.Event, risk *domain.RiskScore) (*domain.Alert, error)
	GetAlert(ctx context.Context, id uuid.UUID) (*domain.Alert, error)
	ListAlerts(ctx context.Context, status *domain.AlertStatus, limit, offset int32) ([]*domain.Alert, error)
	UpdateAlertStatus(ctx context.Context, id uuid.UUID, status domain.AlertStatus) (*domain.Alert, error)
}

type InvestigationService interface {
	GetIndicatorDetails(ctx context.Context, indType domain.IndicatorType, value string) (*domain.Indicator, []*domain.Observation, []*domain.Enrichment, []*domain.Relationship, error)
	GetStats(ctx context.Context) (map[string]interface{}, error)
}
