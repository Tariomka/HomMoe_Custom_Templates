package slice_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenSourceSliceGiven_IterationYieldsAllElementsInOrder(t *testing.T) {
	// Arrange
	firstWord := gofakeit.Word() + "-first"
	secondWord := gofakeit.Word() + "-second"
	thirdWord := gofakeit.Word() + "-third"
	source := []string{firstWord, secondWord, thirdWord}

	// Act
	var collected []string
	for item := range linq.FromSlice(source).Iterate {
		collected = append(collected, item)
	}

	// Assert
	assert.Equal(t, []string{firstWord, secondWord, thirdWord}, collected)
}

func TestWhenYieldReturnsFalse_IterationStopsEarly(t *testing.T) {
	// Arrange
	source := []int{1, 2, 3, 4}

	// Act
	var collected []int
	linq.FromSlice(source).Iterate(func(item int) bool {
		collected = append(collected, item)
		return len(collected) < 2
	})

	// Assert
	assert.Equal(t, []int{1, 2}, collected)
}

func TestWhenSourceSliceIsEmpty_IterationYieldsNothing(t *testing.T) {
	// Arrange
	source := []int{}

	// Act
	yieldCount := 0
	linq.FromSlice(source).Iterate(func(item int) bool {
		yieldCount++
		return true
	})

	// Assert
	assert.Equal(t, 0, yieldCount)
}
