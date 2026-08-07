//go:build !windows

package helpers

func removeReservedNames(input []rune) string { return string(input) }

func isReservedFilename(_ string) bool { return false }
