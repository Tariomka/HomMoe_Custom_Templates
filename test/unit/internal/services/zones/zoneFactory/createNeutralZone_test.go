package zoneFactory_test

import (
	"math"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/stretchr/testify/assert"
)

func TestWhenNeutralZoneIsCreated_PreservesExplicitName(t *testing.T) {
	t.Parallel()
	// Arrange
	factory := newZoneFactory()
	input := models.NeutralZoneCreationRequest{
		Name:               "Neutral-Q",
		Quality:            neutral_zone.QualityMedium,
		Size:               1,
		GuardRandomization: 0.05,
		Tuning:             newUnitTuning(),
	}

	// Act
	zone := factory.CreateNeutralZone(input)

	// Assert
	assert.Equal(t, "Neutral-Q", zone.Name)
}

func TestWhenNeutralZoneIsCreated_RecordsTheRequestedQuality(t *testing.T) {
	t.Parallel()
	// Arrange
	factory := newZoneFactory()
	input := models.NeutralZoneCreationRequest{
		Name:    "Neutral-Q",
		Quality: neutral_zone.QualityLow,
		Size:    1,
		Tuning:  newUnitTuning(),
	}

	// Act
	zone := factory.CreateNeutralZone(input)

	// Assert
	assert.Equal(t, new(neutral_zone.QualityLow), zone.Quality)
}

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
			factory := newZoneFactory()
			input := models.NeutralZoneCreationRequest{
				Name:    "Neutral-Q",
				Quality: neutral_zone.QualityMedium,
				Size:    testCase.zoneSize,
				Tuning:  newUnitTuning(),
			}

			// Act
			zone := factory.CreateNeutralZone(input)

			// Assert
			assert.InDelta(t, testCase.expected, zone.Size, 1e-9)
		})
	}
}
