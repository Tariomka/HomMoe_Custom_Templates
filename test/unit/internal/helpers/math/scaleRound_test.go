package math_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenScaledValueHasAFractionAboveAHalf_RoundsUp(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	scaled := helpers.ScaleRound(10, 1.06)

	// Assert
	assert.Equal(t, 11, scaled)
}

func TestWhenScaledValueHasAFractionBelowAHalf_RoundsDown(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	scaled := helpers.ScaleRound(10, 1.04)

	// Assert
	assert.Equal(t, 10, scaled)
}

func TestWhenScaledValueIsNegative_ClampsToZero(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	scaled := helpers.ScaleRound(10, -2)

	// Assert
	assert.Zero(t, scaled)
}
