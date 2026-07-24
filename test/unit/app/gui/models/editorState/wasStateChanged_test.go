package editorState_test

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWhenNoSnapshotExists_ReportsStateNotChanged(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.UpdateCurrentState(func(dto *dtos.EditorStateDto) { dto.PlayerCount++ })

	// Act
	wasChanged := state.WasStateChanged()

	// Assert
	assert.False(t, wasChanged)
}

func TestWhenStateDivergedFromSnapshot_ReportsStateChanged(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.SnapshotCurrentState()
	state.UpdateCurrentState(func(dto *dtos.EditorStateDto) { dto.PlayerCount++ })

	// Act
	wasChanged := state.WasStateChanged()

	// Assert
	assert.True(t, wasChanged)
}

func TestWhenStateEqualsSnapshot_ReportsStateNotChanged(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.SnapshotCurrentState()

	// Act
	wasChanged := state.WasStateChanged()

	// Assert
	assert.False(t, wasChanged)
}

func TestWhenOnlyManualEditsDifferFromSnapshot_ReportsStateNotChanged(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.SnapshotCurrentState()
	state.SetManualEdits(
		[]entities.Zone{{Name: "Zone A"}},
		[]entities.Connection{{Name: "A-B"}})

	// Act
	wasChanged := state.WasStateChanged()

	// Assert
	assert.False(t, wasChanged)
}
