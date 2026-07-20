package slice_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenAnElementMatchesPredicate_ReturnsTrue(t *testing.T) {
	t.Parallel()
	// Arrange
	match := gofakeit.Word()
	source := []string{gofakeit.Word() + "_other", match}

	// Act
	hasMatch := linq.FromSlice(source).AnyFunc(func(item string) bool { return item == match })

	// Assert
	assert.True(t, hasMatch)
}

func TestWhenNoElementMatchesPredicate_ReturnsFalse(t *testing.T) {
	t.Parallel()
	// Arrange
	source := []int{1, 3, 5}
	isEven := func(number int) bool { return number%2 == 0 }

	// Act
	hasMatch := linq.FromSlice(source).AnyFunc(isEven)

	// Assert
	assert.False(t, hasMatch)
}

func TestWhenSourceIsEmpty_AnyFuncReturnsFalse(t *testing.T) {
	t.Parallel()
	// Arrange
	source := []string{}

	// Act
	hasMatch := linq.FromSlice(source).AnyFunc(func(string) bool { return true })

	// Assert
	assert.False(t, hasMatch)
}
