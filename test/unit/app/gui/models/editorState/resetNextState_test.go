package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenExistingNextStateIsReset_NextStateIsGone(t *testing.T) {
	// Arrange
	state := models.NewEditorState()
	state.SetNextState(state.GetCurrentState())
	require.True(t, state.HasNextState())

	// Act
	state.ResetNextState()

	// Assert
	assert.False(t, state.HasNextState())
}
