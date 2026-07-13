package string_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/utils"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenFormatterIsInvoked_MatchesMultiplierString(t *testing.T) {
	t.Parallel()
	// Arrange
	value := gofakeit.Float32Range(0, 1)
	formatter := utils.MultiplierFormatter(0.5, 1.5)

	// Act
	result := formatter(value)

	// Assert
	assert.Equal(t, utils.MultiplierString(value, 0.5, 1.5), result)
}

func TestWhenSliderIsAtZero_FormatterReturnsBaseMultiplier(t *testing.T) {
	t.Parallel()
	// Arrange
	formatter := utils.MultiplierFormatter(0.5, 1.5)

	// Act
	result := formatter(0)

	// Assert
	assert.Equal(t, "x 0.50", result)
}
