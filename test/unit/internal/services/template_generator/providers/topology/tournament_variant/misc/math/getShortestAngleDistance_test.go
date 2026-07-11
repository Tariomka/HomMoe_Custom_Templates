package math_test

import (
	"math"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/tournament_variant/misc"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenAnglesDiffer_ReturnsExpectedShortestDistance(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name     string
		from     float64
		to       float64
		expected float64
	}{
		{"WhenAnglesAreEqual_ReturnsZero", 1.25, 1.25, 0},
		{"WhenAnglesAreQuarterTurnApart_ReturnsHalfPi", 0, math.Pi / 2, math.Pi / 2},
		{"WhenAnglesAreHalfTurnApart_ReturnsPi", 0, math.Pi, math.Pi},
		{"WhenAnglesWrapAroundFullTurn_ReturnsShortWayDistance", 0.1, 2*math.Pi - 0.1, 0.2},
		{"WhenRawDeltaExceedsPi_ReturnsComplementDistance", 0, 3 * math.Pi / 2, math.Pi / 2},
		{"WhenAnglesAreSeveralTurnsApart_WrapsIntoSingleTurn", 0, 4*math.Pi + 0.3, 0.3},
		{"WhenFromIsNegative_MeasuresAcrossZero", -0.2, 0.3, 0.5},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange

			// Act
			distance := misc.GetShortestAngleDistance(testCase.from, testCase.to)

			// Assert
			assert.InDelta(t, testCase.expected, distance, 1e-12)
		})
	}
}

func TestWhenArbitraryAnglesProvided_DistanceIsNeverNegative(t *testing.T) {
	t.Parallel()
	// Arrange
	fromAngle := gofakeit.Float64Range(-100, 100)
	toAngle := gofakeit.Float64Range(-100, 100)

	// Act
	distance := misc.GetShortestAngleDistance(fromAngle, toAngle)

	// Assert
	assert.GreaterOrEqual(t, distance, 0.0)
}

func TestWhenArbitraryAnglesProvided_DistanceNeverExceedsPi(t *testing.T) {
	t.Parallel()
	// Arrange
	fromAngle := gofakeit.Float64Range(-100, 100)
	toAngle := gofakeit.Float64Range(-100, 100)

	// Act
	distance := misc.GetShortestAngleDistance(fromAngle, toAngle)

	// Assert
	assert.LessOrEqual(t, distance, math.Pi)
}

func TestWhenArgumentsAreSwapped_DistanceIsSymmetric(t *testing.T) {
	t.Parallel()
	// Arrange
	fromAngle := gofakeit.Float64Range(-10, 10)
	toAngle := gofakeit.Float64Range(-10, 10)

	// Act
	forward := misc.GetShortestAngleDistance(fromAngle, toAngle)
	backward := misc.GetShortestAngleDistance(toAngle, fromAngle)

	// Assert
	assert.InDelta(t, forward, backward, 1e-12)
}
