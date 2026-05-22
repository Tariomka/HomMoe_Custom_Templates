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

// ParseBonusesJSON is a thin re-export to keep call sites short.
func ParseBonusesJSON(s string) []generator.BonusEntry { return generator.ParseBonusesJSON(s) }

// SerialiseBonuses is a thin re-export to keep call sites short.
func SerialiseBonuses(entries []generator.BonusEntry) string {
	return generator.SerialiseBonuses(entries)
}

// ParseBonusEntry is a thin re-export to keep call sites short.
func ParseBonusEntry(s string) (generator.BonusEntry, bool) { return generator.ParseBonusEntry(s) }
