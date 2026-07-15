package config_inner

import (
	"math"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
)

type AdvancedSettings struct {
	Enabled                     bool
	NeutralLowestNoCastleCount  int
	NeutralLowestCastleCount    int
	NeutralLowNoCastleCount     int
	NeutralLowCastleCount       int
	NeutralMediumNoCastleCount  int
	NeutralMediumCastleCount    int
	NeutralHighNoCastleCount    int
	NeutralHighCastleCount      int
	NeutralLowestCastlesPerZone int
	NeutralLowCastlesPerZone    int
	NeutralMediumCastlesPerZone int
	NeutralHighCastlesPerZone   int
	PlayerZoneSize              float64
	NeutralZoneSize             float64
	GuardRandomization          float64
}

func (this AdvancedSettings) GetEffectiveGuardRandomization() float64 {
	if !this.Enabled {
		return 0.05
	}

	randomizationValue := this.GuardRandomization
	if math.IsNaN(randomizationValue) || math.IsInf(randomizationValue, 0) {
		return 0.05
	}

	return helpers.RoundWithPrecision(math.Max(0, math.Min(randomizationValue, 0.5)), 3)
}
