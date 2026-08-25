package repository

import (
	"encoding/json"
	"net/netip"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/yehezkiel1086/vortex/internal/adapter/storage/postgres/sqlc"
	"github.com/yehezkiel1086/vortex/internal/core/domain"
)

func uuidToPgUUID(u uuid.UUID) pgtype.UUID {
	if u == uuid.Nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: u, Valid: true}
}

func pgUUIDToUUID(p pgtype.UUID) uuid.UUID {
	if !p.Valid {
		return uuid.Nil
	}
	return uuid.UUID(p.Bytes)
}

func timeToTimestamptz(t time.Time) pgtype.Timestamptz {
	if t.IsZero() {
		return pgtype.Timestamptz{Valid: false}
	}
	return pgtype.Timestamptz{Time: t, Valid: true}
}

func stringToText(s string) pgtype.Text {
	if s == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: s, Valid: true}
}

func intToInt4(i int) pgtype.Int4 {
	if i == 0 {
		return pgtype.Int4{Valid: false}
	}
	return pgtype.Int4{Int32: int32(i), Valid: true}
}

func stringToNetipAddr(s string) *netip.Addr {
	if s == "" {
		return nil
	}
	addr, err := netip.ParseAddr(s)
	if err != nil {
		return nil
	}
	return &addr
}

func toDomainEvent(row *sqlc.Event) *domain.Event {
	if row == nil {
		return nil
	}
	e := &domain.Event{
		ID:         pgUUIDToUUID(row.ID),
		Timestamp:  row.Timestamp.Time,
		Source:     row.Source,
		Severity:   domain.Severity(row.Severity),
		Confidence: row.Confidence,
		CreatedAt:  row.CreatedAt.Time,
	}

	if row.SourceIp != nil {
		e.SourceIP = row.SourceIp.String()
	}
	if row.DestinationIp != nil {
		e.DestinationIP = row.DestinationIp.String()
	}
	if row.SourcePort.Valid {
		e.SourcePort = int(row.SourcePort.Int32)
	}
	if row.DestinationPort.Valid {
		e.DestinationPort = int(row.DestinationPort.Int32)
	}
	if row.Protocol.Valid {
		e.Protocol = row.Protocol.String
	}
	if row.Domain.Valid {
		e.Domain = row.Domain.String
	}
	if row.Url.Valid {
		e.URL = row.Url.String
	}
	if row.FileHash.Valid {
		e.FileHash = row.FileHash.String
	}
	if row.Username.Valid {
		e.Username = row.Username.String
	}
	if row.AttackType.Valid {
		e.AttackType = row.AttackType.String
	}
	if len(row.RawPayload) > 0 {
		e.RawPayload = json.RawMessage(row.RawPayload)
	}
	return e
}

func toDomainIndicator(row *sqlc.Indicator) *domain.Indicator {
	if row == nil {
		return nil
	}
	return &domain.Indicator{
		ID:         pgUUIDToUUID(row.ID),
		Type:       domain.IndicatorType(row.Type),
		Value:      row.Value,
		FirstSeen:  row.FirstSeen.Time,
		LastSeen:   row.LastSeen.Time,
		RiskScore:  row.RiskScore,
		Confidence: row.Confidence,
		Status:     domain.IndicatorStatus(row.Status),
		CreatedAt:  row.CreatedAt.Time,
		UpdatedAt:  row.UpdatedAt.Time,
	}
}

func toDomainObservation(row *sqlc.Observation) *domain.Observation {
	if row == nil {
		return nil
	}
	obs := &domain.Observation{
		ID:          pgUUIDToUUID(row.ID),
		IndicatorID: pgUUIDToUUID(row.IndicatorID),
		EventID:     pgUUIDToUUID(row.EventID),
		AttackType:  row.AttackType,
		Timestamp:   row.Timestamp.Time,
		Severity:    domain.Severity(row.Severity),
		Confidence:  row.Confidence,
		CreatedAt:   row.CreatedAt.Time,
	}
	if row.TechniqueID.Valid {
		obs.TechniqueID = row.TechniqueID.String
	}
	if row.Source.Valid {
		obs.Source = row.Source.String
	}
	return obs
}

func toDomainEnrichment(row *sqlc.Enrichment) *domain.Enrichment {
	if row == nil {
		return nil
	}
	enr := &domain.Enrichment{
		ID:          pgUUIDToUUID(row.ID),
		IndicatorID: pgUUIDToUUID(row.IndicatorID),
		Provider:    domain.Provider(row.Provider),
		Data:        json.RawMessage(row.Data),
		FetchedAt:   row.FetchedAt.Time,
	}
	if row.ExpiresAt.Valid {
		t := row.ExpiresAt.Time
		enr.ExpiresAt = &t
	}
	return enr
}

func toDomainRelationship(row *sqlc.Relationship) *domain.Relationship {
	if row == nil {
		return nil
	}
	return &domain.Relationship{
		ID:                pgUUIDToUUID(row.ID),
		SourceIndicatorID: pgUUIDToUUID(row.SourceIndicatorID),
		TargetIndicatorID: pgUUIDToUUID(row.TargetIndicatorID),
		RelationshipType:  domain.RelationshipType(row.RelationshipType),
		Confidence:        row.Confidence,
		FirstSeen:         row.FirstSeen.Time,
		LastSeen:          row.LastSeen.Time,
	}
}

func toDomainAlert(row *sqlc.Alert) *domain.Alert {
	if row == nil {
		return nil
	}
	alt := &domain.Alert{
		ID:          pgUUIDToUUID(row.ID),
		IndicatorID: pgUUIDToUUID(row.IndicatorID),
		Severity:    domain.Severity(row.Severity),
		RiskScore:   row.RiskScore,
		Confidence:  row.Confidence,
		Title:       row.Title,
		Status:      domain.AlertStatus(row.Status),
		CreatedAt:   row.CreatedAt.Time,
	}
	if row.EventID.Valid {
		eID := pgUUIDToUUID(row.EventID)
		alt.EventID = &eID
	}
	if row.Description.Valid {
		alt.Description = row.Description.String
	}
	if row.ResolvedAt.Valid {
		t := row.ResolvedAt.Time
		alt.ResolvedAt = &t
	}
	return alt
}
