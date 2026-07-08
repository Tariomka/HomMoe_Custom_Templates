package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/stretchr/testify/assert"
)

func TestWhenSnapshotIsTaken_PreviousStateExists(t *testing.T) {
	// Arrange
	state := models.NewEditorState()

	// Act
	state.SnapshotCurrentState()

	// Assert
	assert.True(t, state.HasPreviousState())
}

func TestWhenCurrentStateChangesAfterSnapshot_SnapshotKeepsOldValues(t *testing.T) {
	// Arrange
	state := models.NewEditorState()
	state.SnapshotCurrentState()

	// Act
	state.UpdateCurrentState(func(dto *dtos.EditorStateDto) { dto.PlayerCount++ })

	// Assert - the snapshot still holds the old player count, so the state reads as changed
	assert.True(t, state.WasStateChanged())
}
