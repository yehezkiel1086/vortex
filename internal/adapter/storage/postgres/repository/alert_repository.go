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

type AlertRepository struct {
	queries *sqlc.Queries
	pool    *pgxpool.Pool
}

func NewAlertRepository(pool *pgxpool.Pool) port.AlertRepository {
	return &AlertRepository{
		queries: sqlc.New(pool),
		pool:    pool,
	}
}

func (r *AlertRepository) Create(ctx context.Context, alert *domain.Alert) (*domain.Alert, error) {
	if alert.ID == uuid.Nil {
		alert.ID = uuid.New()
	}
	if alert.CreatedAt.IsZero() {
		alert.CreatedAt = time.Now().UTC()
	}
	if alert.Status == "" {
		alert.Status = domain.AlertStatusOpen
	}
	if alert.Severity == "" {
		alert.Severity = domain.SeverityHigh
	}

	var eventID pgtype.UUID
	if alert.EventID != nil {
		eventID = uuidToPgUUID(*alert.EventID)
	} else {
		eventID = pgtype.UUID{Valid: false}
	}

	params := sqlc.CreateAlertParams{
		ID:          uuidToPgUUID(alert.ID),
		IndicatorID: uuidToPgUUID(alert.IndicatorID),
		EventID:     eventID,
		Severity:    string(alert.Severity),
		RiskScore:   alert.RiskScore,
		Confidence:  alert.Confidence,
		Title:       alert.Title,
		Description: stringToText(alert.Description),
		Status:      string(alert.Status),
	}

	row, err := r.queries.CreateAlert(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create alert: %w", err)
	}

	return toDomainAlert(&row), nil
}

func (r *AlertRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Alert, error) {
	row, err := r.queries.GetAlertByID(ctx, uuidToPgUUID(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get alert by id: %w", err)
	}

	return toDomainAlert(&row), nil
}

func (r *AlertRepository) List(ctx context.Context, limit, offset int32) ([]*domain.Alert, error) {
	rows, err := r.queries.ListAlerts(ctx, sqlc.ListAlertsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list alerts: %w", err)
	}

	alerts := make([]*domain.Alert, len(rows))
	for i, row := range rows {
		alerts[i] = toDomainAlert(&row)
	}
	return alerts, nil
}

func (r *AlertRepository) ListByStatus(ctx context.Context, status domain.AlertStatus, limit, offset int32) ([]*domain.Alert, error) {
	rows, err := r.queries.ListAlertsByStatus(ctx, sqlc.ListAlertsByStatusParams{
		Status: string(status),
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list alerts by status: %w", err)
	}

	alerts := make([]*domain.Alert, len(rows))
	for i, row := range rows {
		alerts[i] = toDomainAlert(&row)
	}
	return alerts, nil
}

func (r *AlertRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.AlertStatus) (*domain.Alert, error) {
	row, err := r.queries.UpdateAlertStatus(ctx, sqlc.UpdateAlertStatusParams{
		ID:     uuidToPgUUID(id),
		Status: string(status),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to update alert status: %w", err)
	}

	return toDomainAlert(&row), nil
}

func (r *AlertRepository) CountByStatus(ctx context.Context) (map[string]int64, error) {
	rows, err := r.queries.CountAlertsByStatus(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count alerts by status: %w", err)
	}

	counts := make(map[string]int64)
	for _, row := range rows {
		counts[row.Status] = row.Count
	}
	return counts, nil
}
