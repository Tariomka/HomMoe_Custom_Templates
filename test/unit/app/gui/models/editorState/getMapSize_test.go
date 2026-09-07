package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenMapSizeWasUpdated_GetMapSizeReturnsIt(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	mapSize := gofakeit.Number(100, 300)
	state.UpdateCurrentState(func(dto *editor_state_model.EditorState) { dto.MapSize = mapSize })

	// Act
	actual := state.GetMapSize()

	// Assert
	assert.Equal(t, mapSize, actual)
}
