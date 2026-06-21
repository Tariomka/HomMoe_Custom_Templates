package utils

import "fmt"

// DenormalizeString returns [Denormalize] result as a string with a "±" prefix
func DenormalizeString(value, low, high float32) string {
	return fmt.Sprintf("± %.2f", Denormalize(value, low, high))
}

// RoundedRangeString returns [RoundedRange] result as a string
func RoundedRangeString(value float32, min, max int) string {
	return fmt.Sprintf("%d", RoundedRange(value, min, max))
}

// RoundedRangePercentString returns [RoundedRange] result as a percentage string, e.g. "50%"
func RoundedRangePercentString(value float32, min, max int) string {
	return fmt.Sprintf("%d%%", RoundedRange(value, min, max))
}

// MultiplierString returns [Multiplier] result as a string with a "x" prefix, e.g. "x 1.50"
func MultiplierString(value, base, factor float32) string {
	return fmt.Sprintf("x %.2f", Multiplier(value, base, factor))
}
