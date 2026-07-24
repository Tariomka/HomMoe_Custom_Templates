package editorState_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWhenExistingSnapshotIsReset_PreviousStateIsGone(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.SnapshotCurrentState()
	require.True(t, state.HasPreviousState())

	// Act
	state.ResetPreviousState()

	// Assert
	assert.False(t, state.HasPreviousState())
}
