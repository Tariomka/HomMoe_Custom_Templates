package slice_test

import (
	"strconv"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenSelectorGiven_ProjectsEveryElementToString(t *testing.T) {
	// Arrange
	firstNumber := gofakeit.Number(1, 1000)
	secondNumber := gofakeit.Number(1, 1000)
	source := []int{firstNumber, secondNumber}

	// Act
	projected := linq.FromSlice(source).SelectString(strconv.Itoa).ToSlice()

	// Assert
	assert.Equal(t, []string{strconv.Itoa(firstNumber), strconv.Itoa(secondNumber)}, projected)
}

func TestWhenSourceIsEmpty_ProjectionYieldsNothing(t *testing.T) {
	// Arrange
	source := []int{}

	// Act
	projected := linq.FromSlice(source).SelectString(strconv.Itoa).ToSlice()

	// Assert
	assert.Empty(t, projected)
}
