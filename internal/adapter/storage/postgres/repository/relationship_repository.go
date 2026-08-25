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

type RelationshipRepository struct {
	queries *sqlc.Queries
	pool    *pgxpool.Pool
}

func NewRelationshipRepository(pool *pgxpool.Pool) port.RelationshipRepository {
	return &RelationshipRepository{
		queries: sqlc.New(pool),
		pool:    pool,
	}
}

func (r *RelationshipRepository) Upsert(ctx context.Context, relationship *domain.Relationship) (*domain.Relationship, error) {
	if relationship.ID == uuid.Nil {
		relationship.ID = uuid.New()
	}
	now := time.Now().UTC()
	if relationship.FirstSeen.IsZero() {
		relationship.FirstSeen = now
	}
	if relationship.LastSeen.IsZero() {
		relationship.LastSeen = now
	}

	params := sqlc.UpsertRelationshipParams{
		ID:                uuidToPgUUID(relationship.ID),
		SourceIndicatorID: uuidToPgUUID(relationship.SourceIndicatorID),
		TargetIndicatorID: uuidToPgUUID(relationship.TargetIndicatorID),
		RelationshipType:  string(relationship.RelationshipType),
		Confidence:        relationship.Confidence,
		FirstSeen:         timeToTimestamptz(relationship.FirstSeen),
		LastSeen:          timeToTimestamptz(relationship.LastSeen),
	}

	row, err := r.queries.UpsertRelationship(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert relationship: %w", err)
	}

	return toDomainRelationship(&row), nil
}

func (r *RelationshipRepository) GetByIndicatorID(ctx context.Context, indicatorID uuid.UUID) ([]*domain.Relationship, error) {
	rows, err := r.queries.GetRelationshipsByIndicatorID(ctx, uuidToPgUUID(indicatorID))
	if err != nil {
		return nil, fmt.Errorf("failed to get relationships by indicator id: %w", err)
	}

	relationships := make([]*domain.Relationship, len(rows))
	for i, row := range rows {
		relationships[i] = toDomainRelationship(&row)
	}
	return relationships, nil
}
