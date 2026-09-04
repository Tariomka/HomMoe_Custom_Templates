package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenZonesAreApplied_ManualZoneSavesAreStoredInCurrentState(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	position := &[2]float64{gofakeit.Float64Range(0, 1), gofakeit.Float64Range(0, 1)}
	zone := template_model.Zone{Name: "Zone A", Size: gofakeit.Float64Range(0.5, 2), ManualPosition: position}

	// Act
	state.SetManualEdits([]template_model.Zone{zone}, nil)

	// Assert
	assert.Equal(t,
		[]editor_state_model.ManualZoneSave{{
			Zone:           template_model.ToZoneEntity(zone),
			ManualPosition: position,
		}},
		state.GetCurrentState().ManualZones)
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
