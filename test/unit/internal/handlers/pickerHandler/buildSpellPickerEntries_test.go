package pickerHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenSpellEntriesAreBuilt_ReturnsTheServiceEntries(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newPickerHandlerFixture()
	spells := []dtos.PickerSpellDto{{Sid: gofakeit.UUID()}}
	expected := []dtos.PickerEntryDto{{ID: gofakeit.UUID()}}
	fixture.service.On("BuildSpellPickerEntries", spells).Return(expected)

	// Act
	entries := fixture.handler.BuildSpellPickerEntries(spells)

	// Assert
	assert.Equal(t, expected, entries)
}
