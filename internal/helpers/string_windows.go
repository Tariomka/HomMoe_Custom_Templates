//go:build windows

package helpers

import (
	"slices"
	"strings"
)

func removeReservedNames(input []rune) string {
	sanitized := string(input)
	reserved := []string{
		"CON", "PRN", "AUX", "NUL",
		"COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9",
		"LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9",
	}
	parts := strings.Split(sanitized, ".")
	if slices.Contains(reserved, strings.ToUpper(strings.TrimSpace(parts[0]))) {
		parts[0] = "_"
	}

	return strings.Join(parts, ".")
}
