package service

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/yehezkiel1086/vortex/internal/core/domain"
)

type mockRelRepo struct {
	saved []*domain.Relationship
}

func (m *mockRelRepo) Upsert(ctx context.Context, rel *domain.Relationship) (*domain.Relationship, error) {
	m.saved = append(m.saved, rel)
	return rel, nil
}

func (m *mockRelRepo) GetByIndicatorID(ctx context.Context, indicatorID uuid.UUID) ([]*domain.Relationship, error) {
	return m.saved, nil
}

func TestCorrelationService(t *testing.T) {
	ctx := context.Background()
	relRepo := &mockRelRepo{}
	svc := NewCorrelationService(relRepo, nil)

	ipIndicator := &domain.Indicator{
		ID:        uuid.New(),
		Type:      domain.IndicatorTypeIP,
		Value:     "185.10.20.30",
		FirstSeen: time.Now(),
	}

	event := &domain.Event{
		ID:         uuid.New(),
		SourceIP:   "185.10.20.30",
		Domain:     "evil-c2.com",
		FileHash:   "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		Confidence: 0.95,
	}

	relationships, err := svc.Correlate(ctx, ipIndicator, event)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(relationships) != 2 {
		t.Fatalf("expected 2 relationships (IP->Domain, IP->Hash), got %d", len(relationships))
	}

	foundDomain := false
	foundHash := false

	for _, r := range relationships {
		if r.RelationshipType == domain.RelIPToDomain {
			foundDomain = true
		}
		if r.RelationshipType == domain.RelIPToHash {
			foundHash = true
		}
	}

	if !foundDomain || !foundHash {
		t.Errorf("failed to produce expected relationship types")
	}
}
