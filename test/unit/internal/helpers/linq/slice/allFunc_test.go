package slice_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenEveryElementMatchesPredicate_ReturnsTrue(t *testing.T) {
	t.Parallel()
	// Arrange
	source := []int{2, 4, gofakeit.Number(1, 100) * 2}
	isEven := func(number int) bool { return number%2 == 0 }

	// Act
	allMatch := linq.FromSlice(source).AllFunc(isEven)

	// Assert
	assert.True(t, allMatch)
}

func TestWhenOneElementFailsPredicate_ReturnsFalse(t *testing.T) {
	t.Parallel()
	// Arrange
	source := []int{2, 4, 5}
	isEven := func(number int) bool { return number%2 == 0 }

	// Act
	allMatch := linq.FromSlice(source).AllFunc(isEven)

	// Assert
	assert.False(t, allMatch)
}

func TestWhenSourceIsEmpty_AllFuncReturnsFalse(t *testing.T) {
	t.Parallel()
	// Arrange
	source := []string{}

	// Act
	allMatch := linq.FromSlice(source).AllFunc(func(string) bool { return true })

	// Assert
	assert.False(t, allMatch)
}
