package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yehezkiel1086/vortex/internal/adapter/storage/postgres/sqlc"
	"github.com/yehezkiel1086/vortex/internal/core/domain"
	"github.com/yehezkiel1086/vortex/internal/core/port"
)

type EnrichmentRepository struct {
	queries *sqlc.Queries
	pool    *pgxpool.Pool
}

func NewEnrichmentRepository(pool *pgxpool.Pool) port.EnrichmentRepository {
	return &EnrichmentRepository{
		queries: sqlc.New(pool),
		pool:    pool,
	}
}

func (r *EnrichmentRepository) Upsert(ctx context.Context, enrichment *domain.Enrichment) (*domain.Enrichment, error) {
	if enrichment.ID == uuid.Nil {
		enrichment.ID = uuid.New()
	}
	if enrichment.FetchedAt.IsZero() {
		enrichment.FetchedAt = time.Now().UTC()
	}

	var expiresAt pgtype.Timestamptz
	if enrichment.ExpiresAt != nil {
		expiresAt = timeToTimestamptz(*enrichment.ExpiresAt)
	} else {
		expiresAt = pgtype.Timestamptz{Valid: false}
	}

	params := sqlc.UpsertEnrichmentParams{
		ID:          uuidToPgUUID(enrichment.ID),
		IndicatorID: uuidToPgUUID(enrichment.IndicatorID),
		Provider:    string(enrichment.Provider),
		Data:        enrichment.Data,
		FetchedAt:   timeToTimestamptz(enrichment.FetchedAt),
		ExpiresAt:   expiresAt,
	}

	row, err := r.queries.UpsertEnrichment(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert enrichment: %w", err)
	}

	return toDomainEnrichment(&row), nil
}

func (r *EnrichmentRepository) GetByIndicatorID(ctx context.Context, indicatorID uuid.UUID) ([]*domain.Enrichment, error) {
	rows, err := r.queries.GetEnrichmentsByIndicatorID(ctx, uuidToPgUUID(indicatorID))
	if err != nil {
		return nil, fmt.Errorf("failed to get enrichments by indicator id: %w", err)
	}

	enrichments := make([]*domain.Enrichment, len(rows))
	for i, row := range rows {
		enrichments[i] = toDomainEnrichment(&row)
	}
	return enrichments, nil
}

func (r *EnrichmentRepository) GetByProvider(ctx context.Context, indicatorID uuid.UUID, provider domain.Provider) (*domain.Enrichment, error) {
	row, err := r.queries.GetEnrichmentByProvider(ctx, sqlc.GetEnrichmentByProviderParams{
		IndicatorID: uuidToPgUUID(indicatorID),
		Provider:    string(provider),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get enrichment by provider: %w", err)
	}

	return toDomainEnrichment(&row), nil
}
