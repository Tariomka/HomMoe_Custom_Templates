package state_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenMapSizeWasUpdated_GetMapSizeReturnsIt(t *testing.T) {
	t.Parallel()
	// Arrange
	state := drivers.NewUIState(
		&test_helpers.TemplateHandlerMock{},
		test_helpers.NewFileSystemHandler(),
		test_helpers.NewRegenerationHandler(),
		false)
	mapSize := gofakeit.Number(100, 300)
	state.UpdateState(func(dto *dtos.EditorStateDto) { dto.MapSize = mapSize })

	// Act
	actual := state.GetMapSize()

	// Assert
	assert.Equal(t, mapSize, actual)
}
