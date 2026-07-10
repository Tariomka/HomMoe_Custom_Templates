package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/models"
	"github.com/stretchr/testify/assert"
)

func TestWhenNextStateWasNeverSet_ReportsNoNextState(t *testing.T) {
	t.Parallel()
	// Arrange
	state := models.NewEditorState()

	// Act
	hasNext := state.HasNextState()

	// Assert
	assert.False(t, hasNext)
}

func TestWhenNextStateWasSet_ReportsNextState(t *testing.T) {
	t.Parallel()
	// Arrange
	state := models.NewEditorState()
	state.SetNextState(state.GetCurrentState())

	// Act
	hasNext := state.HasNextState()

	// Assert
	assert.True(t, hasNext)
}
