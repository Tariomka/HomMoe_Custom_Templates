package mapSizes_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/constants"
	"github.com/stretchr/testify/assert"
)

func TestWhenExperimentalIsFalse_ReturnsBaseSizesOnly(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	result := constants.GetMapSizes(false)

	// Assert
	assert.Equal(t, constants.BaseMapSizes, result)
}

func TestWhenExperimentalIsTrue_ReturnsAllSizes(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	result := constants.GetMapSizes(true)

	// Assert
	assert.Equal(t, constants.AllMapSizes, result)
}
