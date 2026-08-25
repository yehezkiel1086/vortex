package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yehezkiel1086/vortex/internal/adapter/storage/postgres/sqlc"
	"github.com/yehezkiel1086/vortex/internal/core/domain"
	"github.com/yehezkiel1086/vortex/internal/core/port"
)

type ObservationRepository struct {
	queries *sqlc.Queries
	pool    *pgxpool.Pool
}

func NewObservationRepository(pool *pgxpool.Pool) port.ObservationRepository {
	return &ObservationRepository{
		queries: sqlc.New(pool),
		pool:    pool,
	}
}

func (r *ObservationRepository) Create(ctx context.Context, observation *domain.Observation) (*domain.Observation, error) {
	if observation.ID == uuid.Nil {
		observation.ID = uuid.New()
	}
	if observation.Timestamp.IsZero() {
		observation.Timestamp = time.Now().UTC()
	}
	if observation.Severity == "" {
		observation.Severity = domain.SeverityMedium
	}

	params := sqlc.CreateObservationParams{
		ID:          uuidToPgUUID(observation.ID),
		IndicatorID: uuidToPgUUID(observation.IndicatorID),
		EventID:     uuidToPgUUID(observation.EventID),
		AttackType:  observation.AttackType,
		TechniqueID: stringToText(observation.TechniqueID),
		Timestamp:   timeToTimestamptz(observation.Timestamp),
		Severity:    string(observation.Severity),
		Confidence:  observation.Confidence,
		Source:      stringToText(observation.Source),
	}

	row, err := r.queries.CreateObservation(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create observation: %w", err)
	}

	return toDomainObservation(&row), nil
}

func (r *ObservationRepository) ListByIndicatorID(ctx context.Context, indicatorID uuid.UUID) ([]*domain.Observation, error) {
	rows, err := r.queries.ListObservationsByIndicatorID(ctx, uuidToPgUUID(indicatorID))
	if err != nil {
		return nil, fmt.Errorf("failed to list observations by indicator id: %w", err)
	}

	observations := make([]*domain.Observation, len(rows))
	for i, row := range rows {
		observations[i] = toDomainObservation(&row)
	}
	return observations, nil
}

func (r *ObservationRepository) ListByEventID(ctx context.Context, eventID uuid.UUID) ([]*domain.Observation, error) {
	rows, err := r.queries.ListObservationsByEventID(ctx, uuidToPgUUID(eventID))
	if err != nil {
		return nil, fmt.Errorf("failed to list observations by event id: %w", err)
	}

	observations := make([]*domain.Observation, len(rows))
	for i, row := range rows {
		observations[i] = toDomainObservation(&row)
	}
	return observations, nil
}
