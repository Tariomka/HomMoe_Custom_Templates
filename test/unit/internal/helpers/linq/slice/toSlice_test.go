package slice_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenQueryHasElements_ReturnsThemAsSliceInOrder(t *testing.T) {
	// Arrange
	firstWord := gofakeit.Word() + "-first"
	secondWord := gofakeit.Word() + "-second"
	source := []string{firstWord, secondWord}

	// Act
	result := linq.FromSlice(source).ToSlice()

	// Assert
	assert.Equal(t, []string{firstWord, secondWord}, result)
}

func TestWhenResultSliceIsMutated_SourceSliceStaysUnchanged(t *testing.T) {
	// Arrange
	source := []int{10, 20, 30}

	// Act
	result := linq.FromSlice(source).ToSlice()
	result[0] = 999

	// Assert
	assert.Equal(t, []int{10, 20, 30}, source)
}

func TestWhenQueryIsEmpty_ReturnsNilSlice(t *testing.T) {
	// Arrange
	source := []int{}

	// Act
	result := linq.FromSlice(source).ToSlice()

	// Assert
	assert.Nil(t, result)
}
