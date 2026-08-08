package bonusEntryService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/bonuses"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenNoBonusesExist_TheKeySetIsEmpty(t *testing.T) {
	t.Parallel()
	// Arrange
	service := bonuses.NewBonusEntryService()

	// Act
	summary := service.DescribeExistingBonuses(nil)

	// Assert
	assert.Empty(t, summary.Keys)
}

func TestWhenNoBonusesExist_TheSpellListIsEmpty(t *testing.T) {
	t.Parallel()
	// Arrange
	service := bonuses.NewBonusEntryService()

	// Act
	summary := service.DescribeExistingBonuses(nil)

	// Assert
	assert.Empty(t, summary.SpellIDs)
}

func TestWhenABonusExists_ItsHashIsMarkedAsTaken(t *testing.T) {
	t.Parallel()
	// Arrange
	service := bonuses.NewBonusEntryService()
	entry := config.BonusEntry{
		PresetType:     config.BonusStartingGold,
		ReceiverFilter: gofakeit.Word(),
		Param:          gofakeit.DigitN(4),
	}

	// Act
	summary := service.DescribeExistingBonuses([]config.BonusEntry{entry})

	// Assert
	assert.Equal(t, map[string]bool{entry.GetHash(): true}, summary.Keys)
}

func TestWhenASpellBonusExists_ItsSpellIdIsCollected(t *testing.T) {
	t.Parallel()
	// Arrange
	service := bonuses.NewBonusEntryService()
	spellID := gofakeit.UUID()

	// Act
	summary := service.DescribeExistingBonuses([]config.BonusEntry{{
		PresetType: config.BonusSpell,
		Param:      spellID,
	}})

	// Assert
	assert.Equal(t, []string{spellID}, summary.SpellIDs)
}

func TestWhenASpellBonusHasNoSpellId_NoSpellIsCollected(t *testing.T) {
	t.Parallel()
	// Arrange
	service := bonuses.NewBonusEntryService()

	// Act
	summary := service.DescribeExistingBonuses([]config.BonusEntry{{PresetType: config.BonusSpell}})

	// Assert
	assert.Empty(t, summary.SpellIDs)
}

func TestWhenANonSpellBonusExists_NoSpellIsCollected(t *testing.T) {
	t.Parallel()
	// Arrange
	service := bonuses.NewBonusEntryService()

	// Act
	summary := service.DescribeExistingBonuses([]config.BonusEntry{{
		PresetType: config.BonusStartingItem,
		Param:      gofakeit.UUID(),
	}})

	// Assert
	assert.Empty(t, summary.SpellIDs)
}
