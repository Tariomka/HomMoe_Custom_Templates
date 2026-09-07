package bonusEntry_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenThereAreNoBonuses_ReturnsNil(t *testing.T) {
	t.Parallel()
	// Arrange
	var bonuses []editor_state.BonusEntry

	// Act
	models := editor_state_model.ToBonusEntryModels(bonuses)

	// Assert
	assert.Nil(t, models)
}

func TestWhenBonusesArePersisted_EachOneIsWrapped(t *testing.T) {
	t.Parallel()
	// Arrange
	bonuses := []editor_state.BonusEntry{
		{PresetType: editor_state.BonusStartingGold, Param: "9500"},
		{PresetType: editor_state.BonusSpell, Param: "magic_arrow", Param2: "1"},
	}

	// Act
	models := editor_state_model.ToBonusEntryModels(bonuses)

	// Assert
	assert.Equal(
		t,
		[]editor_state_model.BonusEntry{{BonusEntry: bonuses[0]}, {BonusEntry: bonuses[1]}},
		models)
}
