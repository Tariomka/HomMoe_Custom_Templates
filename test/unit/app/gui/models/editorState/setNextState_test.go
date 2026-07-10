package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/stretchr/testify/assert"
)

func TestWhenNextStateIsAssigned_NextStateExists(t *testing.T) {
	t.Parallel()
	// Arrange
	state := models.NewEditorState()

	// Act
	state.SetNextState(state.GetCurrentState())

	// Assert
	assert.True(t, state.HasNextState())
}

func TestWhenAssignedNextStateDiffersFromCurrent_PendingChangesAreReported(t *testing.T) {
	t.Parallel()
	// Arrange
	state := models.NewEditorState()
	divergent := dtos.NewDefaultEditorStateDto()
	divergent.PlayerCount++

	// Act
	state.SetNextState(divergent)

	// Assert
	assert.True(t, state.HasPendingChanges())
}
