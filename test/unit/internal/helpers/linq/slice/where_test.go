package slice_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/stretchr/testify/assert"
)

func TestWhenPredicateMatchesSomeElements_KeepsOnlyMatchingOnesInOrder(t *testing.T) {
	t.Parallel()
	// Arrange
	source := []int{1, 2, 3, 4, 5, 6}
	isEven := func(number int) bool { return number%2 == 0 }

	// Act
	filtered := linq.FromSlice(source).Where(isEven).ToSlice()

	// Assert
	assert.Equal(t, []int{2, 4, 6}, filtered)
}

func TestWhenPredicateMatchesNoElements_ReturnsEmptyResult(t *testing.T) {
	t.Parallel()
	// Arrange
	source := []int{1, 3, 5}
	isEven := func(number int) bool { return number%2 == 0 }

	// Act
	filtered := linq.FromSlice(source).Where(isEven).ToSlice()

	// Assert
	assert.Empty(t, filtered)
}

func TestWhenPredicateMatchesAllElements_ReturnsAllOfThem(t *testing.T) {
	t.Parallel()
	// Arrange
	source := []int{2, 4, 6}
	isEven := func(number int) bool { return number%2 == 0 }

	// Act
	filtered := linq.FromSlice(source).Where(isEven).ToSlice()

	// Assert
	assert.Equal(t, []int{2, 4, 6}, filtered)
}

func TestWhenDownstreamStopsIteration_FilteredIterationStopsEarly(t *testing.T) {
	t.Parallel()
	// Arrange
	source := []int{1, 2, 3, 4, 5, 6}
	isEven := func(number int) bool { return number%2 == 0 }

	// Act
	var collected []int
	linq.FromSlice(source).Where(isEven).Iterate(func(item int) bool {
		collected = append(collected, item)
		return false
	})

	// Assert
	assert.Equal(t, []int{2}, collected)
}
