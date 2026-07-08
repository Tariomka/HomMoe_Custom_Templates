package slice_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenSliceContainsDuplicates_ReturnsFirstOccurrencesInOrder(t *testing.T) {
	// Arrange
	inputSlice := []int{3, 1, 3, 2, 1, 3}

	// Act
	uniqueSlice := helpers.GetUniqueElements(inputSlice)

	// Assert
	assert.Equal(t, []int{3, 1, 2}, uniqueSlice)
}

func TestWhenSliceHasNoDuplicates_ReturnsItUnchanged(t *testing.T) {
	// Arrange
	firstWord := gofakeit.Word() + "-first"
	secondWord := gofakeit.Word() + "-second"
	thirdWord := gofakeit.Word() + "-third"
	inputSlice := []string{firstWord, secondWord, thirdWord}

	// Act
	uniqueSlice := helpers.GetUniqueElements(inputSlice)

	// Assert
	assert.Equal(t, []string{firstWord, secondWord, thirdWord}, uniqueSlice)
}

func TestWhenAllElementsAreEqual_ReturnsSingleElement(t *testing.T) {
	// Arrange
	repeatedNumber := gofakeit.Number(1, 1000)
	inputSlice := []int{repeatedNumber, repeatedNumber, repeatedNumber}

	// Act
	uniqueSlice := helpers.GetUniqueElements(inputSlice)

	// Assert
	assert.Equal(t, []int{repeatedNumber}, uniqueSlice)
}

func TestWhenSliceIsEmpty_ReturnsEmptySlice(t *testing.T) {
	// Arrange
	inputSlice := []int{}

	// Act
	uniqueSlice := helpers.GetUniqueElements(inputSlice)

	// Assert
	assert.Empty(t, uniqueSlice)
}

func TestWhenSliceIsNil_ReturnsEmptySlice(t *testing.T) {
	// Arrange
	var inputSlice []string

	// Act
	uniqueSlice := helpers.GetUniqueElements(inputSlice)

	// Assert
	assert.Empty(t, uniqueSlice)
}
