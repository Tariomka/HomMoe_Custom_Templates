package pickerHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenItemEntriesAreBuilt_ReturnsTheServiceEntries(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newPickerHandlerFixture()
	items := []dtos.PickerItemDto{{Sid: gofakeit.UUID()}}
	expected := []dtos.PickerEntryDto{{ID: gofakeit.UUID()}}
	fixture.service.On("BuildItemPickerEntries", items).Return(expected)

	// Act
	entries := fixture.handler.BuildItemPickerEntries(items)

	// Assert
	assert.Equal(t, expected, entries)
}
