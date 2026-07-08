package config_inner

import "strconv"

// BonusPresetType enumerates the configurable game-start bonus presets.
type BonusPresetType int

const (
	BonusTownPortalFree BonusPresetType = iota
	BonusSpell
	BonusUnitMultiplier
	BonusMovementBonus
	BonusStartingItem
	BonusStartingGold
	BonusStartingGems
	BonusStartingCrystals
	BonusStartingMercury
	BonusStartingWood
	BonusStartingOre
)

func (this BonusPresetType) String() string {
	switch this {
	case BonusTownPortalFree:
		return "TownPortalFree"
	case BonusSpell:
		return "Spell"
	case BonusUnitMultiplier:
		return "UnitMultiplier"
	case BonusMovementBonus:
		return "MovementBonus"
	case BonusStartingItem:
		return "StartingItem"
	case BonusStartingGold:
		return "StartingGold"
	case BonusStartingGems:
		return "StartingGems"
	case BonusStartingCrystals:
		return "StartingCrystals"
	case BonusStartingMercury:
		return "StartingMercury"
	case BonusStartingWood:
		return "StartingWood"
	case BonusStartingOre:
		return "StartingOre"
	}
	return strconv.Itoa(int(this))
}

func (this BonusPresetType) IsResource() bool { return this >= BonusStartingGold }
