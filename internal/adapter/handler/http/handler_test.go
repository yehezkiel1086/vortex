package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yehezkiel1086/vortex/internal/core/domain"
)

type mockIngestionSvc struct {
	ingested *domain.Event
}

func (m *mockIngestionSvc) IngestEvent(ctx context.Context, event *domain.Event) (*domain.Event, error) {
	if event.Source == "" {
		return nil, domain.ErrInvalidInput
	}
	event.ID = uuid.New()
	m.ingested = event
	return event, nil
}

type mockEventRepo struct{}

func (m *mockEventRepo) Create(ctx context.Context, event *domain.Event) (*domain.Event, error) {
	return event, nil
}
func (m *mockEventRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	return &domain.Event{ID: id, Source: "honeypot"}, nil
}
func (m *mockEventRepo) List(ctx context.Context, limit, offset int32) ([]*domain.Event, error) {
	return []*domain.Event{{ID: uuid.New(), Source: "honeypot"}}, nil
}
func (m *mockEventRepo) Count(ctx context.Context) (int64, error) {
	return 1, nil
}

type mockInvestigationSvc struct{}

func (m *mockInvestigationSvc) GetIndicatorDetails(ctx context.Context, indType domain.IndicatorType, value string) (*domain.Indicator, []*domain.Observation, []*domain.Enrichment, []*domain.Relationship, error) {
	return &domain.Indicator{
		ID:    uuid.New(),
		Type:  indType,
		Value: value,
	}, nil, nil, nil, nil
}

func (m *mockInvestigationSvc) GetStats(ctx context.Context) (map[string]interface{}, error) {
	return map[string]interface{}{"total_events": int64(10)}, nil
}

type mockAlertSvc struct{}

func (m *mockAlertSvc) EvaluateAlert(ctx context.Context, indicator *domain.Indicator, event *domain.Event, risk *domain.RiskScore) (*domain.Alert, error) {
	return nil, nil
}
func (m *mockAlertSvc) GetAlert(ctx context.Context, id uuid.UUID) (*domain.Alert, error) {
	return &domain.Alert{ID: id, Severity: domain.SeverityHigh, Status: domain.AlertStatusOpen}, nil
}
func (m *mockAlertSvc) ListAlerts(ctx context.Context, status *domain.AlertStatus, limit, offset int32) ([]*domain.Alert, error) {
	return []*domain.Alert{{ID: uuid.New(), Severity: domain.SeverityHigh, Status: domain.AlertStatusOpen}}, nil
}
func (m *mockAlertSvc) UpdateAlertStatus(ctx context.Context, id uuid.UUID, status domain.AlertStatus) (*domain.Alert, error) {
	return &domain.Alert{ID: id, Status: status}, nil
}

type mockIndRepo struct{}

func (m *mockIndRepo) Upsert(ctx context.Context, ind *domain.Indicator) (*domain.Indicator, error) {
	return ind, nil
}
func (m *mockIndRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Indicator, error) {
	return &domain.Indicator{ID: id}, nil
}
func (m *mockIndRepo) GetByTypeValue(ctx context.Context, indType domain.IndicatorType, value string) (*domain.Indicator, error) {
	return &domain.Indicator{Type: indType, Value: value}, nil
}
func (m *mockIndRepo) List(ctx context.Context, limit, offset int32) ([]*domain.Indicator, error) {
	return []*domain.Indicator{{ID: uuid.New(), Type: domain.IndicatorTypeIP, Value: "185.10.20.30"}}, nil
}
func (m *mockIndRepo) ListByType(ctx context.Context, indType domain.IndicatorType, limit, offset int32) ([]*domain.Indicator, error) {
	return []*domain.Indicator{{ID: uuid.New(), Type: indType, Value: "185.10.20.30"}}, nil
}
func (m *mockIndRepo) UpdateRiskScore(ctx context.Context, id uuid.UUID, riskScore, confidence float64) (*domain.Indicator, error) {
	return &domain.Indicator{ID: id, RiskScore: riskScore}, nil
}
func (m *mockIndRepo) Count(ctx context.Context) (int64, error) {
	return 1, nil
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return NewRouter(Config{
		IngestionService:     &mockIngestionSvc{},
		InvestigationService: &mockInvestigationSvc{},
		AlertService:         &mockAlertSvc{},
		EventRepo:            &mockEventRepo{},
		IndicatorRepo:        &mockIndRepo{},
	})
}

func TestHTTPHandlers(t *testing.T) {
	router := setupTestRouter()

	t.Run("GET /health", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/health", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.Code)
		}
	})

	t.Run("POST /api/v1/events success", func(t *testing.T) {
		payload, _ := json.Marshal(domain.Event{
			Source:   "honeypot",
			SourceIP: "185.10.20.30",
			Severity: domain.SeverityHigh,
		})
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/events", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusCreated {
			t.Errorf("expected 201 Created, got %d: %s", resp.Code, resp.Body.String())
		}
	})

	t.Run("POST /api/v1/events invalid input", func(t *testing.T) {
		payload, _ := json.Marshal(domain.Event{
			Timestamp: time.Now(),
		})
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/events", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusBadRequest {
			t.Errorf("expected 400 Bad Request, got %d", resp.Code)
		}
	})

	t.Run("GET /api/v1/indicators", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/indicators", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.Code)
		}
	})

	t.Run("GET /api/v1/alerts", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/alerts", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.Code)
		}
	})

	t.Run("GET /api/v1/stats", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/stats", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.Code)
		}
	})
}
