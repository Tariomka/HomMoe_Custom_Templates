package string_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/utils"
	"github.com/stretchr/testify/assert"
)

func TestWhenSnappedValueIsFormatted_PercentSuffixIsAppended(t *testing.T) {
	t.Parallel()
	// Arrange
	value := float32(0.5)

	// Act
	result := utils.RoundedRangePercentString(value, 0, 100)

	// Assert
	assert.Equal(t, "50%", result)
}

func TestWhenSliderIsAtZero_LowBoundPercentIsReturned(t *testing.T) {
	t.Parallel()
	// Arrange
	value := float32(0)

	// Act
	result := utils.RoundedRangePercentString(value, 25, 200)

	// Assert
	assert.Equal(t, "25%", result)
}
