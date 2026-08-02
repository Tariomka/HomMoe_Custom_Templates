package editorState_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWhenNoSnapshotWasTaken_ReportsNoPreviousState(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()

	// Act
	hasPrevious := state.HasPreviousState()

	// Assert
	assert.False(t, hasPrevious)
}

func TestWhenSnapshotWasTaken_ReportsPreviousState(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.SnapshotCurrentState()

	// Act
	hasPrevious := state.HasPreviousState()

	// Assert
	assert.True(t, hasPrevious)
}
