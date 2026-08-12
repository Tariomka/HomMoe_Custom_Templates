package mapSizes_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenSizeMatchesKnownSize_ReturnsMatchingEntry(t *testing.T) {
	t.Parallel()
	// Arrange
	mapSizes := common.GetMapSizes(true)
	expected := mapSizes[gofakeit.Number(0, len(mapSizes)-1)]

	// Act
	result := common.GetMapSize(expected.Size)

	// Assert
	assert.Equal(t, expected, result)
}

func TestWhenSizeIsUnknown_ReturnsSmallestBaseSize(t *testing.T) {
	t.Parallel()
	// Arrange
	unknownSize := 999
	expected := common.GetMapSizes(false)[0]

	// Act
	result := common.GetMapSize(unknownSize)

	// Assert
	assert.Equal(t, expected, result)
}
