package editorState_test

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWhenNoNextStateExists_ReportsNoPendingChanges(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()

	// Act
	hasPending := state.HasPendingChanges()

	// Assert
	assert.False(t, hasPending)
}

func TestWhenNextStateEqualsCurrent_ReportsNoPendingChanges(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.SetNextState(state.GetCurrentState())

	// Act
	hasPending := state.HasPendingChanges()

	// Assert
	assert.False(t, hasPending)
}

func TestWhenNextStateDiffersFromCurrent_ReportsPendingChanges(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.SetNextState(state.GetCurrentState())
	state.UpdateCurrentState(func(dto *dtos.EditorStateDto) { dto.PlayerCount++ })

	// Act
	hasPending := state.HasPendingChanges()

	// Assert
	assert.True(t, hasPending)
}
