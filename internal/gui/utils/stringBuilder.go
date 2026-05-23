package utils

import "fmt"

// RoundedRangeString return formatted text of [RoundedRange]
func RoundedRangeString(value float32, min, max int) string {
	return fmt.Sprintf("%d", RoundedRange(value, min, max))
}
