package slice_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenMultipleElementsMatch_ReturnsFirstMatchingElement(t *testing.T) {
	t.Parallel()
	// Arrange
	source := []int{1, 4, 6, 8}
	isEven := func(number int) bool { return number%2 == 0 }

	// Act
	found := linq.FromSlice(source).FirstOrDefault(isEven)

	// Assert
	assert.Equal(t, 4, found)
}

func TestWhenNoElementMatches_ReturnsZeroValue(t *testing.T) {
	t.Parallel()
	// Arrange
	source := []string{gofakeit.Word(), gofakeit.Word()}
	matchesNothing := func(string) bool { return false }

	// Act
	found := linq.FromSlice(source).FirstOrDefault(matchesNothing)

	// Assert
	assert.Empty(t, found)
}

func TestWhenSourceIsEmpty_ReturnsZeroValue(t *testing.T) {
	t.Parallel()
	// Arrange
	source := []int{}
	matchesEverything := func(int) bool { return true }

	// Act
	found := linq.FromSlice(source).FirstOrDefault(matchesEverything)

	// Assert
	assert.Equal(t, 0, found)
}
