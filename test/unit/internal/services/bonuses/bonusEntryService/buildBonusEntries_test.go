package bonusEntryService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/bonuses"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenTheTownPortalPresetIsComposed_EmitsOneParameterlessEntry(t *testing.T) {
	t.Parallel()
	// Arrange
	service := bonuses.NewBonusEntryService()
	receiver := gofakeit.Word()

	// Act
	result := service.BuildBonusEntries(dtos.BonusCompositionRequestDto{
		PresetType:     config.BonusTownPortalFree,
		ReceiverFilter: receiver,
	})

	// Assert
	assert.Equal(t, []config.BonusEntry{{
		PresetType:     config.BonusTownPortalFree,
		ReceiverFilter: receiver,
	}}, result.Entries)
}

func TestWhenNoSpellIsSelected_ReportsTheMissingSpellError(t *testing.T) {
	t.Parallel()
	// Arrange
	service := bonuses.NewBonusEntryService()

	// Act
	result := service.BuildBonusEntries(dtos.BonusCompositionRequestDto{PresetType: config.BonusSpell})

	// Assert
	assert.Equal(t, "Pick at least one spell.", result.Error)
}

func TestWhenNoSpellIsSelected_EmitsNoEntries(t *testing.T) {
	t.Parallel()
	// Arrange
	service := bonuses.NewBonusEntryService()

	// Act
	result := service.BuildBonusEntries(dtos.BonusCompositionRequestDto{PresetType: config.BonusSpell})

	// Assert
	assert.Empty(t, result.Entries)
}

func TestWhenSelectedSpellsAreMadeFree_EachEntryCarriesTheFreeFlag(t *testing.T) {
	t.Parallel()
	// Arrange
	service := bonuses.NewBonusEntryService()
	receiver := gofakeit.Word()
	first := gofakeit.UUID()
	second := gofakeit.UUID()

	// Act
	result := service.BuildBonusEntries(dtos.BonusCompositionRequestDto{
		PresetType:     config.BonusSpell,
		ReceiverFilter: receiver,
		SelectedSpells: []string{first, second},
		MakeSpellsFree: true,
	})

	// Assert
	assert.Equal(t, []config.BonusEntry{
		{PresetType: config.BonusSpell, ReceiverFilter: receiver, Param: first, Param2: "1"},
		{PresetType: config.BonusSpell, ReceiverFilter: receiver, Param: second, Param2: "1"},
	}, result.Entries)
}

func TestWhenSelectedSpellsAreNotFree_EachEntryCarriesTheNotFreeFlag(t *testing.T) {
	t.Parallel()
	// Arrange
	service := bonuses.NewBonusEntryService()
	receiver := gofakeit.Word()
	spell := gofakeit.UUID()

	// Act
	result := service.BuildBonusEntries(dtos.BonusCompositionRequestDto{
		PresetType:     config.BonusSpell,
		ReceiverFilter: receiver,
		SelectedSpells: []string{spell},
	})

	// Assert
	assert.Equal(t, []config.BonusEntry{
		{PresetType: config.BonusSpell, ReceiverFilter: receiver, Param: spell, Param2: "0"},
	}, result.Entries)
}

func TestWhenTheMultiplierIsNotNumeric_ReportsTheMultiplierError(t *testing.T) {
	t.Parallel()
	// Arrange
	service := bonuses.NewBonusEntryService()

	// Act
	result := service.BuildBonusEntries(dtos.BonusCompositionRequestDto{
		PresetType:     config.BonusUnitMultiplier,
		MultiplierText: gofakeit.Word(),
	})

	// Assert
	assert.Equal(t, "Enter a numeric multiplier.", result.Error)
}

func TestWhenTheMultiplierIsBlank_ReportsTheMultiplierError(t *testing.T) {
	t.Parallel()
	// Arrange
	service := bonuses.NewBonusEntryService()

	// Act
	result := service.BuildBonusEntries(dtos.BonusCompositionRequestDto{
		PresetType:     config.BonusUnitMultiplier,
		MultiplierText: "   ",
	})

	// Assert
	assert.Equal(t, "Enter a numeric multiplier.", result.Error)
}

func TestWhenTheMultiplierIsNumeric_EmitsATrimmedMultiplierEntry(t *testing.T) {
	t.Parallel()
	// Arrange
	service := bonuses.NewBonusEntryService()
	receiver := gofakeit.Word()

	// Act
	result := service.BuildBonusEntries(dtos.BonusCompositionRequestDto{
		PresetType:     config.BonusUnitMultiplier,
		ReceiverFilter: receiver,
		MultiplierText: " 3.5 ",
	})

	// Assert
	assert.Equal(t, []config.BonusEntry{{
		PresetType:     config.BonusUnitMultiplier,
		ReceiverFilter: receiver,
		Param:          "3.5",
	}}, result.Entries)
}

func TestWhenTheMovementValueIsNotNumeric_ReportsTheMovementError(t *testing.T) {
	t.Parallel()
	// Arrange
	service := bonuses.NewBonusEntryService()

	// Act
	result := service.BuildBonusEntries(dtos.BonusCompositionRequestDto{
		PresetType:   config.BonusMovementBonus,
		MovementText: gofakeit.Word(),
	})

	// Assert
	assert.Equal(t, "Enter a numeric movement value.", result.Error)
}

func TestWhenTheMovementValueIsNumeric_EmitsAMovementEntry(t *testing.T) {
	t.Parallel()
	// Arrange
	service := bonuses.NewBonusEntryService()
	receiver := gofakeit.Word()

	// Act
	result := service.BuildBonusEntries(dtos.BonusCompositionRequestDto{
		PresetType:     config.BonusMovementBonus,
		ReceiverFilter: receiver,
		MovementText:   "500",
	})

	// Assert
	assert.Equal(t, []config.BonusEntry{{
		PresetType:     config.BonusMovementBonus,
		ReceiverFilter: receiver,
		Param:          "500",
	}}, result.Entries)
}

func TestWhenNoStartingItemIsGiven_ReportsTheItemError(t *testing.T) {
	t.Parallel()
	// Arrange
	service := bonuses.NewBonusEntryService()

	// Act
	result := service.BuildBonusEntries(dtos.BonusCompositionRequestDto{
		PresetType: config.BonusStartingItem,
		ItemText:   "  ",
	})

	// Assert
	assert.Equal(t, "Pick or enter an item.", result.Error)
}

func TestWhenAStartingItemIsGiven_EmitsATrimmedItemEntry(t *testing.T) {
	t.Parallel()
	// Arrange
	service := bonuses.NewBonusEntryService()
	receiver := gofakeit.Word()
	item := gofakeit.UUID()

	// Act
	result := service.BuildBonusEntries(dtos.BonusCompositionRequestDto{
		PresetType:     config.BonusStartingItem,
		ReceiverFilter: receiver,
		ItemText:       " " + item + " ",
	})

	// Assert
	assert.Equal(t, []config.BonusEntry{{
		PresetType:     config.BonusStartingItem,
		ReceiverFilter: receiver,
		Param:          item,
	}}, result.Entries)
}

func TestWhenTheResourceAmountIsNotNumeric_ReportsTheAmountError(t *testing.T) {
	t.Parallel()
	// Arrange
	service := bonuses.NewBonusEntryService()

	// Act
	result := service.BuildBonusEntries(dtos.BonusCompositionRequestDto{
		PresetType:   config.BonusStartingGold,
		ResourceText: gofakeit.Word(),
	})

	// Assert
	assert.Equal(t, "Enter a numeric amount.", result.Error)
}

func TestWhenTheResourceAmountIsNumeric_EmitsAResourceEntry(t *testing.T) {
	t.Parallel()
	// Arrange
	service := bonuses.NewBonusEntryService()
	receiver := gofakeit.Word()

	// Act
	result := service.BuildBonusEntries(dtos.BonusCompositionRequestDto{
		PresetType:     config.BonusStartingGold,
		ReceiverFilter: receiver,
		ResourceText:   "10000",
	})

	// Assert
	assert.Equal(t, []config.BonusEntry{{
		PresetType:     config.BonusStartingGold,
		ReceiverFilter: receiver,
		Param:          "10000",
	}}, result.Entries)
}
