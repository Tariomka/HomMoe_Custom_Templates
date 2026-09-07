package math_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/stretchr/testify/assert"
)

func TestWhenTIsZero_ReturnsTheCurveStart(t *testing.T) {
	t.Parallel()
	// Arrange
	start, ctrl, end := data.NewVec2(10.0, 20.0), data.NewVec2(50.0, 0.0), data.NewVec2(90.0, 20.0)

	// Act
	point := helpers.GetVectorOnQuadraticBezierCurve(start, ctrl, end, 0)

	// Assert
	assert.InDeltaSlice(t, []float64{10.0, 20.0}, []float64{point.X, point.Y}, 1e-9)
}

func TestWhenTIsOne_ReturnsTheCurveEnd(t *testing.T) {
	t.Parallel()
	// Arrange
	start, ctrl, end := data.NewVec2(10.0, 20.0), data.NewVec2(50.0, 0.0), data.NewVec2(90.0, 20.0)

	// Act
	point := helpers.GetVectorOnQuadraticBezierCurve(start, ctrl, end, 1)

	// Assert
	assert.InDeltaSlice(t, []float64{90.0, 20.0}, []float64{point.X, point.Y}, 1e-9)
}

func TestWhenTIsAHalf_ReturnsTheQuarterWeightedMidpoint(t *testing.T) {
	t.Parallel()
	// Arrange
	start, ctrl, end := data.NewVec2(10.0, 20.0), data.NewVec2(50.0, 0.0), data.NewVec2(90.0, 20.0)

	// Act
	point := helpers.GetVectorOnQuadraticBezierCurve(start, ctrl, end, 0.5)

	// Assert
	assert.InDeltaSlice(t, []float64{50.0, 10.0}, []float64{point.X, point.Y}, 1e-9)
}
