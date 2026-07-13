package string_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/utils"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenFormatterIsInvoked_MatchesRoundedRangeString(t *testing.T) {
	t.Parallel()
	// Arrange
	value := gofakeit.Float32Range(0, 1)
	formatter := utils.RoundedRangeFormatter(1, 32)

	// Act
	result := formatter(value)

	// Assert
	assert.Equal(t, utils.RoundedRangeString(value, 1, 32), result)
}

func TestWhenSliderIsAtZero_FormatterReturnsLowBound(t *testing.T) {
	t.Parallel()
	// Arrange
	formatter := utils.RoundedRangeFormatter(3, 30)

	// Act
	result := formatter(0)

	// Assert
	assert.Equal(t, "3 ", result)
}
