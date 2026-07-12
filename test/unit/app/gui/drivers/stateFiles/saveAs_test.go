package stateFiles_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenSaveAsIsCalled_DialogIsOpened(t *testing.T) {
	t.Parallel()
	// Arrange
	state := drivers.NewUIStateWithHandler(&test_helpers.TemplateHandlerMock{})

	// Act
	state.SaveAs(gofakeit.ProductName())

	// Assert
	assert.True(t, state.GetDialogHost().IsOpen())
}
