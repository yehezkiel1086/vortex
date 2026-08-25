package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/yehezkiel1086/vortex/internal/core/domain"
	"github.com/yehezkiel1086/vortex/internal/core/port"
)

type CorrelationService struct {
	relRepo port.RelationshipRepository
	indRepo port.IndicatorRepository
}

func NewCorrelationService(relRepo port.RelationshipRepository, indRepo port.IndicatorRepository) port.CorrelationService {
	return &CorrelationService{
		relRepo: relRepo,
		indRepo: indRepo,
	}
}

func (s *CorrelationService) Correlate(ctx context.Context, indicator *domain.Indicator, event *domain.Event) ([]*domain.Relationship, error) {
	if indicator == nil || event == nil {
		return nil, domain.ErrInvalidInput
	}

	var relationships []*domain.Relationship
	now := time.Now().UTC()

	createRel := func(targetID uuid.UUID, relType domain.RelationshipType) *domain.Relationship {
		return &domain.Relationship{
			ID:                uuid.New(),
			SourceIndicatorID: indicator.ID,
			TargetIndicatorID: targetID,
			RelationshipType:  relType,
			Confidence:        event.Confidence,
			FirstSeen:         now,
			LastSeen:          now,
		}
	}

	// 1. Correlate IP with Domain in same event (IP -> DOMAIN)
	if indicator.Type == domain.IndicatorTypeIP && event.Domain != "" {
		targetID := s.resolveIndicatorID(ctx, domain.IndicatorTypeDomain, event.Domain)
		if targetID != uuid.Nil && targetID != indicator.ID {
			relationships = append(relationships, createRel(targetID, domain.RelIPToDomain))
		}
	}

	// 2. Correlate IP with File Hash in same event (IP -> HASH)
	if indicator.Type == domain.IndicatorTypeIP && event.FileHash != "" {
		targetID := s.resolveIndicatorID(ctx, domain.IndicatorTypeSHA256, event.FileHash)
		if targetID != uuid.Nil && targetID != indicator.ID {
			relationships = append(relationships, createRel(targetID, domain.RelIPToHash))
		}
	}

	// 3. Persist relationships
	if s.relRepo != nil {
		for i, rel := range relationships {
			saved, err := s.relRepo.Upsert(ctx, rel)
			if err == nil {
				relationships[i] = saved
			}
		}
	}

	return relationships, nil
}

func (s *CorrelationService) resolveIndicatorID(ctx context.Context, indType domain.IndicatorType, value string) uuid.UUID {
	if s.indRepo == nil {
		return uuid.New()
	}
	ind, err := s.indRepo.GetByTypeValue(ctx, indType, value)
	if err == nil && ind != nil {
		return ind.ID
	}
	// If not found yet, create placeholder indicator
	newInd := &domain.Indicator{
		ID:         uuid.New(),
		Type:       indType,
		Value:      value,
		FirstSeen:  time.Now().UTC(),
		LastSeen:   time.Now().UTC(),
		Status:     domain.IndicatorStatusActive,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}
	if saved, err := s.indRepo.Upsert(ctx, newInd); err == nil {
		return saved.ID
	}
	return newInd.ID
}
