package manualEditSettings_test

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
)

// newRichZones builds a generated-shaped zone with nil and empty slices mixed
// together plus a sparse zone with a lowest-tier quality.
func newRichZones() []template_model.Zone {
	tuning := models.GenerationTuning{
		ContentScale:                   1.0,
		ResourceDensityMultiplier:      0.5,
		StructureDensityMultiplier:     1.0,
		NeutralStackStrengthMultiplier: 1.0,
		BorderGuardStrengthMultiplier:  1.0,
	}
	generated := test_helpers.NewZoneEditorService().
		NewDefaultNeutralZone("G", neutral_zone.QualityHigh, 2, true, tuning)
	generated.UnguardedContentPool = []string{}
	generated.MandatoryContent = nil
	generated.ManualPosition = new(data.NewVec2(0.3, 0.4))
	generated.GeneratorPosition = new(data.NewVec2(0.1, 0.2))
	generated.GeneratorRing = new(1)

	sparse := template_model.Zone{
		Name:    "Sparse",
		Quality: new(neutral_zone.QualityLowest),
	}
	return []template_model.Zone{generated, sparse}
}

// newRichConnections builds a populated connection and a sparse one with an
// empty (non-nil) rule list and no road flag.
func newRichConnections() []template_model.Connection {
	populated := template_model.Connection{
		Name:           "A-B",
		From:           "Zone A",
		To:             "Zone B",
		ConnectionType: "Proximity",
		Length:         2.5,
		SimTurnSquad:   true,
		Road:           new(false),
		GuardValue:     1000,
		IsUserAdded:    true,
		PortalPlacementRulesFrom: []template_model.PlacementRule{
			newRule("rule"),
		},
	}
	sparse := template_model.Connection{
		Name:                   "B-C",
		From:                   "Zone B",
		To:                     "Zone C",
		PortalPlacementRulesTo: []template_model.PlacementRule{},
	}
	return []template_model.Connection{populated, sparse}
}

func newRule(ruleType string) template_model.PlacementRule {
	var rule template_model.PlacementRule
	rule.Type = ruleType
	rule.Args = []any{"arg"}
	return rule
}
