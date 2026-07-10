package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenStoredManualEditsAreCleared_NoManualEditsRemain(t *testing.T) {
	t.Parallel()
	// Arrange
	state := models.NewEditorState()
	state.SetManualEdits(
		[]entities.Zone{{Name: "Zone A"}},
		[]entities.Connection{{Name: "A-B"}})
	require.True(t, state.HasManualEdits())

	// Act
	state.ClearManualEdits()

	// Assert
	assert.False(t, state.HasManualEdits())
}
