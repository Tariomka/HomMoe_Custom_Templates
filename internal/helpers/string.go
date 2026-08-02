package helpers

import (
	"slices"
	"strings"
	"unicode"
)

func SanitizeFilename(name string) string {
	bad := []rune{'/', '\\', ':', '*', '?', '"', '<', '>', '|'}
	out := []rune(strings.TrimSpace(name))
	for index, runeValue := range out {
		if slices.Contains(bad, runeValue) || runeValue < 32 || runeValue == 127 || !unicode.IsPrint(runeValue) {
			out[index] = '_'
		}
	}

	return removeReservedNames(out)
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
