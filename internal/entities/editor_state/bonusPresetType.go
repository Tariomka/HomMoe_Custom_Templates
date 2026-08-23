package editor_state

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
