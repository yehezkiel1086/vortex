package domain

import (
	"time"

	"github.com/google/uuid"
)

type IndicatorType string

const (
	IndicatorTypeIP     IndicatorType = "ip"
	IndicatorTypeDomain IndicatorType = "domain"
	IndicatorTypeURL    IndicatorType = "url"
	IndicatorTypeSHA256 IndicatorType = "sha256"
	IndicatorTypeSHA1   IndicatorType = "sha1"
	IndicatorTypeMD5    IndicatorType = "md5"
	IndicatorTypeEmail  IndicatorType = "email"
	IndicatorTypeASN    IndicatorType = "asn"
)

type IndicatorStatus string

const (
	IndicatorStatusActive      IndicatorStatus = "active"
	IndicatorStatusExpired     IndicatorStatus = "expired"
	IndicatorStatusWhitelisted IndicatorStatus = "whitelisted"
)

type Indicator struct {
	ID         uuid.UUID       `json:"id"`
	Type       IndicatorType   `json:"type"`
	Value      string          `json:"value"`
	FirstSeen  time.Time       `json:"first_seen"`
	LastSeen   time.Time       `json:"last_seen"`
	RiskScore  float64         `json:"risk_score"`
	Confidence float64         `json:"confidence"`
	Status     IndicatorStatus `json:"status"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}
