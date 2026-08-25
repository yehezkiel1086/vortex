package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yehezkiel1086/vortex/internal/adapter/storage/postgres/sqlc"
	"github.com/yehezkiel1086/vortex/internal/core/domain"
	"github.com/yehezkiel1086/vortex/internal/core/port"
)

type EventRepository struct {
	queries *sqlc.Queries
	pool    *pgxpool.Pool
}

func NewEventRepository(pool *pgxpool.Pool) port.EventRepository {
	return &EventRepository{
		queries: sqlc.New(pool),
		pool:    pool,
	}
}

func (r *EventRepository) Create(ctx context.Context, event *domain.Event) (*domain.Event, error) {
	if event.ID == uuid.Nil {
		event.ID = uuid.New()
	}

	params := sqlc.CreateEventParams{
		ID:              uuidToPgUUID(event.ID),
		Timestamp:       timeToTimestamptz(event.Timestamp),
		Source:          event.Source,
		SourceIp:        stringToNetipAddr(event.SourceIP),
		DestinationIp:   stringToNetipAddr(event.DestinationIP),
		SourcePort:      intToInt4(event.SourcePort),
		DestinationPort: intToInt4(event.DestinationPort),
		Protocol:        stringToText(event.Protocol),
		Domain:          stringToText(event.Domain),
		Url:             stringToText(event.URL),
		FileHash:        stringToText(event.FileHash),
		Username:        stringToText(event.Username),
		AttackType:      stringToText(event.AttackType),
		Severity:        string(event.Severity),
		Confidence:      event.Confidence,
		RawPayload:      event.RawPayload,
	}

	row, err := r.queries.CreateEvent(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	return toDomainEvent(&row), nil
}

func (r *EventRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	row, err := r.queries.GetEventByID(ctx, uuidToPgUUID(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get event by id: %w", err)
	}

	return toDomainEvent(&row), nil
}

func (r *EventRepository) List(ctx context.Context, limit, offset int32) ([]*domain.Event, error) {
	rows, err := r.queries.ListEvents(ctx, sqlc.ListEventsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}

	events := make([]*domain.Event, len(rows))
	for i, row := range rows {
		events[i] = toDomainEvent(&row)
	}
	return events, nil
}

func (r *EventRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.queries.CountEvents(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count events: %w", err)
	}
	return count, nil
}
