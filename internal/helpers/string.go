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
