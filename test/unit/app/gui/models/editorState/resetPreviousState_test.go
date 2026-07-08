package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenExistingSnapshotIsReset_PreviousStateIsGone(t *testing.T) {
	// Arrange
	state := models.NewEditorState()
	state.SnapshotCurrentState()
	require.True(t, state.HasPreviousState())

	// Act
	state.ResetPreviousState()

	// Assert
	assert.False(t, state.HasPreviousState())
}
