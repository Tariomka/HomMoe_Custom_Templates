package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/stretchr/testify/assert"
)

func TestWhenNoNextStateExistsYet_ReportsCapture(t *testing.T) {
	// Arrange
	state := models.NewEditorState()

	// Act
	wasCaptured := state.SetNextFromCurrentIfStateIsBeingUpdated()

	// Assert
	assert.True(t, wasCaptured)
}

func TestWhenNoNextStateExistsYet_NextStateIsSet(t *testing.T) {
	// Arrange
	state := models.NewEditorState()

	// Act
	state.SetNextFromCurrentIfStateIsBeingUpdated()

	// Assert
	assert.True(t, state.HasNextState())
}

func TestWhenNextStateAlreadyMatchesCurrent_ReportsNoCapture(t *testing.T) {
	// Arrange
	state := models.NewEditorState()
	state.SetNextFromCurrentIfStateIsBeingUpdated()

	// Act
	wasCaptured := state.SetNextFromCurrentIfStateIsBeingUpdated()

	// Assert
	assert.False(t, wasCaptured)
}

func TestWhenCurrentStateMovedPastNextState_ReportsCapture(t *testing.T) {
	// Arrange
	state := models.NewEditorState()
	state.SetNextFromCurrentIfStateIsBeingUpdated()
	state.UpdateCurrentState(func(dto *dtos.EditorStateDto) { dto.PlayerCount++ })

	// Act
	wasCaptured := state.SetNextFromCurrentIfStateIsBeingUpdated()

	// Assert
	assert.True(t, wasCaptured)
}

func TestWhenCurrentStateMovedPastNextState_PendingChangesAreCleared(t *testing.T) {
	// Arrange
	state := models.NewEditorState()
	state.SetNextFromCurrentIfStateIsBeingUpdated()
	state.UpdateCurrentState(func(dto *dtos.EditorStateDto) { dto.PlayerCount++ })

	// Act
	state.SetNextFromCurrentIfStateIsBeingUpdated()

	// Assert
	assert.False(t, state.HasPendingChanges())
}
