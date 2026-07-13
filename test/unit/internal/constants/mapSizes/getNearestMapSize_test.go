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
	valid := common.AllMapSizes[gofakeit.Number(0, len(common.AllMapSizes)-1)]

	// Act
	result := common.GetNearestMapSize(valid.Size)

	// Assert
	assert.Equal(t, valid, result)
}

func TestWhenSizeIsBelowSmallest_ReturnsSmallestSize(t *testing.T) {
	t.Parallel()
	// Arrange
	tooSmall := gofakeit.Number(-1000, 63)

	// Act
	result := common.GetNearestMapSize(tooSmall)

	// Assert
	assert.Equal(t, common.BaseMapSizes[0], result)
}

func TestWhenSizeIsAboveLargest_ReturnsLargestSize(t *testing.T) {
	t.Parallel()
	// Arrange
	tooLarge := gofakeit.Number(513, 100000)
	largest := common.ExpandedMapSizes[len(common.ExpandedMapSizes)-1]

	// Act
	result := common.GetNearestMapSize(tooLarge)

	// Assert
	assert.Equal(t, largest, result)
}

func TestWhenSizeIsBetweenTwoSizes_ReturnsClosestSize(t *testing.T) {
	t.Parallel()
	// Arrange
	betweenSize := 230 // between 208 and 240, closer to 240

	// Act
	result := common.GetNearestMapSize(betweenSize)

	// Assert
	assert.Equal(t, 240, result.Size)
}

func TestWhenSizeIsEquidistantBetweenTwoSizes_ReturnsSmallerSize(t *testing.T) {
	t.Parallel()
	// Arrange
	equidistantSize := 72 // exactly between 64 and 80

	// Act
	result := common.GetNearestMapSize(equidistantSize)

	// Assert
	assert.Equal(t, 64, result.Size)
}
