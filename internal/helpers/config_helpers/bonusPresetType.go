package config_helpers

import (
	"strconv"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
)

func GetString(presetType editor_state_model.BonusPresetType) string {
	switch presetType {
	case editor_state_model.BonusTownPortalFree:
		return "TownPortalFree"
	case editor_state_model.BonusSpell:
		return "Spell"
	case editor_state_model.BonusUnitMultiplier:
		return "UnitMultiplier"
	case editor_state_model.BonusMovementBonus:
		return "MovementBonus"
	case editor_state_model.BonusStartingItem:
		return "StartingItem"
	case editor_state_model.BonusStartingGold:
		return "StartingGold"
	case editor_state_model.BonusStartingGems:
		return "StartingGems"
	case editor_state_model.BonusStartingCrystals:
		return "StartingCrystals"
	case editor_state_model.BonusStartingMercury:
		return "StartingMercury"
	case editor_state_model.BonusStartingWood:
		return "StartingWood"
	case editor_state_model.BonusStartingOre:
		return "StartingOre"
	}
	return strconv.Itoa(int(presetType))
}

func IsResource(presetType editor_state_model.BonusPresetType) bool {
	return presetType >= editor_state_model.BonusStartingGold
}
