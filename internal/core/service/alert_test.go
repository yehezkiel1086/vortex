package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/yehezkiel1086/vortex/internal/core/domain"
)

type mockAlertRepo struct {
	createdAlert *domain.Alert
	alerts       []*domain.Alert
}

func (m *mockAlertRepo) Create(ctx context.Context, alert *domain.Alert) (*domain.Alert, error) {
	m.createdAlert = alert
	m.alerts = append(m.alerts, alert)
	return alert, nil
}
func (m *mockAlertRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Alert, error) {
	if m.createdAlert != nil && m.createdAlert.ID == id {
		return m.createdAlert, nil
	}
	return nil, domain.ErrNotFound
}
func (m *mockAlertRepo) List(ctx context.Context, limit, offset int32) ([]*domain.Alert, error) {
	return m.alerts, nil
}
func (m *mockAlertRepo) ListByStatus(ctx context.Context, status domain.AlertStatus, limit, offset int32) ([]*domain.Alert, error) {
	var res []*domain.Alert
	for _, a := range m.alerts {
		if a.Status == status {
			res = append(res, a)
		}
	}
	return res, nil
}
func (m *mockAlertRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.AlertStatus) (*domain.Alert, error) {
	if m.createdAlert != nil && m.createdAlert.ID == id {
		m.createdAlert.Status = status
		return m.createdAlert, nil
	}
	return nil, domain.ErrNotFound
}
func (m *mockAlertRepo) CountByStatus(ctx context.Context) (map[string]int64, error) {
	return map[string]int64{"open": 1}, nil
}

func TestAlertService(t *testing.T) {
	ctx := context.Background()

	t.Run("generate alert for high risk indicator (score >= 70)", func(t *testing.T) {
		repo := &mockAlertRepo{}
		svc := NewAlertService(repo)

		indicator := &domain.Indicator{
			ID:    uuid.New(),
			Type:  domain.IndicatorTypeIP,
			Value: "185.10.20.30",
		}
		event := &domain.Event{
			ID:     uuid.New(),
			Source: "honeypot",
		}
		risk := &domain.RiskScore{
			TotalScore: 88.5,
			Level:      domain.RiskLevelCritical,
			Confidence: 0.94,
		}

		alert, err := svc.EvaluateAlert(ctx, indicator, event, risk)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if alert == nil {
			t.Fatalf("expected alert to be generated for score 88.5")
		}

		if alert.Severity != domain.SeverityCritical {
			t.Errorf("expected Critical severity, got %s", alert.Severity)
		}
		if alert.Status != domain.AlertStatusOpen {
			t.Errorf("expected open status, got %s", alert.Status)
		}
	})

	t.Run("no alert for low risk (score < 70)", func(t *testing.T) {
		repo := &mockAlertRepo{}
		svc := NewAlertService(repo)

		indicator := &domain.Indicator{ID: uuid.New(), Type: domain.IndicatorTypeIP, Value: "1.1.1.1"}
		risk := &domain.RiskScore{TotalScore: 25.0, Level: domain.RiskLevelInformational}

		alert, err := svc.EvaluateAlert(ctx, indicator, nil, risk)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if alert != nil {
			t.Errorf("expected no alert for low risk score")
		}
	})

	t.Run("update alert status", func(t *testing.T) {
		repo := &mockAlertRepo{}
		svc := NewAlertService(repo)
		alertID := uuid.New()
		repo.createdAlert = &domain.Alert{ID: alertID, Status: domain.AlertStatusOpen}

		updated, err := svc.UpdateAlertStatus(ctx, alertID, domain.AlertStatusResolved)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if updated.Status != domain.AlertStatusResolved {
			t.Errorf("expected status resolved, got %s", updated.Status)
		}
	})
}
