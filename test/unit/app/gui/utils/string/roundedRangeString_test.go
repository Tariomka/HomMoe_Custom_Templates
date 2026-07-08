package string_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/utils"
	"github.com/stretchr/testify/assert"
)

func TestWhenSnappedValueIsFormatted_IntegerWithTrailingSpaceIsReturned(t *testing.T) {
	// Arrange
	value := float32(0.5)

	// Act
	result := utils.RoundedRangeString(value, 0, 10)

	// Assert
	assert.Equal(t, "5 ", result)
}

func TestWhenSliderIsAtOne_HighBoundWithTrailingSpaceIsReturned(t *testing.T) {
	// Arrange
	value := float32(1)

	// Act
	result := utils.RoundedRangeString(value, 3, 7)

	// Assert
	assert.Equal(t, "7 ", result)
}
