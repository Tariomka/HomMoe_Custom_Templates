package models

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
)

type GenerationTuning struct {
	ContentScale                   float64
	ResourceDensityMultiplier      float64
	StructureDensityMultiplier     float64
	NeutralStackStrengthMultiplier float64
	BorderGuardStrengthMultiplier  float64
	GuardRandomization             float64
	RemoteFootholdCount            int
	AbandonedOutpostCount          int
	PlayerOwnedCastles             int
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
