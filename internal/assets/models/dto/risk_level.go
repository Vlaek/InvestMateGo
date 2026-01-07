package dto

type RiskLevel string

const (
	RiskLevelUnspecified RiskLevel = "RISK_LEVEL_UNSPECIFIED"
	RiskLevelLow         RiskLevel = "RISK_LEVEL_LOW"
	RiskLevelModerate    RiskLevel = "RISK_LEVEL_MODERATE"
	RiskLevelHigh        RiskLevel = "RISK_LEVEL_HIGH"
)

func (rl RiskLevel) ToString() string {
	return string(rl)
}
