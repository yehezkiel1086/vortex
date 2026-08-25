package service

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/yehezkiel1086/vortex/internal/core/domain"
)

type mockObsRepo struct {
	created []*domain.Observation
}

func (m *mockObsRepo) Create(ctx context.Context, obs *domain.Observation) (*domain.Observation, error) {
	m.created = append(m.created, obs)
	return obs, nil
}

func (m *mockObsRepo) ListByIndicatorID(ctx context.Context, indicatorID uuid.UUID) ([]*domain.Observation, error) {
	return nil, nil
}

func (m *mockObsRepo) ListByEventID(ctx context.Context, eventID uuid.UUID) ([]*domain.Observation, error) {
	return nil, nil
}

func TestDetectionService(t *testing.T) {
	ctx := context.Background()

	t.Run("detects SSH Brute Force and assigns T1110", func(t *testing.T) {
		repo := &mockObsRepo{}
		svc := NewDetectionService(repo, nil)

		event := &domain.Event{
			ID:              uuid.New(),
			Timestamp:       time.Now(),
			Source:          "honeypot",
			SourceIP:        "185.10.20.30",
			DestinationPort: 22,
			AttackType:      "ssh_bruteforce",
		}

		observations, err := svc.Detect(ctx, event)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(observations) != 1 {
			t.Fatalf("expected 1 observation, got %d", len(observations))
		}

		if observations[0].TechniqueID != "T1110" {
			t.Errorf("expected MITRE T1110, got %s", observations[0].TechniqueID)
		}
		if observations[0].Severity != domain.SeverityHigh {
			t.Errorf("expected SeverityHigh, got %s", observations[0].Severity)
		}
	})

	t.Run("detects SQL Injection via URL payload and assigns T1190", func(t *testing.T) {
		repo := &mockObsRepo{}
		svc := NewDetectionService(repo, nil)

		event := &domain.Event{
			ID:        uuid.New(),
			Timestamp: time.Now(),
			Source:    "waf",
			SourceIP:  "198.51.100.22",
			URL:       "http://target.local/search?q=1'%20UNION%20SELECT%20null,username,password%20FROM%20users--",
		}

		observations, err := svc.Detect(ctx, event)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(observations) != 1 {
			t.Fatalf("expected 1 SQLi observation, got %d", len(observations))
		}

		if observations[0].TechniqueID != "T1190" {
			t.Errorf("expected MITRE T1190, got %s", observations[0].TechniqueID)
		}
	})
}
