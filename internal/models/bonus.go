package models

import "github.com/Tariomka/hommoe_custom_templates/internal/models/generator"

// Bonus preset enum re-exports for ergonomic access from the UI / loader
// without pulling in the generator subpackage directly.
const (
	BonusTownPortalFree   = generator.BonusTownPortalFree
	BonusSpell            = generator.BonusSpell
	BonusUnitMultiplier   = generator.BonusUnitMultiplier
	BonusMovementBonus    = generator.BonusMovementBonus
	BonusStartingItem     = generator.BonusStartingItem
	BonusStartingGold     = generator.BonusStartingGold
	BonusStartingGems     = generator.BonusStartingGems
	BonusStartingCrystals = generator.BonusStartingCrystals
	BonusStartingMercury  = generator.BonusStartingMercury
	BonusStartingWood     = generator.BonusStartingWood
	BonusStartingOre      = generator.BonusStartingOre
)

// ParseBonusesJson is a thin re-export to keep call sites short.
func ParseBonusesJson(s string) []BonusEntry { return generator.ParseBonusesJson(s) }

// SerialiseBonuses is a thin re-export to keep call sites short.
func SerialiseBonuses(entries []BonusEntry) string { return generator.SerialiseBonuses(entries) }

// ParseBonusEntry is a thin re-export to keep call sites short.
func ParseBonusEntry(s string) (BonusEntry, bool) { return generator.ParseBonusEntry(s) }
