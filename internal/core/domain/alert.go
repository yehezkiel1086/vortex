package domain

import (
	"time"

	"github.com/google/uuid"
)

type AlertStatus string

const (
	AlertStatusOpen          AlertStatus = "open"
	AlertStatusInvestigating AlertStatus = "investigating"
	AlertStatusResolved      AlertStatus = "resolved"
	AlertStatusFalsePositive AlertStatus = "false_positive"
)

type Alert struct {
	ID          uuid.UUID   `json:"id"`
	IndicatorID uuid.UUID   `json:"indicator_id"`
	EventID     *uuid.UUID  `json:"event_id,omitempty"`
	Severity    Severity    `json:"severity"`
	RiskScore   float64     `json:"risk_score"`
	Confidence  float64     `json:"confidence"`
	Title       string      `json:"title"`
	Description string      `json:"description,omitempty"`
	Status      AlertStatus `json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
	ResolvedAt  *time.Time  `json:"resolved_at,omitempty"`
}
