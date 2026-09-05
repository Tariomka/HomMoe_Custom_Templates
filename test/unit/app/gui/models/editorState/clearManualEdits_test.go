package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenStoredManualEditsAreCleared_NoManualEditsRemain(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.SetManualEdits(
		[]template_model.Zone{{Name: "Zone A"}},
		[]template_model.Connection{{Name: "A-B"}})
	require.True(t, state.HasManualEdits())

	// Act
	state.ClearManualEdits()

	// Assert
	assert.False(t, state.HasManualEdits())
}
