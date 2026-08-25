package service

import (
	"context"

	"github.com/yehezkiel1086/vortex/internal/core/domain"
	"github.com/yehezkiel1086/vortex/internal/core/port"
)

type InvestigationService struct {
	indRepo   port.IndicatorRepository
	obsRepo   port.ObservationRepository
	enrRepo   port.EnrichmentRepository
	relRepo   port.RelationshipRepository
	eventRepo port.EventRepository
	alertRepo port.AlertRepository
}

func NewInvestigationService(
	indRepo port.IndicatorRepository,
	obsRepo port.ObservationRepository,
	enrRepo port.EnrichmentRepository,
	relRepo port.RelationshipRepository,
	eventRepo port.EventRepository,
	alertRepo port.AlertRepository,
) port.InvestigationService {
	return &InvestigationService{
		indRepo:   indRepo,
		obsRepo:   obsRepo,
		enrRepo:   enrRepo,
		relRepo:   relRepo,
		eventRepo: eventRepo,
		alertRepo: alertRepo,
	}
}

func (s *InvestigationService) GetIndicatorDetails(
	ctx context.Context,
	indType domain.IndicatorType,
	value string,
) (*domain.Indicator, []*domain.Observation, []*domain.Enrichment, []*domain.Relationship, error) {
	if s.indRepo == nil {
		return nil, nil, nil, nil, domain.ErrNotFound
	}

	indicator, err := s.indRepo.GetByTypeValue(ctx, indType, value)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	var observations []*domain.Observation
	if s.obsRepo != nil {
		observations, _ = s.obsRepo.ListByIndicatorID(ctx, indicator.ID)
	}

	var enrichments []*domain.Enrichment
	if s.enrRepo != nil {
		enrichments, _ = s.enrRepo.GetByIndicatorID(ctx, indicator.ID)
	}

	var relationships []*domain.Relationship
	if s.relRepo != nil {
		relationships, _ = s.relRepo.GetByIndicatorID(ctx, indicator.ID)
	}

	return indicator, observations, enrichments, relationships, nil
}

func (s *InvestigationService) GetStats(ctx context.Context) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	if s.eventRepo != nil {
		if eventCount, err := s.eventRepo.Count(ctx); err == nil {
			stats["total_events"] = eventCount
		}
	}

	if s.indRepo != nil {
		if indCount, err := s.indRepo.Count(ctx); err == nil {
			stats["total_indicators"] = indCount
		}
	}

	if s.alertRepo != nil {
		if alertCounts, err := s.alertRepo.CountByStatus(ctx); err == nil {
			stats["alerts_by_status"] = alertCounts
		}
	}

	return stats, nil
}
