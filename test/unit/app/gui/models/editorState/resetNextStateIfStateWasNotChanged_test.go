package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/stretchr/testify/assert"
)

func unchangedStateWithNext() *models.EditorState {
	state := newEditorState()
	state.SnapshotCurrentState()
	state.SetNextState(state.GetCurrentState())
	return state
}

func changedStateWithNext() *models.EditorState {
	state := unchangedStateWithNext()
	state.UpdateCurrentState(func(dto *dtos.EditorStateDto) { dto.PlayerCount++ })
	return state
}

func TestWhenStateStillMatchesSnapshot_ReportsReset(t *testing.T) {
	t.Parallel()
	// Arrange
	state := unchangedStateWithNext()

	// Act
	wasReset := state.ResetNextStateIfStateWasNotChanged()

	// Assert
	assert.True(t, wasReset)
}

func TestWhenStateStillMatchesSnapshot_NextStateIsCleared(t *testing.T) {
	t.Parallel()
	// Arrange
	state := unchangedStateWithNext()

	// Act
	state.ResetNextStateIfStateWasNotChanged()

	// Assert
	assert.False(t, state.HasNextState())
}

func TestWhenStateDiffersFromSnapshot_ReportsNoReset(t *testing.T) {
	t.Parallel()
	// Arrange
	state := changedStateWithNext()

	// Act
	wasReset := state.ResetNextStateIfStateWasNotChanged()

	// Assert
	assert.False(t, wasReset)
}

func TestWhenStateDiffersFromSnapshot_NextStateIsKept(t *testing.T) {
	t.Parallel()
	// Arrange
	state := changedStateWithNext()

	// Act
	state.ResetNextStateIfStateWasNotChanged()

	// Assert
	assert.True(t, state.HasNextState())
}

func TestWhenNoSnapshotExistsYet_ReportsNoReset(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.SetNextState(state.GetCurrentState())

	// Act
	wasReset := state.ResetNextStateIfStateWasNotChanged()

	// Assert
	assert.False(t, wasReset)
}
