package mapSizes_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common"
	"github.com/stretchr/testify/assert"
)

func TestWhenExperimentalIsFalse_ReturnsBaseSizesOnly(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	result := common.GetMapSizes(false)

	// Assert
	assert.Equal(t, common.BaseMapSizes, result)
}

func TestWhenExperimentalIsTrue_ReturnsAllSizes(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	result := common.GetMapSizes(true)

	// Assert
	assert.Equal(t, common.AllMapSizes, result)
}
