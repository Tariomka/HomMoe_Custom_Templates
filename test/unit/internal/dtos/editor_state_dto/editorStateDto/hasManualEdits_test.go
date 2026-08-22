package editorStateDto_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/stretchr/testify/assert"
)

func TestWhenNoManualDataIsPresent_ReportsNoManualEdits(t *testing.T) {
	t.Parallel()
	// Arrange
	state := editor_state_dto.NewDefaultEditorStateDto()

	// Act
	hasEdits := state.HasManualEdits()

	// Assert
	assert.False(t, hasEdits)
}

func TestWhenOnlyManualZonesArePresent_ReportsManualEdits(t *testing.T) {
	t.Parallel()
	// Arrange
	state := editor_state_dto.NewDefaultEditorStateDto()
	state.ManualZones = []editor_state.ManualZoneSave{{Zone: entities.Zone{Name: "Zone A"}}}

	// Act
	hasEdits := state.HasManualEdits()

	// Assert
	assert.True(t, hasEdits)
}

func TestWhenOnlyManualConnectionsArePresent_ReportsManualEdits(t *testing.T) {
	t.Parallel()
	// Arrange
	state := editor_state_dto.NewDefaultEditorStateDto()
	state.ManualConnections = []editor_state.ManualConnectionSave{
		{Connection: entities.Connection{Name: "A-B"}},
	}

	// Act
	hasEdits := state.HasManualEdits()

	// Assert
	assert.True(t, hasEdits)
}
