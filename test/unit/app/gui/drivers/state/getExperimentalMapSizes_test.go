package state_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenExperimentalMapSizesWasEnabled_GetExperimentalMapSizesReturnsTrue(t *testing.T) {
	t.Parallel()
	// Arrange
	state := drivers.NewUIState(
		&test_helpers.TemplateHandlerMock{},
		test_helpers.NewFileSystemHandler(),
		test_helpers.NewRegenerationHandler(),
		false)
	state.UpdateState(func(dto *dtos.EditorStateDto) { dto.ExperimentalMapSizes = true })

	// Act
	actual := state.GetExperimentalMapSizes()

	// Assert
	assert.True(t, actual)
}
