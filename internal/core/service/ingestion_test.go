package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/yehezkiel1086/vortex/internal/core/domain"
)

type mockEventRepo struct {
	createdEvent *domain.Event
}

func (m *mockEventRepo) Create(ctx context.Context, event *domain.Event) (*domain.Event, error) {
	m.createdEvent = event
	return event, nil
}

func (m *mockEventRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	return nil, nil
}

func (m *mockEventRepo) List(ctx context.Context, limit, offset int32) ([]*domain.Event, error) {
	return nil, nil
}

func (m *mockEventRepo) Count(ctx context.Context) (int64, error) {
	return 1, nil
}

type mockEventPublisher struct {
	publishedEvent *domain.Event
	shouldFail     bool
}

func (m *mockEventPublisher) PublishEvent(ctx context.Context, event *domain.Event) error {
	if m.shouldFail {
		return errors.New("broker unavailable")
	}
	m.publishedEvent = event
	return nil
}

func TestIngestionService(t *testing.T) {
	ctx := context.Background()

	t.Run("successful ingestion and publish", func(t *testing.T) {
		repo := &mockEventRepo{}
		pub := &mockEventPublisher{}
		svc := NewIngestionService(repo, pub)

		event := &domain.Event{
			Source:     "honeypot",
			SourceIP:   "185.10.20.30",
			Confidence: 90.0, // Percentage, should normalize to 0.9
			Severity:   "HIGH",
		}

		result, err := svc.IngestEvent(ctx, event)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if result.ID == uuid.Nil {
			t.Errorf("expected UUID to be assigned")
		}
		if result.Confidence != 0.9 {
			t.Errorf("expected confidence 0.9, got %f", result.Confidence)
		}
		if result.Severity != domain.SeverityHigh {
			t.Errorf("expected severity high, got %s", result.Severity)
		}
		if pub.publishedEvent == nil || pub.publishedEvent.ID != result.ID {
			t.Errorf("expected event to be published")
		}
	})

	t.Run("missing source validation", func(t *testing.T) {
		repo := &mockEventRepo{}
		svc := NewIngestionService(repo, nil)

		event := &domain.Event{
			Timestamp: time.Now(),
		}

		_, err := svc.IngestEvent(ctx, event)
		if err == nil {
			t.Fatalf("expected error for missing source, got nil")
		}
	})
}
