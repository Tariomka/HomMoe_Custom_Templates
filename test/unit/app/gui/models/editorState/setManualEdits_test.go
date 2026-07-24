package editorState_test

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWhenZonesAreApplied_ManualZoneSavesAreStoredInCurrentState(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	position := &[2]float64{gofakeit.Float64Range(0, 1), gofakeit.Float64Range(0, 1)}
	zone := entities.Zone{Name: "Zone A", Size: gofakeit.Float64Range(0.5, 2), ManualPosition: position}

	// Act
	state.SetManualEdits([]entities.Zone{zone}, nil)

	// Assert
	assert.Equal(t,
		[]editor_state_dto.ManualZoneSave{{Zone: zone, ManualPosition: position}},
		state.GetCurrentState().ManualZones)
}

func TestWhenConnectionsAreApplied_ManualConnectionSavesAreStoredInCurrentState(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	connection := entities.Connection{Name: "A-B", From: "Zone A", To: "Zone B", IsUserAdded: true}

	// Act
	state.SetManualEdits(nil, []entities.Connection{connection})

	// Assert
	assert.Equal(t,
		[]editor_state_dto.ManualConnectionSave{{Connection: connection, IsUserAdded: true}},
		state.GetCurrentState().ManualConnections)
}
