package mapSizes_test

import (
	"testing"

	gui_constants "github.com/Tariomka/hommoe_custom_templates/app/gui/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/common"
	"github.com/stretchr/testify/assert"
)

func TestWhenExperimentalIsFalse_ReturnsSameBaseSizesAsInternalConstants(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	result := gui_constants.GetMapSizes(false)

	// Assert
	assert.Equal(t, common.GetMapSizes(false), result)
}

func TestWhenExperimentalIsTrue_ReturnsSameAllSizesAsInternalConstants(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	result := gui_constants.GetMapSizes(true)

	// Assert
	assert.Equal(t, common.GetMapSizes(true), result)
}
