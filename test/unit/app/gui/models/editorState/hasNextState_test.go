package editorState_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWhenNextStateWasNeverSet_ReportsNoNextState(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()

	// Act
	hasNext := state.HasNextState()

	// Assert
	assert.False(t, hasNext)
}

func TestWhenNextStateWasSet_ReportsNextState(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.SetNextState(state.GetCurrentState())

	// Act
	hasNext := state.HasNextState()

	// Assert
	assert.True(t, hasNext)
}
