export type Severity = "critical" | "high" | "medium" | "low" | "informational"

export type AttackType =
  | "ssh_bruteforce"
  | "port_scan"
  | "credential_stuffing"
  | "sqli"
  | "xss"
  | "path_traversal"
  | "malware_download"
  | "suspicious_command"
  | string

export type IndicatorType = "ip" | "domain" | "url" | "sha256" | "sha1" | "md5" | "email" | "asn"

export type IndicatorStatus = "active" | "expired" | "whitelisted"

export type AlertStatus = "open" | "investigating" | "resolved" | "false_positive"

export interface Event {
  id: string
  timestamp: string
  source: string
  source_ip?: string
  destination_ip?: string
  source_port?: number
  destination_port?: number
  protocol?: string
  domain?: string
  url?: string
  file_hash?: string
  username?: string
  attack_type?: string
  severity: Severity
  confidence: number
  raw_payload?: Record<string, unknown>
  created_at: string
}

export interface Indicator {
  id: string
  type: IndicatorType
  value: string
  first_seen: string
  last_seen: string
  risk_score: number
  confidence: number
  status: IndicatorStatus
  created_at: string
  updated_at: string
}

export interface Observation {
  id: string
  indicator_id: string
  event_id: string
  attack_type: string
  technique_id?: string
  timestamp: string
  severity: Severity
  confidence: number
  source?: string
  created_at: string
}

export interface GeoIPData {
  country?: string
  country_code?: string
  city?: string
  latitude?: number
  longitude?: number
  asn?: string
  org?: string
  isp?: string
}

export interface ThreatIntelData {
  malicious_votes?: number
  harmless_votes?: number
  reputation?: number
  tags?: string[]
  malware_family?: string
  last_analysis?: string
}

export interface Enrichment {
  id: string
  indicator_id: string
  provider: "geoip" | "virustotal" | "abuseipdb" | string
  data: GeoIPData | ThreatIntelData | Record<string, unknown>
  fetched_at: string
  expires_at?: string
}

export interface Relationship {
  id: string
  source_indicator_id: string
  target_indicator_id: string
  relationship_type: "IP->DOMAIN" | "IP->HASH" | "IP->ATTACK" | "DOMAIN->IP" | "HASH->MALWARE" | string
  confidence: number
  first_seen: string
  last_seen: string
}

export interface Alert {
  id: string
  indicator_id: string
  event_id?: string
  severity: Severity
  risk_score: number
  confidence: number
  title: string
  description?: string
  status: AlertStatus
  created_at: string
  resolved_at?: string
}

export interface RiskBreakdown {
  reputation: number
  severity: number
  frequency: number
  confidence: number
  correlation: number
}

export interface RiskScore {
  total_score: number
  level: "critical" | "high" | "medium" | "low" | "informational"
  confidence: number
  breakdown: RiskBreakdown
}

export interface InvestigationContext {
  indicator: Indicator
  observations: Observation[]
  enrichments: Enrichment[]
  relationships: Relationship[]
}

export interface DashboardStats {
  total_events?: number
  total_indicators?: number
  alerts_by_status?: Record<string, number>
}
