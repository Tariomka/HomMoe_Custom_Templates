package string_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/utils"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenFormatterIsInvoked_MatchesDenormalizeString(t *testing.T) {
	t.Parallel()
	// Arrange
	value := gofakeit.Float32Range(0, 1)
	formatter := utils.DenormalizeFormatter(0, 0.5)

	// Act
	result := formatter(value)

	// Assert
	assert.Equal(t, utils.DenormalizeString(value, 0, 0.5), result)
}

func TestWhenSliderIsAtOne_FormatterReturnsHighBoundWithPlusMinusPrefix(t *testing.T) {
	t.Parallel()
	// Arrange
	formatter := utils.DenormalizeFormatter(0, 0.5)

	// Act
	result := formatter(1)

	// Assert
	assert.Equal(t, "± 0.50", result)
}
