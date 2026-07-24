package editorState_test

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWhenStoredManualEditsAreCleared_NoManualEditsRemain(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.SetManualEdits(
		[]entities.Zone{{Name: "Zone A"}},
		[]entities.Connection{{Name: "A-B"}})
	require.True(t, state.HasManualEdits())

	// Act
	state.ClearManualEdits()

	// Assert
	assert.False(t, state.HasManualEdits())
}
