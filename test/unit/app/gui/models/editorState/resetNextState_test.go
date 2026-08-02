package editorState_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenExistingNextStateIsReset_NextStateIsGone(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.SetNextState(state.GetCurrentState())
	require.True(t, state.HasNextState())

	// Act
	state.ResetNextState()

	// Assert
	assert.False(t, state.HasNextState())
}
