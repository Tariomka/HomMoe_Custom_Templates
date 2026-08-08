package pickerHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenVisibleRowsAreRequested_ReturnsTheServiceRowModel(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newPickerHandlerFixture()
	entries := []dtos.PickerEntryDto{{ID: gofakeit.UUID()}}
	filter := gofakeit.Word()
	expected := []dtos.PickerRowDto{{IsGroupHeader: true, Group: gofakeit.Word()}}
	fixture.service.On("GetVisiblePickerRows", entries, filter, true).Return(expected)

	// Act
	rows := fixture.handler.GetVisiblePickerRows(entries, filter, true)

	// Assert
	assert.Equal(t, expected, rows)
}
