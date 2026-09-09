package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenManualZonesWereStored_ZonesRoundTripWithManualPositions(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	zones := []template_model.Zone{
		{
			Name:           "Zone A",
			Size:           gofakeit.Float64Range(0.5, 2),
			ManualPosition: new(data.NewVec2(gofakeit.Float64Range(0, 1), gofakeit.Float64Range(0, 1))),
		},
		{Name: "Zone B"},
	}
	state.SetManualEdits(zones, nil)

	// Act
	restored := state.GetManualZones()

	// Assert
	assert.Equal(t, zones, restored)
}

func TestWhenNoManualZonesWereStored_NilZonesAreReturned(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()

	// Act
	restored := state.GetManualZones()

	// Assert
	assert.Nil(t, restored)
}

// Returned zones are the editor's working copies; editing them must not
// rewrite the committed snapshot behind the state's back.
func TestWhenAReturnedZoneIsMutated_TheStoredSnapshotIsUnchanged(t *testing.T) {
	t.Parallel()
	for caseName, fieldCase := range manualZoneFieldCases() {
		t.Run(caseName, func(t *testing.T) {
			t.Parallel()
			// Arrange
			state := newEditorState()
			state.SetManualEdits([]template_model.Zone{newPopulatedManualZone()}, nil)
			expected := fieldCase.read(state.GetCurrentState().ManualZones[0])

			// Act
			fieldCase.mutate(state.GetManualZones()[0])

			// Assert
			assert.Equal(t, expected, fieldCase.read(state.GetCurrentState().ManualZones[0]))
		})
	}
}
