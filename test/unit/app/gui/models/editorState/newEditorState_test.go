package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/stretchr/testify/assert"
)

func TestWhenStateIsConstructed_CurrentStateIsDefault(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := dtos.NewDefaultEditorStateDto()

	// Act
	state := models.NewEditorState()

	// Assert
	assert.Equal(t, expected, state.GetCurrentState())
}

func TestWhenStateIsConstructed_ThereIsNoPreviousState(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	state := models.NewEditorState()

	// Assert
	assert.False(t, state.HasPreviousState())
}

func TestWhenStateIsConstructed_ThereIsNoNextState(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	state := models.NewEditorState()

	// Assert
	assert.False(t, state.HasNextState())
}
