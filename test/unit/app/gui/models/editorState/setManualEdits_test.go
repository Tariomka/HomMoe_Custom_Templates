package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenZonesAreApplied_ManualZoneSavesAreStoredInCurrentState(t *testing.T) {
	// Arrange
	state := models.NewEditorState()
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
	// Arrange
	state := models.NewEditorState()
	connection := entities.Connection{Name: "A-B", From: "Zone A", To: "Zone B", IsUserAdded: true}

	// Act
	state.SetManualEdits(nil, []entities.Connection{connection})

	// Assert
	assert.Equal(t,
		[]editor_state_dto.ManualConnectionSave{{Connection: connection, IsUserAdded: true}},
		state.GetCurrentState().ManualConnections)
}
