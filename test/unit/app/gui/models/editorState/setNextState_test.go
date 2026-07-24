package editorState_test

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWhenNextStateIsAssigned_NextStateExists(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()

	// Act
	state.SetNextState(state.GetCurrentState())

	// Assert
	assert.True(t, state.HasNextState())
}

func TestWhenAssignedNextStateDiffersFromCurrent_PendingChangesAreReported(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	divergent := dtos.NewDefaultEditorStateDto()
	divergent.PlayerCount++

	// Act
	state.SetNextState(divergent)

	// Assert
	assert.True(t, state.HasPendingChanges())
}
