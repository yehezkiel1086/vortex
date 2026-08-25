package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yehezkiel1086/vortex/internal/adapter/storage/postgres/sqlc"
	"github.com/yehezkiel1086/vortex/internal/core/domain"
	"github.com/yehezkiel1086/vortex/internal/core/port"
)

type IndicatorRepository struct {
	queries *sqlc.Queries
	pool    *pgxpool.Pool
}

func NewIndicatorRepository(pool *pgxpool.Pool) port.IndicatorRepository {
	return &IndicatorRepository{
		queries: sqlc.New(pool),
		pool:    pool,
	}
}

func (r *IndicatorRepository) Upsert(ctx context.Context, indicator *domain.Indicator) (*domain.Indicator, error) {
	if indicator.ID == uuid.Nil {
		indicator.ID = uuid.New()
	}
	now := time.Now().UTC()
	if indicator.FirstSeen.IsZero() {
		indicator.FirstSeen = now
	}
	if indicator.LastSeen.IsZero() {
		indicator.LastSeen = now
	}
	if indicator.Status == "" {
		indicator.Status = domain.IndicatorStatusActive
	}

	params := sqlc.UpsertIndicatorParams{
		ID:         uuidToPgUUID(indicator.ID),
		Type:       string(indicator.Type),
		Value:      indicator.Value,
		FirstSeen:  timeToTimestamptz(indicator.FirstSeen),
		LastSeen:   timeToTimestamptz(indicator.LastSeen),
		RiskScore:  indicator.RiskScore,
		Confidence: indicator.Confidence,
		Status:     string(indicator.Status),
	}

	row, err := r.queries.UpsertIndicator(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert indicator: %w", err)
	}

	return toDomainIndicator(&row), nil
}

func (r *IndicatorRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Indicator, error) {
	row, err := r.queries.GetIndicatorByID(ctx, uuidToPgUUID(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get indicator by id: %w", err)
	}

	return toDomainIndicator(&row), nil
}

func (r *IndicatorRepository) GetByTypeValue(ctx context.Context, indType domain.IndicatorType, value string) (*domain.Indicator, error) {
	row, err := r.queries.GetIndicatorByTypeValue(ctx, sqlc.GetIndicatorByTypeValueParams{
		Type:  string(indType),
		Value: value,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get indicator by type and value: %w", err)
	}

	return toDomainIndicator(&row), nil
}

func (r *IndicatorRepository) List(ctx context.Context, limit, offset int32) ([]*domain.Indicator, error) {
	rows, err := r.queries.ListIndicators(ctx, sqlc.ListIndicatorsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list indicators: %w", err)
	}

	indicators := make([]*domain.Indicator, len(rows))
	for i, row := range rows {
		indicators[i] = toDomainIndicator(&row)
	}
	return indicators, nil
}

func (r *IndicatorRepository) ListByType(ctx context.Context, indType domain.IndicatorType, limit, offset int32) ([]*domain.Indicator, error) {
	rows, err := r.queries.ListIndicatorsByType(ctx, sqlc.ListIndicatorsByTypeParams{
		Type:   string(indType),
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list indicators by type: %w", err)
	}

	indicators := make([]*domain.Indicator, len(rows))
	for i, row := range rows {
		indicators[i] = toDomainIndicator(&row)
	}
	return indicators, nil
}

func (r *IndicatorRepository) UpdateRiskScore(ctx context.Context, id uuid.UUID, riskScore, confidence float64) (*domain.Indicator, error) {
	row, err := r.queries.UpdateIndicatorRiskScore(ctx, sqlc.UpdateIndicatorRiskScoreParams{
		ID:         uuidToPgUUID(id),
		RiskScore:  riskScore,
		Confidence: confidence,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to update indicator risk score: %w", err)
	}

	return toDomainIndicator(&row), nil
}

func (r *IndicatorRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.queries.CountIndicators(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count indicators: %w", err)
	}
	return count, nil
}
