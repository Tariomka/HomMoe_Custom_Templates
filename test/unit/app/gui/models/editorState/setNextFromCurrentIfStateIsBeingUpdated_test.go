package editorState_test

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWhenNoNextStateExistsYet_ReportsCapture(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()

	// Act
	wasCaptured := state.SetNextFromCurrentIfStateIsBeingUpdated()

	// Assert
	assert.True(t, wasCaptured)
}

func TestWhenNoNextStateExistsYet_NextStateIsSet(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()

	// Act
	state.SetNextFromCurrentIfStateIsBeingUpdated()

	// Assert
	assert.True(t, state.HasNextState())
}

func TestWhenNextStateAlreadyMatchesCurrent_ReportsNoCapture(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.SetNextFromCurrentIfStateIsBeingUpdated()

	// Act
	wasCaptured := state.SetNextFromCurrentIfStateIsBeingUpdated()

	// Assert
	assert.False(t, wasCaptured)
}

func TestWhenCurrentStateMovedPastNextState_ReportsCapture(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.SetNextFromCurrentIfStateIsBeingUpdated()
	state.UpdateCurrentState(func(dto *dtos.EditorStateDto) { dto.PlayerCount++ })

	// Act
	wasCaptured := state.SetNextFromCurrentIfStateIsBeingUpdated()

	// Assert
	assert.True(t, wasCaptured)
}

func TestWhenCurrentStateMovedPastNextState_PendingChangesAreCleared(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.SetNextFromCurrentIfStateIsBeingUpdated()
	state.UpdateCurrentState(func(dto *dtos.EditorStateDto) { dto.PlayerCount++ })

	// Act
	state.SetNextFromCurrentIfStateIsBeingUpdated()

	// Assert
	assert.False(t, state.HasPendingChanges())
}
