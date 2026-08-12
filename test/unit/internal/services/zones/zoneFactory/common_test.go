package zoneFactory_test

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones/zone_interfaces"
)

func newZoneFactory() zone_interfaces.IZoneFactory {
	return zones.NewZoneFactory(zones.NewCastleFactory(), zones.NewRoadFactory())
}

func newUnitTuning() models.GenerationTuning {
	return models.GenerationTuning{
		ContentScale:                   1,
		ResourceDensityMultiplier:      1,
		StructureDensityMultiplier:     1,
		NeutralStackStrengthMultiplier: 1,
		BorderGuardStrengthMultiplier:  1,
	}
}
