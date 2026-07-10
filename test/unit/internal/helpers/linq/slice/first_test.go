package slice_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenMultipleElementsMatch_ReturnsFirstMatchWithOkTrue(t *testing.T) {
	t.Parallel()
	// Arrange
	source := []int{1, 4, 6, 8}
	isEven := func(number int) bool { return number%2 == 0 }

	// Act
	found, ok := linq.FromSlice(source).First(isEven)

	// Assert
	assert.Equal(t, []any{4, true}, []any{found, ok})
}

func TestWhenNoElementMatches_ReturnsZeroValueWithOkFalse(t *testing.T) {
	t.Parallel()
	// Arrange
	source := []string{gofakeit.Word(), gofakeit.Word()}
	matchesNothing := func(string) bool { return false }

	// Act
	found, ok := linq.FromSlice(source).First(matchesNothing)

	// Assert
	assert.Equal(t, []any{"", false}, []any{found, ok})
}

func TestWhenSourceIsEmpty_ReturnsZeroValueWithOkFalse(t *testing.T) {
	t.Parallel()
	// Arrange
	source := []int{}
	matchesEverything := func(int) bool { return true }

	// Act
	found, ok := linq.FromSlice(source).First(matchesEverything)

	// Assert
	assert.Equal(t, []any{0, false}, []any{found, ok})
}
