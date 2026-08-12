package pickerHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenTheSelectionIsRequested_ReturnsTheServiceIds(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newPickerHandlerFixture()
	entries := []dtos.PickerEntryDto{{ID: gofakeit.UUID()}}
	selected := map[string]bool{gofakeit.UUID(): true}
	expected := []string{gofakeit.UUID()}
	fixture.service.On("GetSelectedPickerIDs", entries, selected).Return(expected)

	// Act
	ids := fixture.handler.GetSelectedPickerIDs(entries, selected)

	// Assert
	assert.Equal(t, expected, ids)
}
