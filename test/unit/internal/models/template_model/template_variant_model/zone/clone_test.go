package zone_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/stretchr/testify/assert"
)

// referenceFieldCase reads a value reachable only through a slice or pointer,
// and mutates that same value in place.
type referenceFieldCase struct {
	read   func(zone template_model.Zone) any
	mutate func(zone template_model.Zone)
}

func TestWhenZoneIsCloned_ScalarFieldsAreCopied(t *testing.T) {
	t.Parallel()
	// Arrange
	zone := newPopulatedZone()

	// Act
	clone := zone.Clone()

	// Assert
	assert.Equal(t, zone, clone)
}

// TestWhenAReferenceFieldIsMutatedInPlaceOnTheClone_SourceIsUnchanged walks
// every slice and pointer reachable from a zone. The manual-edit snapshot is
// deep-copied through this method on the UI hot path, so a missed field here
// is a live aliasing bug.
func TestWhenAReferenceFieldIsMutatedInPlaceOnTheClone_SourceIsUnchanged(t *testing.T) {
	t.Parallel()
	for caseName, fieldCase := range referenceFieldCases() {
		t.Run(caseName, func(t *testing.T) {
			t.Parallel()
			// Arrange
			zone := newPopulatedZone()
			clone := zone.Clone()
			expected := fieldCase.read(zone)

			// Act
			fieldCase.mutate(clone)

			// Assert
			assert.Equal(t, expected, fieldCase.read(zone))
		})
	}
}

func referenceFieldCases() map[string]referenceFieldCase {
	return map[string]referenceFieldCase{
		"WhenQualityIsMutated_SourceIsUnchanged": {
			read:   func(zone template_model.Zone) any { return *zone.Quality },
			mutate: func(zone template_model.Zone) { *zone.Quality = neutral_zone.QualityHighest },
		},
		"WhenGeneratorPositionIsMutated_SourceIsUnchanged": {
			read:   func(zone template_model.Zone) any { return *zone.GeneratorPosition },
			mutate: func(zone template_model.Zone) { zone.GeneratorPosition.X = 9 },
		},
		"WhenGeneratorRingIsMutated_SourceIsUnchanged": {
			read:   func(zone template_model.Zone) any { return *zone.GeneratorRing },
			mutate: func(zone template_model.Zone) { *zone.GeneratorRing = 9 },
		},
		"WhenManualPositionIsMutated_SourceIsUnchanged": {
			read:   func(zone template_model.Zone) any { return *zone.ManualPosition },
			mutate: func(zone template_model.Zone) { zone.ManualPosition.Y = 9 },
		},
		"WhenEncounterHolesSettingsIsMutated_SourceIsUnchanged": {
			read:   func(zone template_model.Zone) any { return *zone.EncounterHolesSettings },
			mutate: func(zone template_model.Zone) { zone.EncounterHolesSettings.TwoHoleEncounters = 9 },
		},
		"WhenCrossroadsPositionIsMutated_SourceIsUnchanged": {
			read:   func(zone template_model.Zone) any { return *zone.CrossroadsPosition },
			mutate: func(zone template_model.Zone) { *zone.CrossroadsPosition = 9 },
		},
		"WhenGuardReactionDistributionIsMutated_SourceIsUnchanged": {
			read:   func(zone template_model.Zone) any { return zone.GuardReactionDistribution[0] },
			mutate: func(zone template_model.Zone) { zone.GuardReactionDistribution[0] = 9 },
		},
		"WhenRandomHireEnableWeeklyUnitIncrementIsMutated_SourceIsUnchanged": {
			read:   func(zone template_model.Zone) any { return zone.RandomHireEnableWeeklyUnitIncrement[0] },
			mutate: func(zone template_model.Zone) { zone.RandomHireEnableWeeklyUnitIncrement[0] = false },
		},
		"WhenRandomHireInitialUnitIncrementIsMutated_SourceIsUnchanged": {
			read:   func(zone template_model.Zone) any { return zone.RandomHireInitialUnitIncrement[0] },
			mutate: func(zone template_model.Zone) { zone.RandomHireInitialUnitIncrement[0] = 9 },
		},
		"WhenGuardedContentPoolIsMutated_SourceIsUnchanged": {
			read:   func(zone template_model.Zone) any { return zone.GuardedContentPool[0] },
			mutate: func(zone template_model.Zone) { zone.GuardedContentPool[0] = "changed" },
		},
		"WhenUnguardedContentPoolIsMutated_SourceIsUnchanged": {
			read:   func(zone template_model.Zone) any { return zone.UnguardedContentPool[0] },
			mutate: func(zone template_model.Zone) { zone.UnguardedContentPool[0] = "changed" },
		},
		"WhenResourcesContentPoolIsMutated_SourceIsUnchanged": {
			read:   func(zone template_model.Zone) any { return zone.ResourcesContentPool[0] },
			mutate: func(zone template_model.Zone) { zone.ResourcesContentPool[0] = "changed" },
		},
		"WhenMandatoryContentIsMutated_SourceIsUnchanged": {
			read:   func(zone template_model.Zone) any { return zone.MandatoryContent[0] },
			mutate: func(zone template_model.Zone) { zone.MandatoryContent[0] = "changed" },
		},
		"WhenContentCountLimitsIsMutated_SourceIsUnchanged": {
			read:   func(zone template_model.Zone) any { return zone.ContentCountLimits[0] },
			mutate: func(zone template_model.Zone) { zone.ContentCountLimits[0] = "changed" },
		},
		"WhenZoneBiomeArgsIsMutated_SourceIsUnchanged": {
			read:   func(zone template_model.Zone) any { return zone.ZoneBiome.Args[0] },
			mutate: func(zone template_model.Zone) { zone.ZoneBiome.Args[0] = "changed" },
		},
		"WhenContentBiomeArgsIsMutated_SourceIsUnchanged": {
			read:   func(zone template_model.Zone) any { return zone.ContentBiome.Args[0] },
			mutate: func(zone template_model.Zone) { zone.ContentBiome.Args[0] = "changed" },
		},
		"WhenMetaObjectsBiomeArgsIsMutated_SourceIsUnchanged": {
			read:   func(zone template_model.Zone) any { return zone.MetaObjectsBiome.Args[0] },
			mutate: func(zone template_model.Zone) { zone.MetaObjectsBiome.Args[0] = "changed" },
		},
		"WhenMainObjectIsMutated_SourceIsUnchanged": {
			read:   func(zone template_model.Zone) any { return zone.MainObjects[0].Type },
			mutate: func(zone template_model.Zone) { zone.MainObjects[0].Type = "changed" },
		},
		"WhenMainObjectFactionsIsMutated_SourceIsUnchanged": {
			read:   func(zone template_model.Zone) any { return zone.MainObjects[0].Factions[0] },
			mutate: func(zone template_model.Zone) { zone.MainObjects[0].Factions[0] = "changed" },
		},
		"WhenMainObjectPlacementArgsIsMutated_SourceIsUnchanged": {
			read:   func(zone template_model.Zone) any { return zone.MainObjects[0].PlacementArgs[0] },
			mutate: func(zone template_model.Zone) { zone.MainObjects[0].PlacementArgs[0] = "changed" },
		},
		"WhenMainObjectFactionArgsIsMutated_SourceIsUnchanged": {
			read:   func(zone template_model.Zone) any { return zone.MainObjects[0].Faction.Args[0] },
			mutate: func(zone template_model.Zone) { zone.MainObjects[0].Faction.Args[0] = "changed" },
		},
		"WhenRoadIsMutated_SourceIsUnchanged": {
			read:   func(zone template_model.Zone) any { return zone.Roads[0].Type },
			mutate: func(zone template_model.Zone) { zone.Roads[0].Type = "changed" },
		},
		"WhenRoadFromArgsIsMutated_SourceIsUnchanged": {
			read:   func(zone template_model.Zone) any { return zone.Roads[0].From.Args[0] },
			mutate: func(zone template_model.Zone) { zone.Roads[0].From.Args[0] = "changed" },
		},
		"WhenRoadToArgsIsMutated_SourceIsUnchanged": {
			read:   func(zone template_model.Zone) any { return zone.Roads[0].To.Args[0] },
			mutate: func(zone template_model.Zone) { zone.Roads[0].To.Args[0] = "changed" },
		},
		"WhenRoadRoadFlagIsMutated_SourceIsUnchanged": {
			read:   func(zone template_model.Zone) any { return *zone.Roads[0].Road },
			mutate: func(zone template_model.Zone) { *zone.Roads[0].Road = false },
		},
	}
}

// newPopulatedZone builds a zone whose every reference-typed field carries
// data, so that a missed copy in Clone shows up as shared storage.
func newPopulatedZone() template_model.Zone {
	return template_model.Zone{
		Name:              "zone",
		Quality:           new(neutral_zone.QualityLowest),
		GeneratorPosition: new(data.NewVec2(1.0, 2.0)),
		GeneratorRing:     new(3),
		ManualPosition:    new(data.NewVec2(4.0, 5.0)),
		EncounterHolesSettings: &template_model.EncounterHolesSettings{
			AffectedEncounters: 1,
			TwoHoleEncounters:  2,
		},
		CrossroadsPosition:                  new(6),
		GuardReactionDistribution:           []int{1, 2},
		RandomHireEnableWeeklyUnitIncrement: []bool{true},
		RandomHireInitialUnitIncrement:      []int{1},
		GuardedContentPool:                  []string{"guarded"},
		UnguardedContentPool:                []string{"unguarded"},
		ResourcesContentPool:                []string{"resources"},
		MandatoryContent:                    template_model.StringList{"mandatory"},
		ContentCountLimits:                  template_model.StringList{"limit"},
		ZoneBiome:                           template_model.TypedRef{Type: "zone", Args: []string{"zoneArg"}},
		ContentBiome:                        template_model.TypedRef{Type: "content", Args: []string{"contentArg"}},
		MetaObjectsBiome:                    template_model.TypedRef{Type: "meta", Args: []string{"metaArg"}},
		MainObjects: []template_model.MainObject{{
			Type:          "City",
			Factions:      []string{"faction"},
			PlacementArgs: []string{"arg"},
			Faction:       &template_model.TypedRef{Type: "FromList", Args: []string{"factionArg"}},
		}},
		Roads: []template_model.Road{{
			Type: "road",
			From: template_model.TypedRef{Type: "from", Args: []string{"fromArg"}},
			To:   template_model.TypedRef{Type: "to", Args: []string{"toArg"}},
			Road: new(true),
		}},
	}
}
