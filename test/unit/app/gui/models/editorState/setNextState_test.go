package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
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

func TestWhenAssignedNextStateDiffersFromCurrent_AssignedStateIsStored(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	divergent := editor_state_model.NewDefaultEditorStateModel()
	divergent.PlayerCount++

	// Act
	state.SetNextState(divergent)

	// Assert
	assert.Equal(t, &divergent, state.GetNextState())
}
