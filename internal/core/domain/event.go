package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Severity string

const (
	SeverityInformational Severity = "informational"
	SeverityLow           Severity = "low"
	SeverityMedium        Severity = "medium"
	SeverityHigh          Severity = "high"
	SeverityCritical      Severity = "critical"
)

type AttackType string

const (
	AttackTypeSSHBruteForce   AttackType = "ssh_bruteforce"
	AttackTypePortScan        AttackType = "port_scan"
	AttackTypeCredentialStuff AttackType = "credential_stuffing"
	AttackTypeSQLInjection    AttackType = "sqli"
	AttackTypeXSS             AttackType = "xss"
	AttackTypePathTraversal   AttackType = "path_traversal"
	AttackTypeMalwareDownload AttackType = "malware_download"
	AttackTypeSuspiciousCmd   AttackType = "suspicious_command"
	AttackTypeUnknown         AttackType = "unknown"
)

type Event struct {
	ID              uuid.UUID       `json:"id"`
	Timestamp       time.Time       `json:"timestamp"`
	Source          string          `json:"source"`
	SourceIP        string          `json:"source_ip,omitempty"`
	DestinationIP   string          `json:"destination_ip,omitempty"`
	SourcePort      int             `json:"source_port,omitempty"`
	DestinationPort int             `json:"destination_port,omitempty"`
	Protocol        string          `json:"protocol,omitempty"`
	Domain          string          `json:"domain,omitempty"`
	URL             string          `json:"url,omitempty"`
	FileHash        string          `json:"file_hash,omitempty"`
	Username        string          `json:"username,omitempty"`
	AttackType      string          `json:"attack_type,omitempty"`
	Severity        Severity        `json:"severity"`
	Confidence      float64         `json:"confidence"`
	RawPayload      json.RawMessage `json:"raw_payload,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
}
