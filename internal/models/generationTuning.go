package models

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/utils"
)

type GenerationTuning struct {
	ContentScale                   float64
	ResourceDensityMultiplier      float64
	StructureDensityMultiplier     float64
	NeutralStackStrengthMultiplier float64
	BorderGuardStrengthMultiplier  float64
	GuardRandomization             float64
	SpawnAbandonedOutposts         bool
	PlayerOwnedCastles             int
}

// NewGenerationTuning builds the content/guard scaling factors for the given configuration
func NewGenerationTuning(configuration *config.GeneratorConfig, totalZoneCount int) GenerationTuning {
	return GenerationTuning{
		ContentScale:                   utils.ComputeContentScale(configuration.MapSize, totalZoneCount),
		ResourceDensityMultiplier:      float64(configuration.ZoneConfiguration.ResourceDensityPercent) / 200.0,
		StructureDensityMultiplier:     float64(configuration.ZoneConfiguration.StructureDensityPercent) / 100.0,
		NeutralStackStrengthMultiplier: float64(configuration.ZoneConfiguration.NeutralStackStrengthPercent) / 100.0,
		BorderGuardStrengthMultiplier:  float64(configuration.ZoneConfiguration.BorderGuardStrengthPercent) / 100.0,
		GuardRandomization:             configuration.ZoneConfiguration.Advanced.GetEffectiveGuardRandomization(),
		SpawnAbandonedOutposts:         configuration.ZoneConfiguration.SpawnAbandonedOutposts,
		PlayerOwnedCastles:             configuration.ZoneConfiguration.PlayerOwnedCastles,
	}
}

func (this GenerationTuning) ScaleByStructureDensity(value float64) int {
	return helpers.Scale(value, this.StructureDensityMultiplier)
}

func (this GenerationTuning) ScaleByResourceDensity(value float64) int {
	return helpers.Scale(value, this.ResourceDensityMultiplier)
}

func (this GenerationTuning) ScaleByNeutralGuardStrength(value int) int {
	return helpers.Scale(float64(value), this.NeutralStackStrengthMultiplier)
}

func (this GenerationTuning) ScaleByBorderGuardStrength(value int) int {
	return helpers.Scale(float64(value), this.BorderGuardStrengthMultiplier)
}

func (this GenerationTuning) ScaleByNeutralGuardStrengthPrecise(value float64) float64 {
	return helpers.RoundWithPrecision(value*this.NeutralStackStrengthMultiplier, 3)
}
