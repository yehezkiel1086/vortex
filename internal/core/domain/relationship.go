package domain

import (
	"time"

	"github.com/google/uuid"
)

type RelationshipType string

const (
	RelIPToDomain    RelationshipType = "IP->DOMAIN"
	RelIPToHash      RelationshipType = "IP->HASH"
	RelIPToAttack    RelationshipType = "IP->ATTACK"
	RelDomainToIP    RelationshipType = "DOMAIN->IP"
	RelHashToMalware RelationshipType = "HASH->MALWARE"
)

type Relationship struct {
	ID                uuid.UUID        `json:"id"`
	SourceIndicatorID uuid.UUID        `json:"source_indicator_id"`
	TargetIndicatorID uuid.UUID        `json:"target_indicator_id"`
	RelationshipType  RelationshipType `json:"relationship_type"`
	Confidence        float64          `json:"confidence"`
	FirstSeen         time.Time        `json:"first_seen"`
	LastSeen          time.Time        `json:"last_seen"`
}
