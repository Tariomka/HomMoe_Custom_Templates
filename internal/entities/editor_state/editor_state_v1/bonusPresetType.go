package editor_state_v1

// BonusPresetType is the frozen v1 snapshot - see editorState.go. Never edit.
// The ordinals are the wire values a v1 file carries, so they stay pinned here
// even if the current enum is reordered.
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
