package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenMapSizeWasUpdated_GetMapSizeReturnsIt(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	mapSize := gofakeit.Number(100, 300)
	state.UpdateCurrentState(func(dto *dtos.EditorStateDto) { dto.MapSize = mapSize })

	// Act
	actual := state.GetMapSize()

	// Assert
	assert.Equal(t, mapSize, actual)
}
