package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/stretchr/testify/assert"
)

func TestWhenExperimentalMapSizesWasEnabled_GetExperimentalMapSizesReturnsTrue(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.UpdateCurrentState(func(dto *editor_state_dto.EditorStateDto) { dto.ExperimentalMapSizes = true })

	// Act
	actual := state.GetExperimentalMapSizes()

	// Assert
	assert.True(t, actual)
}
