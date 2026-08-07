package state_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenStatusIsSet_GetStatusReturnsMessageAndSeverity(t *testing.T) {
	t.Parallel()
	// Arrange
	state := drivers.NewUIState(
		&test_helpers.TemplateHandlerMock{},
		test_helpers.NewFileSystemHandler(),
		test_helpers.NewRegenerationHandler(),
		false)
	expectedMessage := gofakeit.Sentence(5)
	expectedIsError := gofakeit.Bool()

	// Act
	state.SetStatus(expectedMessage, expectedIsError)

	// Assert
	message, isError := state.GetStatus()
	assert.Equal(t, []any{expectedMessage, expectedIsError}, []any{message, isError})
}
