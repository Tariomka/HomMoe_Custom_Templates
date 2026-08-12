package mapSizes_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenSizeIsAlreadyValid_ReturnsSameSize(t *testing.T) {
	t.Parallel()
	// Arrange
	mapSizes := common.GetMapSizes(true)
	valid := mapSizes[gofakeit.Number(0, len(mapSizes)-1)]

	// Act
	result := common.GetNearestMapSize(valid.Size)

	// Assert
	assert.Equal(t, valid, result)
}

func TestWhenSizeIsBelowSmallest_ReturnsSmallestSize(t *testing.T) {
	t.Parallel()
	// Arrange
	tooSmall := gofakeit.Number(-1000, 63)
	expected := common.GetMapSizes(false)[0]

	// Act
	result := common.GetNearestMapSize(tooSmall)

	// Assert
	assert.Equal(t, expected, result)
}

func TestWhenSizeIsAboveLargest_ReturnsLargestSize(t *testing.T) {
	t.Parallel()
	// Arrange
	tooLarge := gofakeit.Number(513, 100000)
	mapSizes := common.GetMapSizes(true)
	largest := mapSizes[len(mapSizes)-1]

	// Act
	result := common.GetNearestMapSize(tooLarge)

	// Assert
	assert.Equal(t, largest, result)
}

func TestWhenSizeIsBetweenTwoSizes_ReturnsClosestSize(t *testing.T) {
	t.Parallel()
	// Arrange
	betweenSize := 230

	// Act
	result := common.GetNearestMapSize(betweenSize)

	// Assert
	assert.Equal(t, 240, result.Size)
}

func TestWhenSizeIsEquidistantBetweenTwoSizes_ReturnsSmallerSize(t *testing.T) {
	t.Parallel()
	// Arrange
	equidistantSize := 72

	// Act
	result := common.GetNearestMapSize(equidistantSize)

	// Assert
	assert.Equal(t, 64, result.Size)
}
