// Package manualZoneEdit_test contains shared arrangement helpers for the
// manualZoneEdit.go unit tests.
package manualZoneEdit_test

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
)

func defaultTuning() models.GenerationTuning {
	return models.GenerationTuning{
		ContentScale:                   1.0,
		ResourceDensityMultiplier:      0.5,
		StructureDensityMultiplier:     1.0,
		NeutralStackStrengthMultiplier: 1.0,
		BorderGuardStrengthMultiplier:  1.0,
	}
}

// newDerivedZone builds a generator-shaped zone to diff against.
func newDerivedZone() template_model.Zone {
	return test_helpers.NewZoneEditorService().
		NewDefaultNeutralZone("D", neutral_zone.QualityLow, 1, false, defaultTuning())
}

// newRichEditedZone builds a zone that differs from newDerivedZone in tier,
// castles, roads, scalars, pools and structure fields.
func newRichEditedZone() template_model.Zone {
	zone := test_helpers.NewZoneEditorService().
		NewDefaultNeutralZone("G", neutral_zone.QualityHigh, 2, true, defaultTuning())
	zone.Size = 3.5
	zone.Layout = "edited_layout"
	zone.GuardCutoffValue = 0
	zone.DiplomacyModifier = 0.25
	zone.GuardedContentPool = []string{"edited_pool_a", "edited_pool_b"}
	zone.UnguardedContentPool = []string{}
	zone.EncounterHolesSettings = &template_model.EncounterHolesSettings{
		AffectedEncounters: 1,
		TwoHoleEncounters:  2,
	}
	zone.ZoneBiome = template_model.TypedRef{}
	zone.ZoneBiome.Type = "FromList"
	zone.ZoneBiome.Args = []string{"biome"}
	zone.CrossroadsPosition = new(2)
	zone.GeneratorPosition = new(data.NewVec2(0.1, 0.2))
	zone.GeneratorRing = new(2)
	zone.ManualPosition = new(data.NewVec2(0.3, 0.4))
	return zone
}

// newPopulatedEdit builds an added-zone edit whose every reference field carries data.
func newPopulatedEdit() editor_state_model.ManualZoneEdit {
	edit := editor_state_model.NewManualZoneEdit(newDerivedZone(), newPopulatedZone(), true)
	edit.FieldOverrides.MainObjects = []template_entity.MainObject{{Type: "City", Factions: []string{"faction"}}}
	edit.Unset = []string{"layout"}
	edit = edit.AsAdded(2)
	return edit
}

// newPopulatedZone builds a zone whose every reference field carries data.
func newPopulatedZone() template_model.Zone {
	return template_model.Zone{
		Name:               "Zone A",
		Quality:            new(neutral_zone.QualityHigh),
		GeneratorPosition:  new(data.NewVec2(1.0, 2.0)),
		GeneratorRing:      new(3),
		ManualPosition:     new(data.NewVec2(4.0, 5.0)),
		Size:               1.5,
		GuardedContentPool: []string{"guarded"},
		EncounterHolesSettings: &template_model.EncounterHolesSettings{
			AffectedEncounters: 1,
			TwoHoleEncounters:  2,
		},
		MainObjects: []template_model.MainObject{{
			Type:     "City",
			Factions: []string{"faction"},
		}},
		Roads: []template_model.Road{{Type: "road", Road: new(true)}},
	}
}
