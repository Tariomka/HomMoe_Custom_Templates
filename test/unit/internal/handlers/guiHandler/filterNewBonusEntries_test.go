package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/config_helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenBonusEntriesAreFiltered_DropsTheOnesThatAlreadyExist(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	known := config.BonusEntry{PresetType: config.BonusStartingGold, Param: gofakeit.DigitN(4)}

	// Act
	fresh := handler.FilterNewBonusEntries(
		[]config.BonusEntry{known},
		map[string]bool{config_helpers.GetHash(known): true})

	// Assert
	assert.Empty(t, fresh)
}
