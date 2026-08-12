package bonusEntryService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/bonuses"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenAnEntryIsAlreadyKnown_ItIsDroppedFromTheResult(t *testing.T) {
	t.Parallel()
	// Arrange
	service := bonuses.NewBonusEntryService()
	known := config.BonusEntry{PresetType: config.BonusStartingGold, Param: gofakeit.DigitN(3)}

	// Act
	fresh := service.FilterNewBonusEntries(
		[]config.BonusEntry{known},
		map[string]bool{known.GetHash(): true})

	// Assert
	assert.Empty(t, fresh)
}

func TestWhenAnEntryIsUnknown_ItIsKeptInTheResult(t *testing.T) {
	t.Parallel()
	// Arrange
	service := bonuses.NewBonusEntryService()
	known := config.BonusEntry{PresetType: config.BonusStartingGold, Param: gofakeit.DigitN(3)}
	unknown := config.BonusEntry{PresetType: config.BonusMovementBonus, Param: gofakeit.DigitN(3)}

	// Act
	fresh := service.FilterNewBonusEntries(
		[]config.BonusEntry{known, unknown},
		map[string]bool{known.GetHash(): true})

	// Assert
	assert.Equal(t, []config.BonusEntry{unknown}, fresh)
}

func TestWhenNothingIsKnown_EveryEntryIsKept(t *testing.T) {
	t.Parallel()
	// Arrange
	service := bonuses.NewBonusEntryService()
	entries := []config.BonusEntry{
		{PresetType: config.BonusStartingGold, Param: gofakeit.DigitN(3)},
		{PresetType: config.BonusMovementBonus, Param: gofakeit.DigitN(3)},
	}

	// Act
	fresh := service.FilterNewBonusEntries(entries, nil)

	// Assert
	assert.Equal(t, entries, fresh)
}
