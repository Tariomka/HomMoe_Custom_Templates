package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/stretchr/testify/assert"
)

func TestWhenNoNextStateExists_ReportsNoPendingChanges(t *testing.T) {
	t.Parallel()
	// Arrange
	state := models.NewEditorState()

	// Act
	hasPending := state.HasPendingChanges()

	// Assert
	assert.False(t, hasPending)
}

func TestWhenNextStateEqualsCurrent_ReportsNoPendingChanges(t *testing.T) {
	t.Parallel()
	// Arrange
	state := models.NewEditorState()
	state.SetNextState(state.GetCurrentState())

	// Act
	hasPending := state.HasPendingChanges()

	// Assert
	assert.False(t, hasPending)
}

func TestWhenNextStateDiffersFromCurrent_ReportsPendingChanges(t *testing.T) {
	t.Parallel()
	// Arrange
	state := models.NewEditorState()
	state.SetNextState(state.GetCurrentState())
	state.UpdateCurrentState(func(dto *dtos.EditorStateDto) { dto.PlayerCount++ })

	// Act
	hasPending := state.HasPendingChanges()

	// Assert
	assert.True(t, hasPending)
}
