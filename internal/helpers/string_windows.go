//go:build windows

package helpers

import (
	"strings"
)

func removeReservedNames(input []rune) string {
	sanitized := string(input)

	parts := strings.Split(sanitized, ".")
	if isReservedFilename(parts[0]) {
		parts[0] = "_"
	}

	return strings.Join(parts, ".")
}

func isReservedFilename(name string) bool {
	stem, _, _ := strings.Cut(name, ".")
	switch strings.ToUpper(strings.TrimSpace(stem)) {
	case "CON", "PRN", "AUX", "NUL",
		"COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9",
		"LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9":
		return true
	default:
		return false
	}
}
