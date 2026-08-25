package repository

import (
	"encoding/json"
	"net/netip"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/yehezkiel1086/vortex/internal/adapter/storage/postgres/sqlc"
	"github.com/yehezkiel1086/vortex/internal/core/domain"
)

func TestMapper(t *testing.T) {
	testUUID := uuid.New()
	now := time.Now().Truncate(time.Microsecond)
	testIP, _ := netip.ParseAddr("192.168.1.100")

	t.Run("toDomainEvent", func(t *testing.T) {
		sqlcEvent := &sqlc.Event{
			ID:              pgtype.UUID{Bytes: testUUID, Valid: true},
			Timestamp:       pgtype.Timestamptz{Time: now, Valid: true},
			Source:          "honeypot",
			SourceIp:        &testIP,
			DestinationPort: pgtype.Int4{Int32: 22, Valid: true},
			Protocol:        pgtype.Text{String: "tcp", Valid: true},
			AttackType:      pgtype.Text{String: "ssh_bruteforce", Valid: true},
			Severity:        "high",
			Confidence:      0.95,
			RawPayload:      []byte(`{"test":true}`),
			CreatedAt:       pgtype.Timestamptz{Time: now, Valid: true},
		}

		domainEvent := toDomainEvent(sqlcEvent)
		if domainEvent.ID != testUUID {
			t.Errorf("expected ID %v, got %v", testUUID, domainEvent.ID)
		}
		if domainEvent.SourceIP != "192.168.1.100" {
			t.Errorf("expected SourceIP 192.168.1.100, got %s", domainEvent.SourceIP)
		}
		if domainEvent.DestinationPort != 22 {
			t.Errorf("expected DestinationPort 22, got %d", domainEvent.DestinationPort)
		}
		if domainEvent.Severity != domain.SeverityHigh {
			t.Errorf("expected Severity high, got %s", domainEvent.Severity)
		}
	})

	t.Run("toDomainIndicator", func(t *testing.T) {
		sqlcInd := &sqlc.Indicator{
			ID:         pgtype.UUID{Bytes: testUUID, Valid: true},
			Type:       "ip",
			Value:      "185.10.20.30",
			FirstSeen:  pgtype.Timestamptz{Time: now, Valid: true},
			LastSeen:   pgtype.Timestamptz{Time: now, Valid: true},
			RiskScore:  85.0,
			Confidence: 0.9,
			Status:     "active",
			CreatedAt:  pgtype.Timestamptz{Time: now, Valid: true},
			UpdatedAt:  pgtype.Timestamptz{Time: now, Valid: true},
		}

		ind := toDomainIndicator(sqlcInd)
		if ind.Type != domain.IndicatorTypeIP {
			t.Errorf("expected Type ip, got %s", ind.Type)
		}
		if ind.Value != "185.10.20.30" {
			t.Errorf("expected Value 185.10.20.30, got %s", ind.Value)
		}
		if ind.RiskScore != 85.0 {
			t.Errorf("expected RiskScore 85.0, got %f", ind.RiskScore)
		}
	})

	t.Run("toDomainEnrichment", func(t *testing.T) {
		exp := now.Add(24 * time.Hour)
		sqlcEnr := &sqlc.Enrichment{
			ID:          pgtype.UUID{Bytes: testUUID, Valid: true},
			IndicatorID: pgtype.UUID{Bytes: testUUID, Valid: true},
			Provider:    "geoip",
			Data:        json.RawMessage(`{"country":"US"}`),
			FetchedAt:   pgtype.Timestamptz{Time: now, Valid: true},
			ExpiresAt:   pgtype.Timestamptz{Time: exp, Valid: true},
		}

		enr := toDomainEnrichment(sqlcEnr)
		if enr.Provider != domain.ProviderGeoIP {
			t.Errorf("expected Provider geoip, got %s", enr.Provider)
		}
		if enr.ExpiresAt == nil || !enr.ExpiresAt.Equal(exp) {
			t.Errorf("expected ExpiresAt %v, got %v", exp, enr.ExpiresAt)
		}
	})

	t.Run("toDomainAlert", func(t *testing.T) {
		eventUUID := uuid.New()
		sqlcAlert := &sqlc.Alert{
			ID:          pgtype.UUID{Bytes: testUUID, Valid: true},
			IndicatorID: pgtype.UUID{Bytes: testUUID, Valid: true},
			EventID:     pgtype.UUID{Bytes: eventUUID, Valid: true},
			Severity:    "critical",
			RiskScore:   92.0,
			Confidence:  0.95,
			Title:       "Critical Threat Detected",
			Description: pgtype.Text{String: "Suspicious activity detected", Valid: true},
			Status:      "open",
			CreatedAt:   pgtype.Timestamptz{Time: now, Valid: true},
		}

		alert := toDomainAlert(sqlcAlert)
		if alert.Severity != domain.SeverityCritical {
			t.Errorf("expected Severity critical, got %s", alert.Severity)
		}
		if alert.EventID == nil || *alert.EventID != eventUUID {
			t.Errorf("expected EventID %v, got %v", eventUUID, alert.EventID)
		}
		if alert.Description != "Suspicious activity detected" {
			t.Errorf("expected Description 'Suspicious activity detected', got '%s'", alert.Description)
		}
	})
}
