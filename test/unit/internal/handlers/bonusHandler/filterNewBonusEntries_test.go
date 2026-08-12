package bonusHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenBonusEntriesAreFiltered_ReturnsTheServiceSelection(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newBonusHandlerFixture()
	entries := []config.BonusEntry{{PresetType: config.BonusStartingGold, Param: gofakeit.DigitN(4)}}
	existingKeys := map[string]bool{gofakeit.UUID(): true}
	expected := []config.BonusEntry{{PresetType: config.BonusMovementBonus, Param: gofakeit.DigitN(3)}}
	fixture.bonusService.On("FilterNewBonusEntries", entries, existingKeys).Return(expected)

	// Act
	fresh := fixture.handler.FilterNewBonusEntries(entries, existingKeys)

	// Assert
	assert.Equal(t, expected, fresh)
}
