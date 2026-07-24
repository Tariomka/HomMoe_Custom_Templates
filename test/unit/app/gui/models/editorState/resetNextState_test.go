package editorState_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
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
