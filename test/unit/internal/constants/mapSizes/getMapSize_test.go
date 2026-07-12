package mapSizes_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/constants"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenSizeMatchesKnownSize_ReturnsMatchingEntry(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := constants.AllMapSizes[gofakeit.Number(0, len(constants.AllMapSizes)-1)]

	// Act
	result := constants.GetMapSize(expected.Size)

	// Assert
	assert.Equal(t, expected, result)
}

func TestWhenSizeIsUnknown_ReturnsSmallestBaseSize(t *testing.T) {
	t.Parallel()
	// Arrange
	unknownSize := 999

	// Act
	result := constants.GetMapSize(unknownSize)

	// Assert
	assert.Equal(t, constants.BaseMapSizes[0], result)
}
