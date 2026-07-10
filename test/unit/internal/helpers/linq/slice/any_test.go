package slice_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenQueryHasElements_ReturnsTrue(t *testing.T) {
	t.Parallel()
	// Arrange
	source := []string{gofakeit.Word()}

	// Act
	hasAny := linq.FromSlice(source).Any()

	// Assert
	assert.True(t, hasAny)
}

func TestWhenQueryIsEmpty_ReturnsFalse(t *testing.T) {
	t.Parallel()
	// Arrange
	source := []string{}

	// Act
	hasAny := linq.FromSlice(source).Any()

	// Assert
	assert.False(t, hasAny)
}

func TestWhenFilteredQueryHasNoMatches_ReturnsFalse(t *testing.T) {
	t.Parallel()
	// Arrange
	source := []int{1, 3, 5}
	isEven := func(number int) bool { return number%2 == 0 }

	// Act
	hasAny := linq.FromSlice(source).Where(isEven).Any()

	// Assert
	assert.False(t, hasAny)
}
