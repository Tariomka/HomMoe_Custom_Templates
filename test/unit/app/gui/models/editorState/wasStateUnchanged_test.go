package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/stretchr/testify/assert"
)

func TestWhenNoSnapshotExists_DoesNotReportUnchangedState(t *testing.T) {
	// Arrange
	state := models.NewEditorState()

	// Act
	wasUnchanged := state.WasStateUnchanged()

	// Assert
	assert.False(t, wasUnchanged)
}

func TestWhenStateEqualsSnapshot_ReportsUnchangedState(t *testing.T) {
	// Arrange
	state := models.NewEditorState()
	state.SnapshotCurrentState()

	// Act
	wasUnchanged := state.WasStateUnchanged()

	// Assert
	assert.True(t, wasUnchanged)
}

func TestWhenStateDivergedFromSnapshot_DoesNotReportUnchangedState(t *testing.T) {
	// Arrange
	state := models.NewEditorState()
	state.SnapshotCurrentState()
	state.UpdateCurrentState(func(dto *dtos.EditorStateDto) { dto.PlayerCount++ })

	// Act
	wasUnchanged := state.WasStateUnchanged()

	// Assert
	assert.False(t, wasUnchanged)
}
