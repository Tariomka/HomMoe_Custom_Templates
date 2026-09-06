package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenNoManualDataIsPresent_ReportsNoManualEdits(t *testing.T) {
	t.Parallel()
	// Arrange
	state := editor_state_model.NewDefaultEditorStateModel()

	// Act
	hasEdits := state.HasManualEdits()

	// Assert
	assert.False(t, hasEdits)
}

func TestWhenOnlyManualZonesArePresent_ReportsManualEdits(t *testing.T) {
	t.Parallel()
	// Arrange
	state := editor_state_model.NewDefaultEditorStateModel()
	state.ManualZones = []template_model.Zone{{Name: "Zone A"}}

	// Act
	hasEdits := state.HasManualEdits()

	// Assert
	assert.True(t, hasEdits)
}

func TestWhenOnlyManualConnectionsArePresent_ReportsManualEdits(t *testing.T) {
	t.Parallel()
	// Arrange
	state := editor_state_model.NewDefaultEditorStateModel()
	state.ManualConnections = []editor_state_model.ManualConnectionSave{
		{Connection: template_entity.Connection{Name: "A-B"}},
	}

	// Act
	hasEdits := state.HasManualEdits()

	// Assert
	assert.True(t, hasEdits)
}
