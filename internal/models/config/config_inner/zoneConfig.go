package config_inner

import (
	"math"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
)

type ZoneConfig struct {
	NeutralZoneCount int

	PlayerOwnedCastles int
	PlayerZoneCastles  int
	NeutralZoneCastles int

	SpawnAbandonedOutposts bool
	AbandonedOutpostCount  int

	ResourceDensityPercent  int
	StructureDensityPercent int

	NeutralStackStrengthPercent int
	BorderGuardStrengthPercent  int
	GuardRandomization          float64

	HubZoneSize     float64
	PlayerZoneSize  float64
	NeutralZoneSize float64

	Advanced AdvancedSettings
}

func (this ZoneConfig) GetEffectiveGuardRandomization() float64 {
	if !this.Advanced.Enabled {
		return 0.05
	}

	randomizationValue := this.GuardRandomization
	if math.IsNaN(randomizationValue) || math.IsInf(randomizationValue, 0) {
		return 0.05
	}

	return helpers.RoundWithPrecision(math.Max(0, math.Min(randomizationValue, 0.5)), 3)
}
