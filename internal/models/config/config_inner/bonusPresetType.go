package config_inner

import (
	"strconv"
	"strings"
)

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

func (this BonusPresetType) IsResource() bool {
	return this >= BonusStartingGold
}

func parseBonusPresetType(s string) (BonusPresetType, bool) {
	if n, err := strconv.Atoi(s); err == nil {
		return BonusPresetType(n), true
	}
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "townportalfree":
		return BonusTownPortalFree, true
	case "spell":
		return BonusSpell, true
	case "unitmultiplier":
		return BonusUnitMultiplier, true
	case "movementbonus":
		return BonusMovementBonus, true
	case "startingitem":
		return BonusStartingItem, true
	case "startinggold":
		return BonusStartingGold, true
	case "startinggems":
		return BonusStartingGems, true
	case "startingcrystals":
		return BonusStartingCrystals, true
	case "startingmercury":
		return BonusStartingMercury, true
	case "startingwood":
		return BonusStartingWood, true
	case "startingore":
		return BonusStartingOre, true
	}
	return 0, false
}
