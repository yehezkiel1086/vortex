package domain

import (
	"time"

	"github.com/google/uuid"
)

type Observation struct {
	ID          uuid.UUID `json:"id"`
	IndicatorID uuid.UUID `json:"indicator_id"`
	EventID     uuid.UUID `json:"event_id"`
	AttackType  string    `json:"attack_type"`
	TechniqueID string    `json:"technique_id,omitempty"` // MITRE ATT&CK ID e.g. T1110
	Timestamp   time.Time `json:"timestamp"`
	Severity    Severity  `json:"severity"`
	Confidence  float64   `json:"confidence"`
	Source      string    `json:"source,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}
