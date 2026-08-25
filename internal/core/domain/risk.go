package domain

type RiskLevel string

const (
	RiskLevelInformational RiskLevel = "informational" // 0-29
	RiskLevelLow           RiskLevel = "low"           // 30-49
	RiskLevelMedium        RiskLevel = "medium"        // 50-69
	RiskLevelHigh          RiskLevel = "high"          // 70-84
	RiskLevelCritical      RiskLevel = "critical"      // 85-100
)

type RiskScore struct {
	TotalScore  float64   `json:"total_score"`  // 0 - 100
	Level       RiskLevel `json:"level"`
	Confidence  float64   `json:"confidence"`   // 0.0 - 1.0 (or percentage 0-100)
	Breakdown   RiskBreakdown `json:"breakdown"`
}

type RiskBreakdown struct {
	Reputation  float64 `json:"reputation"`   // max 30
	Severity    float64 `json:"severity"`     // max 25
	Frequency   float64 `json:"frequency"`    // max 20
	Confidence  float64 `json:"confidence"`   // max 15
	Correlation float64 `json:"correlation"`  // max 10
}

func CalculateRiskLevel(score float64) RiskLevel {
	switch {
	case score >= 85:
		return RiskLevelCritical
	case score >= 70:
		return RiskLevelHigh
	case score >= 50:
		return RiskLevelMedium
	case score >= 30:
		return RiskLevelLow
	default:
		return RiskLevelInformational
	}
}
