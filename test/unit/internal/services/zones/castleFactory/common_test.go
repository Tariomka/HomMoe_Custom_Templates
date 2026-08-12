package castleFactory_test

import "github.com/Tariomka/hommoe_custom_templates/internal/models"

func newUnitTuning() models.GenerationTuning {
	return models.GenerationTuning{
		ContentScale:                   1,
		ResourceDensityMultiplier:      1,
		StructureDensityMultiplier:     1,
		NeutralStackStrengthMultiplier: 1,
		BorderGuardStrengthMultiplier:  1,
	}
}
