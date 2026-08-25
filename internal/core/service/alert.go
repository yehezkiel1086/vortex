package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/yehezkiel1086/vortex/internal/core/domain"
	"github.com/yehezkiel1086/vortex/internal/core/port"
)

const AlertThreshold = 70.0

type AlertService struct {
	alertRepo port.AlertRepository
}

func NewAlertService(alertRepo port.AlertRepository) port.AlertService {
	return &AlertService{
		alertRepo: alertRepo,
	}
}

func (s *AlertService) EvaluateAlert(
	ctx context.Context,
	indicator *domain.Indicator,
	event *domain.Event,
	risk *domain.RiskScore,
) (*domain.Alert, error) {
	if indicator == nil || risk == nil {
		return nil, domain.ErrInvalidInput
	}

	if risk.TotalScore < AlertThreshold {
		return nil, nil // below threshold, no alert generated
	}

	severity := domain.SeverityHigh
	if risk.Level == domain.RiskLevelCritical {
		severity = domain.SeverityCritical
	}

	title := fmt.Sprintf("[%s THREAT] %s: %s", strings.ToUpper(string(severity)), strings.ToUpper(string(indicator.Type)), indicator.Value)
	description := fmt.Sprintf(
		"High risk security indicator detected. Risk Score: %.1f/100 (%s), Confidence: %.0f%%. Ingested from source '%s'.",
		risk.TotalScore,
		risk.Level,
		risk.Confidence*100,
		event.Source,
	)

	var eventID *uuid.UUID
	if event != nil && event.ID != uuid.Nil {
		eventID = &event.ID
	}

	alert := &domain.Alert{
		ID:          uuid.New(),
		IndicatorID: indicator.ID,
		EventID:     eventID,
		Severity:    severity,
		RiskScore:   risk.TotalScore,
		Confidence:  risk.Confidence,
		Title:       title,
		Description: description,
		Status:      domain.AlertStatusOpen,
		CreatedAt:   time.Now().UTC(),
	}

	if s.alertRepo != nil {
		return s.alertRepo.Create(ctx, alert)
	}

	return alert, nil
}

func (s *AlertService) GetAlert(ctx context.Context, id uuid.UUID) (*domain.Alert, error) {
	if s.alertRepo == nil {
		return nil, domain.ErrNotFound
	}
	return s.alertRepo.GetByID(ctx, id)
}

func (s *AlertService) ListAlerts(ctx context.Context, status *domain.AlertStatus, limit, offset int32) ([]*domain.Alert, error) {
	if s.alertRepo == nil {
		return nil, nil
	}
	if limit <= 0 {
		limit = 20
	}
	if status != nil && *status != "" {
		return s.alertRepo.ListByStatus(ctx, *status, limit, offset)
	}
	return s.alertRepo.List(ctx, limit, offset)
}

func (s *AlertService) UpdateAlertStatus(ctx context.Context, id uuid.UUID, status domain.AlertStatus) (*domain.Alert, error) {
	if s.alertRepo == nil {
		return nil, domain.ErrNotFound
	}
	return s.alertRepo.UpdateStatus(ctx, id, status)
}
