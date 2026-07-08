package editorStateDto_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/stretchr/testify/assert"
)

func TestWhenNoManualDataIsPresent_ReportsNoManualEdits(t *testing.T) {
	// Arrange
	state := dtos.NewDefaultEditorStateDto()

	// Act
	hasEdits := state.HasManualEdits()

	// Assert
	assert.False(t, hasEdits)
}

func TestWhenOnlyManualZonesArePresent_ReportsManualEdits(t *testing.T) {
	// Arrange
	state := dtos.NewDefaultEditorStateDto()
	state.ManualZones = []editor_state_dto.ManualZoneSave{{Zone: entities.Zone{Name: "Zone A"}}}

	// Act
	hasEdits := state.HasManualEdits()

	// Assert
	assert.True(t, hasEdits)
}

func TestWhenOnlyManualConnectionsArePresent_ReportsManualEdits(t *testing.T) {
	// Arrange
	state := dtos.NewDefaultEditorStateDto()
	state.ManualConnections = []editor_state_dto.ManualConnectionSave{
		{Connection: entities.Connection{Name: "A-B"}},
	}

	// Act
	hasEdits := state.HasManualEdits()

	// Assert
	assert.True(t, hasEdits)
}
