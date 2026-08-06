// Package gladiatorArenaProvider_test contains shared arrangement helpers for
// the gladiatorArenaProvider.go unit tests.
package gladiatorArenaProvider_test

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers"
	zone_services "github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
)

const arenaObjectType = "GladiatorArena"

func newProvider() *providers.GladiatorArenaProvider {
	return providers.NewGladiatorArenaProvider(zone_services.NewZoneClassifier())
}

// newArenaConfiguration returns a configuration that asks for the arena win
// condition, so PlaceArena is expected to act on it.
func newArenaConfiguration() config.GeneratorConfig {
	configuration := config.NewGeneratorConfig()
	configuration.GladiatorArenaRules.Enabled = true
	return *configuration
}

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

// newNeutralZone builds a generator-shaped neutral zone whose content pools
// classify back to the requested quality.
func newNeutralZone(label string, quality neutral_zone.Quality) entities.Zone {
	return test_helpers.NewZoneEditorService().
		NewDefaultNeutralZone(label, quality, 0, false, defaultTuning())
}

func countArenaMainObjects(zone entities.Zone) int {
	count := 0
	for _, mainObject := range zone.MainObjects {
		if mainObject.Type == arenaObjectType {
			count++
		}
	}
	return count
}

func countArenaConnections(variant entities.Variant) int {
	count := 0
	for _, connection := range variant.Connections {
		if connection.ConnectionType == arenaObjectType {
			count++
		}
	}
	return count
}
