package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenZonesAreApplied_ManualZonesAreStoredInCurrentState(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	position := new(data.NewVec2(gofakeit.Float64Range(0, 1), gofakeit.Float64Range(0, 1)))
	zone := template_model.Zone{Name: "Zone A", Size: gofakeit.Float64Range(0.5, 2), ManualPosition: position}

	// Act
	state.SetManualEdits([]template_model.Zone{zone}, nil)

	// Assert
	assert.Equal(t, []template_model.Zone{zone}, state.GetCurrentState().ManualZones)
}

func TestWhenConnectionsAreApplied_ManualConnectionSavesAreStoredInCurrentState(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	connection := template_model.Connection{Name: "A-B", From: "Zone A", To: "Zone B", IsUserAdded: true}

	// Act
	state.SetManualEdits(nil, []template_model.Connection{connection})

	// Assert
	assert.Equal(t,
		[]editor_state_model.ManualConnectionSave{
			{Connection: template_model.ToConnectionEntity(connection), IsUserAdded: true}},
		state.GetCurrentState().ManualConnections)
}

// The editor keeps editing its own zones after applying; the committed
// snapshot must not follow those later mutations.
func TestWhenAnInputZoneIsMutatedAfterStoring_TheStoredSnapshotIsUnchanged(t *testing.T) {
	t.Parallel()
	for caseName, fieldCase := range manualZoneFieldCases() {
		t.Run(caseName, func(t *testing.T) {
			t.Parallel()
			// Arrange
			state := newEditorState()
			zone := newPopulatedManualZone()
			state.SetManualEdits([]template_model.Zone{zone}, nil)
			expected := fieldCase.read(zone)

			// Act
			fieldCase.mutate(zone)

			// Assert
			assert.Equal(t, expected, fieldCase.read(state.GetCurrentState().ManualZones[0]))
		})
	}
}

// A first snapshot is a change even when it merely pins the layout the
// generator already produced: the file did not carry it before.
func TestWhenTheFirstSnapshotIsStored_TheChangeIsReported(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()

	// Act
	changed := state.SetManualEdits([]template_model.Zone{newPopulatedManualZone()}, nil)

	// Assert
	assert.True(t, changed)
}

func TestWhenTheSameSnapshotIsStoredAgain_NoChangeIsReported(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.SetManualEdits([]template_model.Zone{newPopulatedManualZone()}, nil)

	// Act
	changed := state.SetManualEdits([]template_model.Zone{newPopulatedManualZone()}, nil)

	// Assert
	assert.False(t, changed)
}

func TestWhenOnlyAZoneMoved_TheChangeIsReported(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.SetManualEdits([]template_model.Zone{newPopulatedManualZone()}, nil)
	movedZone := newPopulatedManualZone()
	movedZone.ManualPosition = new(data.NewVec2(9.0, 9.0))

	// Act
	changed := state.SetManualEdits([]template_model.Zone{movedZone}, nil)

	// Assert
	assert.True(t, changed)
}

// Length is written by the connection converter, so it is part of what the
// file holds and part of what makes an apply dirty.
func TestWhenOnlyAConnectionLengthChanged_TheChangeIsReported(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	connection := template_model.Connection{Name: "A-B", From: "Zone A", To: "Zone B", Length: 2}
	state.SetManualEdits(nil, []template_model.Connection{connection})
	connection.Length = 3

	// Act
	changed := state.SetManualEdits(nil, []template_model.Connection{connection})

	// Assert
	assert.True(t, changed)
}

// IsUserAdded has nowhere to live in the template schema, which is exactly why
// the snapshot carries it.
func TestWhenOnlyIsUserAddedChanged_TheChangeIsReported(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	connection := template_model.Connection{Name: "A-B", From: "Zone A", To: "Zone B"}
	state.SetManualEdits(nil, []template_model.Connection{connection})
	connection.IsUserAdded = true

	// Act
	changed := state.SetManualEdits(nil, []template_model.Connection{connection})

	// Assert
	assert.True(t, changed)
}

func TestWhenAnEmptyLayoutIsStoredOverNothing_NoChangeIsReported(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()

	// Act
	changed := state.SetManualEdits([]template_model.Zone{}, []template_model.Connection{})

	// Assert
	assert.False(t, changed)
}
