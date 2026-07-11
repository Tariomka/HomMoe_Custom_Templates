package slice_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenQueryContainsDuplicates_KeepsFirstOccurrencesInOrder(t *testing.T) {
	t.Parallel()
	// Arrange
	source := []int{5, 1, 5, 2, 1}

	// Act
	distinct := linq.FromSlice(source).Distinct().ToSlice()

	// Assert
	assert.Equal(t, []int{5, 1, 2}, distinct)
}

func TestWhenQueryHasNoDuplicates_ReturnsAllElements(t *testing.T) {
	t.Parallel()
	// Arrange
	firstWord := gofakeit.Word() + "-first"
	secondWord := gofakeit.Word() + "-second"
	source := []string{firstWord, secondWord}

	// Act
	distinct := linq.FromSlice(source).Distinct().ToSlice()

	// Assert
	assert.Equal(t, []string{firstWord, secondWord}, distinct)
}

func TestWhenQueryIsEmpty_ReturnsEmptyResult(t *testing.T) {
	t.Parallel()
	// Arrange
	source := []int{}

	// Act
	distinct := linq.FromSlice(source).Distinct().ToSlice()

	// Assert
	assert.Empty(t, distinct)
}
