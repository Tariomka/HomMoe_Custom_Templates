package editorState_test

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWhenStateIsConstructed_CurrentStateIsDefault(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := dtos.NewDefaultEditorStateDto()

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
