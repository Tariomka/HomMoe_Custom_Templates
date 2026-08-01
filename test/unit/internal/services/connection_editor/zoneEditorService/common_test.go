// Package zoneEditor_test contains shared arrangement helpers for the
// zoneEditorService.go unit tests.
package zoneEditorService_test

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

// defaultTuning returns a neutral generation tuning so profile values are not
// scaled by density/guard multipliers.
func defaultTuning() models.GenerationTuning {
	return models.GenerationTuning{
		ContentScale:                   1.0,
		ResourceDensityMultiplier:      0.5,
		StructureDensityMultiplier:     1.0,
		NeutralStackStrengthMultiplier: 1.0,
		BorderGuardStrengthMultiplier:  1.0,
	}
}

// roadTargets collects the first argument of every road endpoint of the given
// TypedRef type, from both road directions.
func roadTargets(zone entities.Zone, refType string) map[string]bool {
	targets := map[string]bool{}
	for _, road := range zone.Roads {
		if road.To.Type == refType && len(road.To.Args) > 0 {
			targets[road.To.Args[0]] = true
		}
		if road.From.Type == refType && len(road.From.Args) > 0 {
			targets[road.From.Args[0]] = true
		}
	}
	return targets
}

// castleRoadTargets returns the MainObject indices linked by the zone's stone
// castle<->castle roads from the primary main object.
func castleRoadTargets(zone entities.Zone) []string {
	var targets []string
	for _, road := range zone.Roads {
		if road.From.Type == "MainObject" && road.To.Type == "MainObject" && len(road.To.Args) > 0 {
			targets = append(targets, road.To.Args[0])
		}
	}
	return targets
}

// mainObjectZeroRef is the road endpoint pointing at a zone's primary main object.
func mainObjectZeroRef() entities.TypedRef {
	return entities.TypedRef{Type: "MainObject", Args: []string{"0"}}
}
