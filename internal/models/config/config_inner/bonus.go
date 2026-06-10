package config_inner

import (
	"strconv"
	"strings"
)

// BonusPresetType enumerates the configurable game-start bonus presets.
type BonusPresetType int

const (
	BonusTownPortalFree   BonusPresetType = 0
	BonusSpell            BonusPresetType = 1
	BonusUnitMultiplier   BonusPresetType = 2
	BonusMovementBonus    BonusPresetType = 3
	BonusStartingItem     BonusPresetType = 4
	BonusStartingGold     BonusPresetType = 5
	BonusStartingGems     BonusPresetType = 6
	BonusStartingCrystals BonusPresetType = 7
	BonusStartingMercury  BonusPresetType = 8
	BonusStartingWood     BonusPresetType = 9
	BonusStartingOre      BonusPresetType = 10
)

func (t BonusPresetType) String() string {
	switch t {
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
	return strconv.Itoa(int(t))
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

// BonusEntry is the editor-side view-model for a single configurable game-start bonus
type BonusEntry struct {
	PresetType BonusPresetType
	// "start_hero" or "all_heroes".
	ReceiverFilter string
	// Spell sid / item sid / numeric value depending on type.
	Param string
	// For Spell: "1" = free, "0" = normal. Unused for other types.
	Param2 string
}

func (b BonusEntry) String() string {
	return b.PresetType.String() + "|" + b.ReceiverFilter + "|" + b.Param + "|" + b.Param2
}

// ParseBonusEntry deserializes a single line produced by BonusEntry.String.
// Returns ok=false for empty or malformed input.
func ParseBonusEntry(s string) (BonusEntry, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return BonusEntry{}, false
	}
	parts := strings.Split(s, "|")
	if len(parts) < 4 {
		return BonusEntry{}, false
	}
	preset, ok := parseBonusPresetType(parts[0])
	if !ok {
		return BonusEntry{}, false
	}
	return BonusEntry{
		PresetType:     preset,
		ReceiverFilter: parts[1],
		Param:          parts[2],
		Param2:         parts[3],
	}, true
}

// ParseBonusesJSON splits the persisted BonusesJSON string (newline-separated
// pipe-encoded entries) into individual BonusEntry values.
func ParseBonusesJSON(s string) []BonusEntry {
	if s == "" {
		return nil
	}
	var out []BonusEntry
	for line := range strings.SplitSeq(s, "\n") {
		entry, ok := ParseBonusEntry(strings.TrimRight(line, "\r"))
		if ok {
			out = append(out, entry)
		}
	}
	return out
}

// SerializeBonuses joins a slice of bonus entries back to the canonical
// newline-separated string used by BonusesJson.
func SerializeBonuses(entries []BonusEntry) string {
	if len(entries) == 0 {
		return ""
	}
	lines := make([]string, 0, len(entries))
	for _, e := range entries {
		lines = append(lines, e.String())
	}
	return strings.Join(lines, "\n")
}
