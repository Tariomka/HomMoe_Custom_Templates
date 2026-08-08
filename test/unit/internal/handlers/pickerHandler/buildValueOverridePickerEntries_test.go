package pickerHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenValueOverrideEntriesAreBuilt_ReturnsTheServiceEntries(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newPickerHandlerFixture()
	sids := []string{gofakeit.UUID()}
	expected := []dtos.PickerEntryDto{{ID: gofakeit.UUID()}}
	fixture.service.On("BuildValueOverridePickerEntries", sids).Return(expected)

	// Act
	entries := fixture.handler.BuildValueOverridePickerEntries(sids)

	// Assert
	assert.Equal(t, expected, entries)
}
