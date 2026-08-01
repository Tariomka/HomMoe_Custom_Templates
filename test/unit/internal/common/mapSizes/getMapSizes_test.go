package mapSizes_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common"
	"github.com/stretchr/testify/assert"
)

func TestWhenExperimentalIsFalse_ReturnsBaseSizesOnly(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedCount := 11

	// Act
	result := common.GetMapSizes(false)

	// Assert
	assert.Len(t, result, expectedCount)
}

func TestWhenExperimentalIsTrue_ReturnsAllSizes(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedCount := 28

	// Act
	result := common.GetMapSizes(true)

	// Assert
	assert.Len(t, result, expectedCount)
}

func TestWhenBaseResultIsMutated_KeepsCatalogUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := common.GetMapSizes(false)
	result := common.GetMapSizes(false)

	// Act
	result[0].Size = 0

	// Assert
	assert.Equal(t, expected, common.GetMapSizes(false))
}

func TestWhenAllResultIsMutated_KeepsCatalogUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := common.GetMapSizes(true)
	result := common.GetMapSizes(true)

	// Act
	result[len(result)-1].Size = 0

	// Assert
	assert.Equal(t, expected, common.GetMapSizes(true))
}
