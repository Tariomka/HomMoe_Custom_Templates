package topologyConnectionService_test

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

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

func connectionRoadTargets(zoneList []template_model.Zone) []string {
	var names []string
	for _, zone := range zoneList {
		for _, road := range zone.Roads {
			if road.From.Type == "Connection" && len(road.From.Args) > 0 {
				names = append(names, road.From.Args[0])
			}
			if road.To.Type == "Connection" && len(road.To.Args) > 0 {
				names = append(names, road.To.Args[0])
			}
		}
	}
	return names
}
