package slice_test

import (
	"strconv"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenSelectorGiven_ProjectsEveryElement(t *testing.T) {
	t.Parallel()
	// Arrange
	firstNumber := gofakeit.Number(1, 1000)
	secondNumber := gofakeit.Number(1, 1000)
	source := []int{firstNumber, secondNumber}

	// Act
	projected := linq.FromSlice(source).Select(strconv.Itoa).ToSlice()

	// Assert
	assert.Equal(t, []string{strconv.Itoa(firstNumber), strconv.Itoa(secondNumber)}, projected)
}

func TestWhenSelectorProjectsToStruct_ProjectsEveryElement(t *testing.T) {
	t.Parallel()
	// Arrange
	type wrapper struct {
		Value string
	}

	firstWord := gofakeit.Word()
	secondWord := gofakeit.Word()
	source := []string{firstWord, secondWord}

	// Act
	projected := linq.FromSlice(source).
		Select(func(word string) wrapper { return wrapper{Value: word} }).
		ToSlice()

	// Assert
	assert.Equal(t, []wrapper{{Value: firstWord}, {Value: secondWord}}, projected)
}

func TestWhenSelectIsChained_ProjectsThroughBothSelectors(t *testing.T) {
	t.Parallel()
	// Arrange
	number := gofakeit.Number(1, 1000)
	source := []int{number}

	// Act
	projected := linq.FromSlice(source).
		Select(strconv.Itoa).
		Select(func(text string) int { return len(text) }).
		ToSlice()

	// Assert
	assert.Equal(t, []int{len(strconv.Itoa(number))}, projected)
}

func TestWhenSourceIsEmpty_ProjectionYieldsNothing(t *testing.T) {
	t.Parallel()
	// Arrange
	source := []int{}

	// Act
	projected := linq.FromSlice(source).Select(strconv.Itoa).ToSlice()

	// Assert
	assert.Empty(t, projected)
}
