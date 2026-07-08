package string_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/utils"
	"github.com/stretchr/testify/assert"
)

func TestWhenValueIsFormatted_PlusMinusPrefixAndTwoDecimalsAreUsed(t *testing.T) {
	// Arrange
	value := float32(0.5)

	// Act
	result := utils.DenormalizeString(value, 0, 10)

	// Assert
	assert.Equal(t, "± 5.00", result)
}

func TestWhenValueExceedsOne_ClampedHighBoundIsFormatted(t *testing.T) {
	// Arrange
	value := float32(1.5)

	// Act
	result := utils.DenormalizeString(value, 0, 10)

	// Assert
	assert.Equal(t, "± 10.00", result)
}
