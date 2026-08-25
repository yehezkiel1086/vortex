package service

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/yehezkiel1086/vortex/internal/core/domain"
)

type mockIndRepo struct {
	updatedScore float64
}

func (m *mockIndRepo) Upsert(ctx context.Context, indicator *domain.Indicator) (*domain.Indicator, error) {
	return indicator, nil
}
func (m *mockIndRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Indicator, error) {
	return nil, nil
}
func (m *mockIndRepo) GetByTypeValue(ctx context.Context, indType domain.IndicatorType, value string) (*domain.Indicator, error) {
	return nil, nil
}
func (m *mockIndRepo) List(ctx context.Context, limit, offset int32) ([]*domain.Indicator, error) {
	return nil, nil
}
func (m *mockIndRepo) ListByType(ctx context.Context, indType domain.IndicatorType, limit, offset int32) ([]*domain.Indicator, error) {
	return nil, nil
}
func (m *mockIndRepo) UpdateRiskScore(ctx context.Context, id uuid.UUID, riskScore, confidence float64) (*domain.Indicator, error) {
	m.updatedScore = riskScore
	return &domain.Indicator{ID: id, RiskScore: riskScore, Confidence: confidence}, nil
}
func (m *mockIndRepo) Count(ctx context.Context) (int64, error) {
	return 1, nil
}

func TestRiskScoringService(t *testing.T) {
	ctx := context.Background()
	indRepo := &mockIndRepo{}
	svc := NewRiskScoringService(indRepo)

	indicator := &domain.Indicator{
		ID:         uuid.New(),
		Type:       domain.IndicatorTypeIP,
		Value:      "185.10.20.30",
		Confidence: 0.95,
	}

	vtData, _ := json.Marshal(domain.ThreatIntelData{
		MaliciousVotes: 50, // Reputation max 30
		Reputation:     -30,
	})

	enrichments := []*domain.Enrichment{
		{
			ID:          uuid.New(),
			IndicatorID: indicator.ID,
			Provider:    domain.ProviderVirusTotal,
			Data:        vtData,
		},
	}

	observations := []*domain.Observation{
		{
			ID:          uuid.New(),
			AttackType:  "ssh_bruteforce",
			Severity:    domain.SeverityHigh, // 20
			Confidence:  0.9,
			Timestamp:   time.Now(),
		},
		{
			ID:          uuid.New(),
			AttackType:  "port_scan",
			Severity:    domain.SeverityMedium,
			Confidence:  0.9,
			Timestamp:   time.Now(),
		},
		{
			ID:          uuid.New(),
			AttackType:  "malware_download",
			Severity:    domain.SeverityHigh,
			Confidence:  0.95,
			Timestamp:   time.Now(),
		},
	}

	relationships := []*domain.Relationship{
		{ID: uuid.New(), RelationshipType: domain.RelIPToDomain},
		{ID: uuid.New(), RelationshipType: domain.RelIPToHash},
	}

	riskScore, err := svc.CalculateRisk(ctx, indicator, observations, enrichments, relationships)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if riskScore.Level != domain.RiskLevelCritical && riskScore.Level != domain.RiskLevelHigh {
		t.Errorf("expected high/critical risk level, got %s", riskScore.Level)
	}

	if riskScore.TotalScore < 70.0 {
		t.Errorf("expected high risk total score >= 70, got %f", riskScore.TotalScore)
	}

	if indRepo.updatedScore != riskScore.TotalScore {
		t.Errorf("expected indicator repo to be updated with %f, got %f", riskScore.TotalScore, indRepo.updatedScore)
	}
}
