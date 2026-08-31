package slice_test

import (
	"strconv"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenSelectSliceGetsElements_ProjectsEveryElementInOrder(t *testing.T) {
	t.Parallel()
	// Arrange
	firstNumber := gofakeit.Number(1, 1000)
	secondNumber := gofakeit.Number(1, 1000)
	source := []int{firstNumber, secondNumber}

	// Act
	projected := linq.SelectSlice(source, strconv.Itoa)

	// Assert
	assert.Equal(t, []string{strconv.Itoa(firstNumber), strconv.Itoa(secondNumber)}, projected)
}

func TestWhenSelectSliceGetsEmptySlice_ReturnsNilSlice(t *testing.T) {
	t.Parallel()
	// Arrange
	source := []int{}

	// Act
	projected := linq.SelectSlice(source, strconv.Itoa)

	// Assert
	assert.Nil(t, projected)
}

func TestWhenSelectSliceGetsNilSlice_ReturnsNilSlice(t *testing.T) {
	t.Parallel()
	// Arrange
	var source []int

	// Act
	projected := linq.SelectSlice(source, strconv.Itoa)

	// Assert
	assert.Nil(t, projected)
}

func TestWhenSelectSliceResultIsMutated_SourceSliceStaysUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	source := []int{10, 20, 30}

	// Act
	projected := linq.SelectSlice(source, func(number int) int { return number })
	projected[0] = 999

	// Assert
	assert.Equal(t, []int{10, 20, 30}, source)
}

func TestWhenSelectSliceGetsNamedSliceType_ProjectsEveryElement(t *testing.T) {
	t.Parallel()
	// Arrange
	type numbers []int

	number := gofakeit.Number(1, 1000)
	source := numbers{number}

	// Act
	projected := linq.SelectSlice(source, strconv.Itoa)

	// Assert
	assert.Equal(t, []string{strconv.Itoa(number)}, projected)
}
