package editorState_test

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWhenNoSnapshotExists_DoesNotReportUnchangedState(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()

	// Act
	wasUnchanged := state.WasStateUnchanged()

	// Assert
	assert.False(t, wasUnchanged)
}

func TestWhenStateEqualsSnapshot_ReportsUnchangedState(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.SnapshotCurrentState()

	// Act
	wasUnchanged := state.WasStateUnchanged()

	// Assert
	assert.True(t, wasUnchanged)
}

func TestWhenStateDivergedFromSnapshot_DoesNotReportUnchangedState(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.SnapshotCurrentState()
	state.UpdateCurrentState(func(dto *dtos.EditorStateDto) { dto.PlayerCount++ })

	// Act
	wasUnchanged := state.WasStateUnchanged()

	// Assert
	assert.False(t, wasUnchanged)
}
