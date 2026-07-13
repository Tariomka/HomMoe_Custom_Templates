package string_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/utils"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenFormatterIsInvoked_MatchesRoundedRangePercentString(t *testing.T) {
	t.Parallel()
	// Arrange
	value := gofakeit.Float32Range(0, 1)
	formatter := utils.RoundedRangePercentFormatter(25, 200)

	// Act
	result := formatter(value)

	// Assert
	assert.Equal(t, utils.RoundedRangePercentString(value, 25, 200), result)
}

func TestWhenSliderIsAtOne_FormatterReturnsHighBoundPercent(t *testing.T) {
	t.Parallel()
	// Arrange
	formatter := utils.RoundedRangePercentFormatter(25, 200)

	// Act
	result := formatter(1)

	// Assert
	assert.Equal(t, "200%", result)
}
