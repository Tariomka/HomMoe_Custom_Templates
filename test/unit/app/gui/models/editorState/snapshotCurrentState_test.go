package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/stretchr/testify/assert"
)

func TestWhenSnapshotIsTaken_PreviousStateExists(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()

	// Act
	state.SnapshotCurrentState()

	// Assert
	assert.True(t, state.HasPreviousState())
}

func TestWhenCurrentStateChangesAfterSnapshot_SnapshotKeepsOldValues(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.SnapshotCurrentState()

	// Act
	state.UpdateCurrentState(func(dto *dtos.EditorStateDto) { dto.PlayerCount++ })

	// Assert - the snapshot still holds the old player count, so the state reads as changed
	assert.True(t, state.WasStateChanged())
}
