package helpers

import (
	"runtime"
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

	sanitized := string(out)
	if runtime.GOOS == windowsOS {
		reserved := []string{
			"CON", "PRN", "AUX", "NUL",
			"COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9",
			"LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9",
		}
		parts := strings.Split(sanitized, ".")
		if slices.Contains(reserved, strings.ToUpper(strings.TrimSpace(parts[0]))) {
			parts[0] = "_"
		}
		sanitized = strings.Join(parts, ".")
	}

	return sanitized
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
