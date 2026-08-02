package editorState_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
