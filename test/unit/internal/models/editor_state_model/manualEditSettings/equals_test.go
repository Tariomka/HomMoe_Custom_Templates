package manualEditSettings_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenTheSnapshotsHoldTheSameData_TheyAreEqual(t *testing.T) {
	t.Parallel()
	// Arrange
	stored := newPopulatedManualEditSettings()

	// Act
	equal := stored.Equals(newPopulatedManualEditSettings())

	// Assert
	assert.True(t, equal)
}

// The save format omits both, so the difference cannot reach the file and
// must not count as an edit.
func TestWhenATopLevelListIsEmptyInsteadOfNil_TheSnapshotsAreEqual(t *testing.T) {
	t.Parallel()
	// Arrange
	empty := editor_state_model.ManualEditSettings{
		ManualZones:       []template_model.Zone{},
		ManualConnections: []editor_state_model.ManualConnectionSave{},
	}

	// Act
	equal := empty.Equals(editor_state_model.ManualEditSettings{})

	// Assert
	assert.True(t, equal)
}

func TestWhenTheFirstSnapshotIsComparedWithAnEmptyOne_TheyAreNotEqual(t *testing.T) {
	t.Parallel()
	// Arrange
	stored := editor_state_model.ManualEditSettings{}

	// Act
	equal := stored.Equals(newPopulatedManualEditSettings())

	// Assert
	assert.False(t, equal)
}

func TestWhenAZoneIsAdded_TheSnapshotsAreNotEqual(t *testing.T) {
	t.Parallel()
	// Arrange
	stored := newPopulatedManualEditSettings()
	candidate := newPopulatedManualEditSettings()
	candidate.ManualZones = append(candidate.ManualZones, template_model.Zone{Name: "Zone B"})

	// Act
	equal := stored.Equals(candidate)

	// Assert
	assert.False(t, equal)
}

// TestWhenAPersistedFieldDiffers_TheSnapshotsAreNotEqual pins the projection
// the dirty flag is decided on. Every field the save format carries has to
// count as an edit; a converter that stops writing one trips this test.
func TestWhenAPersistedFieldDiffers_TheSnapshotsAreNotEqual(t *testing.T) {
	t.Parallel()
	for caseName, mutate := range persistedFieldCases() {
		t.Run(caseName, func(t *testing.T) {
			t.Parallel()
			// Arrange
			stored := newPopulatedManualEditSettings()
			candidate := newPopulatedManualEditSettings()

			// Act
			mutate(candidate)

			// Assert
			assert.False(t, stored.Equals(candidate))
		})
	}
}

// persistedFieldCases mutates one saved field per case, in place through the
// snapshot's own storage.
func persistedFieldCases() map[string]func(editor_state_model.ManualEditSettings) {
	return map[string]func(editor_state_model.ManualEditSettings){
		"WhenTheZoneNameDiffers_TheSnapshotsAreNotEqual": func(settings editor_state_model.ManualEditSettings) {
			settings.ManualZones[0].Name = "changed"
		},
		"WhenTheZoneSizeDiffers_TheSnapshotsAreNotEqual": func(settings editor_state_model.ManualEditSettings) {
			settings.ManualZones[0].Size = 9
		},
		"WhenTheManualPositionDiffers_TheSnapshotsAreNotEqual": func(
			settings editor_state_model.ManualEditSettings) {
			settings.ManualZones[0].ManualPosition.X = 9
		},
		"WhenTheGeneratorPositionDiffers_TheSnapshotsAreNotEqual": func(
			settings editor_state_model.ManualEditSettings) {
			settings.ManualZones[0].GeneratorPosition.Y = 9
		},
		"WhenTheGeneratorRingDiffers_TheSnapshotsAreNotEqual": func(
			settings editor_state_model.ManualEditSettings) {
			*settings.ManualZones[0].GeneratorRing = 9
		},
		"WhenTheQualityDiffers_TheSnapshotsAreNotEqual": func(settings editor_state_model.ManualEditSettings) {
			*settings.ManualZones[0].Quality = neutral_zone.QualityHighest
		},
		"WhenTheEncounterHolesDiffer_TheSnapshotsAreNotEqual": func(
			settings editor_state_model.ManualEditSettings) {
			settings.ManualZones[0].EncounterHolesSettings.TwoHoleEncounters = 9
		},
		"WhenAContentPoolEntryDiffers_TheSnapshotsAreNotEqual": func(
			settings editor_state_model.ManualEditSettings) {
			settings.ManualZones[0].GuardedContentPool[0] = "changed"
		},
		"WhenAMainObjectDiffers_TheSnapshotsAreNotEqual": func(settings editor_state_model.ManualEditSettings) {
			settings.ManualZones[0].MainObjects[0].Factions[0] = "changed"
		},
		"WhenAZoneRoadDiffers_TheSnapshotsAreNotEqual": func(settings editor_state_model.ManualEditSettings) {
			settings.ManualZones[0].Roads[0].Type = "changed"
		},
		"WhenTheConnectionLengthDiffers_TheSnapshotsAreNotEqual": func(
			settings editor_state_model.ManualEditSettings) {
			settings.ManualConnections[0].Connection.Length = 9
		},
		"WhenTheConnectionRoadFlagDiffers_TheSnapshotsAreNotEqual": func(
			settings editor_state_model.ManualEditSettings) {
			*settings.ManualConnections[0].Connection.Road = false
		},
		"WhenTheConnectionGuardValueDiffers_TheSnapshotsAreNotEqual": func(
			settings editor_state_model.ManualEditSettings) {
			settings.ManualConnections[0].Connection.GuardValue = 9
		},
		"WhenAPortalPlacementRuleDiffers_TheSnapshotsAreNotEqual": func(
			settings editor_state_model.ManualEditSettings) {
			settings.ManualConnections[0].Connection.PortalPlacementRulesFrom[0].Type = "changed"
		},
		"WhenIsUserAddedDiffers_TheSnapshotsAreNotEqual": func(settings editor_state_model.ManualEditSettings) {
			settings.ManualConnections[0].IsUserAdded = false
		},
	}
}

// newPopulatedManualEditSettings builds a snapshot whose saved fields all
// carry data, so a comparison that skips one shows up as a missed edit.
func newPopulatedManualEditSettings() editor_state_model.ManualEditSettings {
	return editor_state_model.ManualEditSettings{
		ManualZones: []template_model.Zone{{
			Name:              "Zone A",
			Size:              1.5,
			Quality:           new(neutral_zone.QualityLowest),
			GeneratorPosition: new(data.NewVec2(1.0, 2.0)),
			GeneratorRing:     new(3),
			ManualPosition:    new(data.NewVec2(4.0, 5.0)),
			EncounterHolesSettings: &template_model.EncounterHolesSettings{
				AffectedEncounters: 1,
				TwoHoleEncounters:  2,
			},
			GuardedContentPool: []string{"guarded"},
			MainObjects: []template_model.MainObject{{
				Type:     "City",
				Factions: []string{"faction"},
			}},
			Roads: []template_model.Road{{Type: "road", Road: new(true)}},
		}},
		ManualConnections: editor_state_model.ToManualConnectionSaves([]template_model.Connection{{
			Name:                     "A-B",
			From:                     "Zone A",
			To:                       "Zone B",
			Length:                   2.5,
			GuardValue:               1000,
			Road:                     new(true),
			PortalPlacementRulesFrom: []template_model.PlacementRule{{Type: "rule"}},
			IsUserAdded:              true,
		}}),
	}
}
