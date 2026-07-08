package slice_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/stretchr/testify/assert"
)

func TestWhenSelectorKeysCollide_KeepsFirstElementPerKeyInOrder(t *testing.T) {
	// Arrange
	source := []string{"apple", "avocado", "banana", "cherry", "citrus"}
	firstLetter := func(word string) any { return word[0] }

	// Act
	distinct := linq.FromSlice(source).DistinctBy(firstLetter).ToSlice()

	// Assert
	assert.Equal(t, []string{"apple", "banana", "cherry"}, distinct)
}

func TestWhenAllSelectorKeysAreUnique_ReturnsAllElements(t *testing.T) {
	// Arrange
	source := []int{10, 21, 32}
	moduloTen := func(number int) any { return number % 10 }

	// Act
	distinct := linq.FromSlice(source).DistinctBy(moduloTen).ToSlice()

	// Assert
	assert.Equal(t, []int{10, 21, 32}, distinct)
}

func TestWhenQueryIsEmptyAndSelectorProvided_ReturnsEmptyResult(t *testing.T) {
	// Arrange
	source := []int{}
	identity := func(number int) any { return number }

	// Act
	distinct := linq.FromSlice(source).DistinctBy(identity).ToSlice()

	// Assert
	assert.Empty(t, distinct)
}
