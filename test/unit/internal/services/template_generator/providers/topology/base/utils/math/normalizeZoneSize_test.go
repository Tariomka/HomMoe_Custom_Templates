package math_test

import (
	"math"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base/utils"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenZoneSizeIsProvided_NormalizesIntoSupportedRange(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name     string
		zoneSize float64
		expected float64
	}{
		{"WhenSizeIsNaN_ReturnsDefaultSize", math.NaN(), 1.0},
		{"WhenSizeIsPositiveInfinity_ReturnsDefaultSize", math.Inf(1), 1.0},
		{"WhenSizeIsNegativeInfinity_ReturnsDefaultSize", math.Inf(-1), 1.0},
		{"WhenSizeIsBelowMinimum_ClampsToMinimum", 0.01, 0.1},
		{"WhenSizeIsNegative_ClampsToMinimum", -3.5, 0.1},
		{"WhenSizeExceedsMaximum_ClampsToMaximum", 5.0, 2.0},
		{"WhenSizeIsWithinRange_KeepsValue", 0.75, 0.75},
		{"WhenSizeHasExtraPrecision_RoundsToTwoDecimals", 1.234567, 1.23},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange

			// Act
			normalized := utils.NormalizeZoneSize(testCase.zoneSize)

			// Assert
			assert.InDelta(t, testCase.expected, normalized, test_helpers.Delta)
		})
	}
}

func TestWhenArbitraryFiniteSizeIsProvided_ResultStaysWithinClampBounds(t *testing.T) {
	t.Parallel()
	// Arrange
	zoneSize := gofakeit.Float64Range(-10, 10)

	// Act
	normalized := utils.NormalizeZoneSize(zoneSize)

	// Assert
	assert.GreaterOrEqual(t, normalized, 0.1)
	assert.LessOrEqual(t, normalized, 2.0)
}
