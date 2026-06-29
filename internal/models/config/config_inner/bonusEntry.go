package config_inner

import "strings"

// BonusEntry is the editor-side view-model for a single configurable game-start bonus
type BonusEntry struct {
	PresetType     BonusPresetType
	ReceiverFilter string // "start_hero" or "all_heroes".
	Param          string // Spell sid / item sid / numeric value depending on type.
	Param2         string // For Spell: "1" = free, "0" = normal. Unused for other types.
}

func (b BonusEntry) String() string {
	return b.PresetType.String() + "|" + b.ReceiverFilter + "|" + b.Param + "|" + b.Param2
}

// DeserializeBonusEntry deserializes a single line produced by BonusEntry.String.
// Returns ok=false for empty or malformed input.
func DeserializeBonusEntry(s string) (BonusEntry, bool) {
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

// DeserializeBonusEntries splits the persisted BonusesJSON string (newline-separated
// pipe-encoded entries) into individual BonusEntry values.
func DeserializeBonusEntries(s string) []BonusEntry {
	if s == "" {
		return nil
	}
	var out []BonusEntry
	for line := range strings.SplitSeq(s, "\n") {
		entry, ok := DeserializeBonusEntry(strings.TrimRight(line, "\r"))
		if ok {
			out = append(out, entry)
		}
	}
	return out
}

// SerializeBonusEntries joins a slice of bonus entries back to the canonical
// newline-separated string used by BonusesJson.
func SerializeBonusEntries(entries []BonusEntry) string {
	if len(entries) == 0 {
		return ""
	}
	lines := make([]string, 0, len(entries))
	for _, e := range entries {
		lines = append(lines, e.String())
	}
	return strings.Join(lines, "\n")
}
