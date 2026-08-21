package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/stretchr/testify/assert"
)

func TestWhenStateIsConstructed_CurrentStateIsDefault(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := editor_state_dto.NewDefaultEditorStateDto()

	// Act
	state := newEditorState()

	// Assert
	assert.Equal(t, expected, state.GetCurrentState())
}

func TestWhenStateIsConstructed_ThereIsNoPreviousState(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	state := newEditorState()

	// Assert
	assert.False(t, state.HasPreviousState())
}

func TestWhenStateIsConstructed_ThereIsNoNextState(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	state := newEditorState()

	// Assert
	assert.False(t, state.HasNextState())
}
