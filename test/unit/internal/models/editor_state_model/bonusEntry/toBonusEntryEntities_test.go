package bonusEntry_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenThereAreNoBonusModels_ReturnsNil(t *testing.T) {
	t.Parallel()
	// Arrange
	var bonuses []editor_state_model.BonusEntry

	// Act
	entities := editor_state_model.ToBonusEntryEntities(bonuses)

	// Assert
	assert.Nil(t, entities)
}

func TestWhenBonusModelsAreUnwrapped_TheEntitiesAreCarried(t *testing.T) {
	t.Parallel()
	// Arrange
	bonuses := []editor_state_model.BonusEntry{
		{PresetType: editor_state.BonusStartingGold, Param: "9500", ReceiverFilter: "all_heroes"},
	}

	// Act
	entities := editor_state_model.ToBonusEntryEntities(bonuses)

	// Assert
	assert.Equal(t, []editor_state.BonusEntry{bonuses[0].BonusEntry}, entities)
}
