// Package manualReapply_test contains shared arrangement helpers for the
// manualReapply.go unit tests.
package manualReapplyService_test

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/generation_tuning"
	zone_services "github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
)

// newManualReapplyService builds the service with the same collaborators the
// application wires.
func newManualReapplyService() connection_editor.IManualReapplyService {
	return connection_editor.NewManualReapplyService(
		test_helpers.NewZoneEditorService(),
		zone_services.NewCastleFactory(),
		zone_services.NewZoneClassifier(),
		generation_tuning.NewGenerationTuningFactory(),
	)
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

// makeNeutralZone builds a generator-shaped neutral zone for the given quality
// and castle count.
func makeNeutralZone(label string, quality neutral_zone.Quality, castleCount int) entities.Zone {
	return test_helpers.NewZoneEditorService().
		NewDefaultNeutralZone(label, quality, castleCount, false, defaultTuning())
}

// makeSpawnZone builds a player spawn zone with the spawn castle as the primary
// main object plus the requested extra city castles.
func makeSpawnZone(label, playerName string, extraCastleCount int) entities.Zone {
	mainObjects := []entities.MainObject{{Type: "Spawn", Spawn: playerName}}
	for range extraCastleCount {
		mainObjects = append(mainObjects, entities.MainObject{Type: "City"})
	}
	return entities.Zone{Name: "Spawn-" + label, MainObjects: mainObjects}
}
