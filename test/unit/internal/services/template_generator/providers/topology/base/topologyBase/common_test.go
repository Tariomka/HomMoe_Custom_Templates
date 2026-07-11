package topologyBase_test

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

// newUnitTuning builds tuning with every multiplier at 1.0 so that guard and
// content values pass through the scaling helpers unchanged, keeping the
// expected values in assertions readable.
func newUnitTuning() models.GenerationTuning {
	return models.GenerationTuning{
		ContentScale:                   1.0,
		ResourceDensityMultiplier:      1.0,
		StructureDensityMultiplier:     1.0,
		NeutralStackStrengthMultiplier: 1.0,
		BorderGuardStrengthMultiplier:  1.0,
		GuardRandomization:             0.05,
	}
}
