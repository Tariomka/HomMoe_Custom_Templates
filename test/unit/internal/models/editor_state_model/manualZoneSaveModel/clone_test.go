package manualZoneSaveModel_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
)

// referenceFieldCase reads a value reachable only through a slice or pointer,
// and mutates that same value in place.
type referenceFieldCase struct {
	read   func(save editor_state_model.ManualZoneSaveModel) any
	mutate func(save editor_state_model.ManualZoneSaveModel)
}

func TestWhenSaveIsCloned_ScalarFieldsAreCopied(t *testing.T) {
	t.Parallel()
	// Arrange
	save := newPopulatedSave()

	// Act
	clone := save.Clone()

	// Assert
	assert.Equal(t, save, clone)
}

// TestWhenAReferenceFieldIsMutatedInPlaceOnTheClone_SourceIsUnchanged walks
// every slice and pointer reachable from entities.Zone. That entity lives in
// the protected template tree and cannot carry a Clone of its own, so this is
// the only place its copy semantics are pinned.
func TestWhenAReferenceFieldIsMutatedInPlaceOnTheClone_SourceIsUnchanged(t *testing.T) {
	t.Parallel()
	for caseName, fieldCase := range referenceFieldCases() {
		t.Run(caseName, func(t *testing.T) {
			t.Parallel()
			// Arrange
			save := newPopulatedSave()
			clone := save.Clone()
			expected := fieldCase.read(save)

			// Act
			fieldCase.mutate(clone)

			// Assert
			assert.Equal(t, expected, fieldCase.read(save))
		})
	}
}

func referenceFieldCases() map[string]referenceFieldCase {
	return map[string]referenceFieldCase{
		"WhenSaveManualPositionIsMutated_SourceIsUnchanged": {
			read:   func(save editor_state_model.ManualZoneSaveModel) any { return *save.ManualPosition },
			mutate: func(save editor_state_model.ManualZoneSaveModel) { save.ManualPosition[0] = 9 },
		},
		"WhenGeneratorPositionIsMutated_SourceIsUnchanged": {
			read:   func(save editor_state_model.ManualZoneSaveModel) any { return *save.Zone.GeneratorPosition },
			mutate: func(save editor_state_model.ManualZoneSaveModel) { save.Zone.GeneratorPosition[0] = 9 },
		},
		"WhenGeneratorRingIsMutated_SourceIsUnchanged": {
			read:   func(save editor_state_model.ManualZoneSaveModel) any { return *save.Zone.GeneratorRing },
			mutate: func(save editor_state_model.ManualZoneSaveModel) { *save.Zone.GeneratorRing = 9 },
		},
		"WhenZoneManualPositionIsMutated_SourceIsUnchanged": {
			read:   func(save editor_state_model.ManualZoneSaveModel) any { return *save.Zone.ManualPosition },
			mutate: func(save editor_state_model.ManualZoneSaveModel) { save.Zone.ManualPosition[1] = 9 },
		},
		"WhenEncounterHolesSettingsIsMutated_SourceIsUnchanged": {
			read: func(save editor_state_model.ManualZoneSaveModel) any { return *save.Zone.EncounterHolesSettings },
			mutate: func(save editor_state_model.ManualZoneSaveModel) {
				save.Zone.EncounterHolesSettings.TwoHoleEncounters = 9
			},
		},
		"WhenCrossroadsPositionIsMutated_SourceIsUnchanged": {
			read:   func(save editor_state_model.ManualZoneSaveModel) any { return *save.Zone.CrossroadsPosition },
			mutate: func(save editor_state_model.ManualZoneSaveModel) { *save.Zone.CrossroadsPosition = 9 },
		},
		"WhenGuardReactionDistributionIsMutated_SourceIsUnchanged": {
			read:   func(save editor_state_model.ManualZoneSaveModel) any { return save.Zone.GuardReactionDistribution[0] },
			mutate: func(save editor_state_model.ManualZoneSaveModel) { save.Zone.GuardReactionDistribution[0] = 9 },
		},
		"WhenRandomHireEnableWeeklyUnitIncrementIsMutated_SourceIsUnchanged": {
			read: func(save editor_state_model.ManualZoneSaveModel) any {
				return save.Zone.RandomHireEnableWeeklyUnitIncrement[0]
			},
			mutate: func(save editor_state_model.ManualZoneSaveModel) {
				save.Zone.RandomHireEnableWeeklyUnitIncrement[0] = false
			},
		},
		"WhenRandomHireInitialUnitIncrementIsMutated_SourceIsUnchanged": {
			read: func(save editor_state_model.ManualZoneSaveModel) any {
				return save.Zone.RandomHireInitialUnitIncrement[0]
			},
			mutate: func(save editor_state_model.ManualZoneSaveModel) { save.Zone.RandomHireInitialUnitIncrement[0] = 9 },
		},
		"WhenGuardedContentPoolIsMutated_SourceIsUnchanged": {
			read:   func(save editor_state_model.ManualZoneSaveModel) any { return save.Zone.GuardedContentPool[0] },
			mutate: func(save editor_state_model.ManualZoneSaveModel) { save.Zone.GuardedContentPool[0] = "changed" },
		},
		"WhenUnguardedContentPoolIsMutated_SourceIsUnchanged": {
			read:   func(save editor_state_model.ManualZoneSaveModel) any { return save.Zone.UnguardedContentPool[0] },
			mutate: func(save editor_state_model.ManualZoneSaveModel) { save.Zone.UnguardedContentPool[0] = "changed" },
		},
		"WhenResourcesContentPoolIsMutated_SourceIsUnchanged": {
			read:   func(save editor_state_model.ManualZoneSaveModel) any { return save.Zone.ResourcesContentPool[0] },
			mutate: func(save editor_state_model.ManualZoneSaveModel) { save.Zone.ResourcesContentPool[0] = "changed" },
		},
		"WhenMandatoryContentIsMutated_SourceIsUnchanged": {
			read:   func(save editor_state_model.ManualZoneSaveModel) any { return save.Zone.MandatoryContent[0] },
			mutate: func(save editor_state_model.ManualZoneSaveModel) { save.Zone.MandatoryContent[0] = "changed" },
		},
		"WhenContentCountLimitsIsMutated_SourceIsUnchanged": {
			read:   func(save editor_state_model.ManualZoneSaveModel) any { return save.Zone.ContentCountLimits[0] },
			mutate: func(save editor_state_model.ManualZoneSaveModel) { save.Zone.ContentCountLimits[0] = "changed" },
		},
		"WhenZoneBiomeArgsIsMutated_SourceIsUnchanged": {
			read:   func(save editor_state_model.ManualZoneSaveModel) any { return save.Zone.ZoneBiome.Args[0] },
			mutate: func(save editor_state_model.ManualZoneSaveModel) { save.Zone.ZoneBiome.Args[0] = "changed" },
		},
		"WhenContentBiomeArgsIsMutated_SourceIsUnchanged": {
			read:   func(save editor_state_model.ManualZoneSaveModel) any { return save.Zone.ContentBiome.Args[0] },
			mutate: func(save editor_state_model.ManualZoneSaveModel) { save.Zone.ContentBiome.Args[0] = "changed" },
		},
		"WhenMetaObjectsBiomeArgsIsMutated_SourceIsUnchanged": {
			read:   func(save editor_state_model.ManualZoneSaveModel) any { return save.Zone.MetaObjectsBiome.Args[0] },
			mutate: func(save editor_state_model.ManualZoneSaveModel) { save.Zone.MetaObjectsBiome.Args[0] = "changed" },
		},
		"WhenMainObjectIsMutated_SourceIsUnchanged": {
			read:   func(save editor_state_model.ManualZoneSaveModel) any { return save.Zone.MainObjects[0].Type },
			mutate: func(save editor_state_model.ManualZoneSaveModel) { save.Zone.MainObjects[0].Type = "changed" },
		},
		"WhenMainObjectFactionsIsMutated_SourceIsUnchanged": {
			read:   func(save editor_state_model.ManualZoneSaveModel) any { return save.Zone.MainObjects[0].Factions[0] },
			mutate: func(save editor_state_model.ManualZoneSaveModel) { save.Zone.MainObjects[0].Factions[0] = "changed" },
		},
		"WhenMainObjectPlacementArgsIsMutated_SourceIsUnchanged": {
			read: func(save editor_state_model.ManualZoneSaveModel) any {
				return save.Zone.MainObjects[0].PlacementArgs[0]
			},
			mutate: func(save editor_state_model.ManualZoneSaveModel) {
				save.Zone.MainObjects[0].PlacementArgs[0] = "changed"
			},
		},
		"WhenMainObjectFactionArgsIsMutated_SourceIsUnchanged": {
			read: func(save editor_state_model.ManualZoneSaveModel) any { return save.Zone.MainObjects[0].Faction.Args[0] },
			mutate: func(save editor_state_model.ManualZoneSaveModel) {
				save.Zone.MainObjects[0].Faction.Args[0] = "changed"
			},
		},
		"WhenRoadIsMutated_SourceIsUnchanged": {
			read:   func(save editor_state_model.ManualZoneSaveModel) any { return save.Zone.Roads[0].Type },
			mutate: func(save editor_state_model.ManualZoneSaveModel) { save.Zone.Roads[0].Type = "changed" },
		},
		"WhenRoadFromArgsIsMutated_SourceIsUnchanged": {
			read:   func(save editor_state_model.ManualZoneSaveModel) any { return save.Zone.Roads[0].From.Args[0] },
			mutate: func(save editor_state_model.ManualZoneSaveModel) { save.Zone.Roads[0].From.Args[0] = "changed" },
		},
		"WhenRoadToArgsIsMutated_SourceIsUnchanged": {
			read:   func(save editor_state_model.ManualZoneSaveModel) any { return save.Zone.Roads[0].To.Args[0] },
			mutate: func(save editor_state_model.ManualZoneSaveModel) { save.Zone.Roads[0].To.Args[0] = "changed" },
		},
		"WhenRoadRoadFlagIsMutated_SourceIsUnchanged": {
			read:   func(save editor_state_model.ManualZoneSaveModel) any { return *save.Zone.Roads[0].Road },
			mutate: func(save editor_state_model.ManualZoneSaveModel) { *save.Zone.Roads[0].Road = false },
		},
	}
}

// newPopulatedSave builds a save whose every reference-typed field carries data,
// so that a missed copy in cloneZone shows up as shared storage.
func newPopulatedSave() editor_state_model.ManualZoneSaveModel {
	zone := entities.Zone{
		Name:              "zone",
		GeneratorPosition: &[2]float64{1, 2},
		GeneratorRing:     new(3),
		ManualPosition:    &[2]float64{4, 5},
		EncounterHolesSettings: &entities.EncounterHolesSettings{
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
		MandatoryContent:                    entities.StringList{"mandatory"},
		ContentCountLimits:                  entities.StringList{"limit"},
		ZoneBiome:                           entities.TypedRef{Type: "zone", Args: []string{"zoneArg"}},
		ContentBiome:                        entities.TypedRef{Type: "content", Args: []string{"contentArg"}},
		MetaObjectsBiome:                    entities.TypedRef{Type: "meta", Args: []string{"metaArg"}},
		MainObjects: []entities.MainObject{{
			Type:          "City",
			Factions:      []string{"faction"},
			PlacementArgs: []string{"arg"},
			Faction:       &entities.TypedRef{Type: "FromList", Args: []string{"factionArg"}},
		}},
		Roads: []entities.Road{{
			Type: "road",
			From: entities.TypedRef{Type: "from", Args: []string{"fromArg"}},
			To:   entities.TypedRef{Type: "to", Args: []string{"toArg"}},
			Road: new(true),
		}},
	}

	return editor_state_model.ManualZoneSaveModel{
		Zone: zone, ManualPosition: &[2]float64{7, 8},
	}
}
