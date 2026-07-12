package mapSizes_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/constants"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenSizeIsAlreadyValid_ReturnsSameSize(t *testing.T) {
	t.Parallel()
	// Arrange
	valid := constants.AllMapSizes[gofakeit.Number(0, len(constants.AllMapSizes)-1)]

	// Act
	result := constants.GetNearestMapSize(valid.Size)

	// Assert
	assert.Equal(t, valid, result)
}

func TestWhenSizeIsBelowSmallest_ReturnsSmallestSize(t *testing.T) {
	t.Parallel()
	// Arrange
	tooSmall := gofakeit.Number(-1000, 63)

	// Act
	result := constants.GetNearestMapSize(tooSmall)

	// Assert
	assert.Equal(t, constants.BaseMapSizes[0], result)
}

func TestWhenSizeIsAboveLargest_ReturnsLargestSize(t *testing.T) {
	t.Parallel()
	// Arrange
	tooLarge := gofakeit.Number(513, 100000)
	largest := constants.ExpandedMapSizes[len(constants.ExpandedMapSizes)-1]

	// Act
	result := constants.GetNearestMapSize(tooLarge)

	// Assert
	assert.Equal(t, largest, result)
}

func TestWhenSizeIsBetweenTwoSizes_ReturnsClosestSize(t *testing.T) {
	t.Parallel()
	// Arrange
	betweenSize := 230 // between 208 and 240, closer to 240

	// Act
	result := constants.GetNearestMapSize(betweenSize)

	// Assert
	assert.Equal(t, 240, result.Size)
}

func TestWhenSizeIsEquidistantBetweenTwoSizes_ReturnsSmallerSize(t *testing.T) {
	t.Parallel()
	// Arrange
	equidistantSize := 72 // exactly between 64 and 80

	// Act
	result := constants.GetNearestMapSize(equidistantSize)

	// Assert
	assert.Equal(t, 64, result.Size)
}
