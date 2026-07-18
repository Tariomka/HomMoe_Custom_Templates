package helpers

import "strings"

// SanitizeFilename replaces filesystem-unsafe runes in name with underscores
// and trims surrounding whitespace.
func SanitizeFilename(name string) string {
	bad := []rune{'/', '\\', ':', '*', '?', '"', '<', '>', '|'}
	out := []rune(strings.TrimSpace(name))
	for index, runeValue := range out {
		for _, badRune := range bad {
			if runeValue == badRune {
				out[index] = '_'
			}
		}
	}
	return string(out)
}

// GetZoneLabel returns the trailing label portion of a zone name like
// "Spawn-A" → "A" or "Neutral-C" → "C". Plain names (e.g. "Hub") pass through.
func GetZoneLabel(zoneName string) string {
	if after, ok := strings.CutPrefix(zoneName, "Spawn-"); ok {
		return after
	}

	if after, ok := strings.CutPrefix(zoneName, "Neutral-"); ok {
		return after
	}

	if _, after, ok := strings.Cut(zoneName, "-"); ok {
		return after
	}

	return zoneName
}
