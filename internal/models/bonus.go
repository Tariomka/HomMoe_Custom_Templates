package models

import "github.com/Tariomka/hommoe_custom_templates/internal/models/generator"

// ParseBonusesJSON is a thin re-export to keep call sites short.
func ParseBonusesJSON(s string) []generator.BonusEntry { return generator.ParseBonusesJSON(s) }

// SerializeBonuses is a thin re-export to keep call sites short.
func SerializeBonuses(entries []generator.BonusEntry) string {
	return generator.SerializeBonuses(entries)
}

// ParseBonusEntry is a thin re-export to keep call sites short.
func ParseBonusEntry(s string) (generator.BonusEntry, bool) { return generator.ParseBonusEntry(s) }
