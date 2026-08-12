package pickerHandler_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenTheFilterIsNormalized_ReturnsTheServiceFilter(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newPickerHandlerFixture()
	text := gofakeit.Word()
	expected := gofakeit.Word()
	fixture.service.On("NormalizePickerFilter", text).Return(expected)

	// Act
	filter := fixture.handler.NormalizePickerFilter(text)

	// Assert
	assert.Equal(t, expected, filter)
}
