package string_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/utils"
	"github.com/stretchr/testify/assert"
)

func TestWhenMultiplierIsFormatted_XPrefixAndTwoDecimalsAreUsed(t *testing.T) {
	t.Parallel()
	// Arrange
	value := float32(0.5)

	// Act
	result := utils.MultiplierString(value, 1, 1)

	// Assert
	assert.Equal(t, "x 1.50", result)
}

func TestWhenFactorIsZero_BaseIsFormatted(t *testing.T) {
	t.Parallel()
	// Arrange
	value := float32(0.75)

	// Act
	result := utils.MultiplierString(value, 2, 0)

	// Assert
	assert.Equal(t, "x 2.00", result)
}
