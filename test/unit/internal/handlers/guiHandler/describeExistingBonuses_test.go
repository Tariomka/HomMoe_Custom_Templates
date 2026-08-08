package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenExistingBonusesDescribed_CollectsTheirSpellIds(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	spellID := gofakeit.UUID()

	// Act
	summary := handler.DescribeExistingBonuses([]config.BonusEntry{{
		PresetType: config.BonusSpell,
		Param:      spellID,
	}})

	// Assert
	assert.Equal(t, []string{spellID}, summary.SpellIDs)
}
