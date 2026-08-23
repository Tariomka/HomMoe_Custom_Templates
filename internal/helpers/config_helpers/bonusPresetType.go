package config_helpers

import (
	"strconv"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
)

func GetString(presetType editor_state.BonusPresetType) string {
	switch presetType {
	case editor_state.BonusTownPortalFree:
		return "TownPortalFree"
	case editor_state.BonusSpell:
		return "Spell"
	case editor_state.BonusUnitMultiplier:
		return "UnitMultiplier"
	case editor_state.BonusMovementBonus:
		return "MovementBonus"
	case editor_state.BonusStartingItem:
		return "StartingItem"
	case editor_state.BonusStartingGold:
		return "StartingGold"
	case editor_state.BonusStartingGems:
		return "StartingGems"
	case editor_state.BonusStartingCrystals:
		return "StartingCrystals"
	case editor_state.BonusStartingMercury:
		return "StartingMercury"
	case editor_state.BonusStartingWood:
		return "StartingWood"
	case editor_state.BonusStartingOre:
		return "StartingOre"
	}
	return strconv.Itoa(int(presetType))
}

func IsResource(presetType editor_state.BonusPresetType) bool {
	return presetType >= editor_state.BonusStartingGold
}
