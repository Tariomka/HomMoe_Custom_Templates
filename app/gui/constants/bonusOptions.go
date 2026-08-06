package constants

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
)

// BonusTypeOption pairs a friendly dropdown label with its preset type.
type BonusTypeOption struct {
	Label      string
	PresetType config.BonusPresetType
}

func GetBonusTypeOptions() []BonusTypeOption {
	return []BonusTypeOption{
		{"Free Town Portal", config.BonusTownPortalFree},
		{"Spell", config.BonusSpell},
		{"Unit Multiplier", config.BonusUnitMultiplier},
		{"Movement Bonus", config.BonusMovementBonus},
		{"Starting Item", config.BonusStartingItem},
		{"Starting Gold", config.BonusStartingGold},
		{"Starting Gems", config.BonusStartingGems},
		{"Starting Crystals", config.BonusStartingCrystals},
		{"Starting Mercury", config.BonusStartingMercury},
		{"Starting Wood", config.BonusStartingWood},
		{"Starting Ore", config.BonusStartingOre},
	}
}

func GetBonusReceiverOptions() []string {
	receiversFilters := registry.GetReceiversFilterValues()
	return []string{receiversFilters.StartingHero, receiversFilters.AllHeroes}
}

func GetBonusResourceDefaults() map[config.BonusPresetType]string {
	return map[config.BonusPresetType]string{
		config.BonusStartingGold:     "10000",
		config.BonusStartingGems:     "15",
		config.BonusStartingCrystals: "15",
		config.BonusStartingMercury:  "15",
		config.BonusStartingWood:     "20",
		config.BonusStartingOre:      "20",
	}
}
