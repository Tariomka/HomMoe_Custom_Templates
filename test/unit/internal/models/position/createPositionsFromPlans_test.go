package position_test

import (
	"math"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestWhenNoLabelsAreProvided_ReturnsNil(t *testing.T) {
	t.Parallel()
	// Arrange
	var orderedLabels []string

	// Act
	positions := models.CreatePositionsFromPlans(orderedLabels, nil, models.NeutralZonePlans{})

	// Assert
	assert.Nil(t, positions)
}

func TestWhenSinglePlayerLabelIsProvided_PlacesItOnPlayerRingRadius(t *testing.T) {
	t.Parallel()
	// Arrange
	orderedLabels := []string{"P1"}
	playerLabels := []string{"P1"}
	expected := models.Positions{data.NewVec2(
		0.5+math.Cos(-0.008)*0.38,
		0.5+math.Sin(-0.008)*0.38,
	)}

	// Act
	positions := models.CreatePositionsFromPlans(orderedLabels, playerLabels, models.NeutralZonePlans{})

	// Assert
	assert.Equal(t, expected, positions)
}

func TestWhenLabelsSpanEveryTier_ReturnsOnePositionPerLabel(t *testing.T) {
	t.Parallel()
	// Arrange
	orderedLabels := []string{"P1", "P2", "L1", "M1", "H1"}
	playerLabels := []string{"P1", "P2"}
	plans := models.NeutralZonePlans{
		{Label: "L1", Quality: models.QualityLow},
		{Label: "M1", Quality: models.QualityMedium},
		{Label: "H1", Quality: models.QualityHigh},
	}

	// Act
	positions := models.CreatePositionsFromPlans(orderedLabels, playerLabels, plans)

	// Assert
	assert.Len(t, positions, 5)
}

func TestWhenManyLabelsArePlaced_ClampsEveryPositionInsideCanvasMargins(t *testing.T) {
	t.Parallel()
	// Arrange
	orderedLabels := []string{"P1", "P2", "P3", "L1", "L2", "M1", "M2", "H1", "H2"}
	playerLabels := []string{"P1", "P2", "P3"}
	plans := models.NeutralZonePlans{
		{Label: "L1", Quality: models.QualityLow},
		{Label: "L2", Quality: models.QualityLow},
		{Label: "M1", Quality: models.QualityMedium},
		{Label: "M2", Quality: models.QualityMedium},
		{Label: "H1", Quality: models.QualityHigh},
		{Label: "H2", Quality: models.QualityHigh},
	}

	// Act
	positions := models.CreatePositionsFromPlans(orderedLabels, playerLabels, plans)

	// Assert
	outOfBoundsCount := 0
	for _, position := range positions {
		if position.X < 0.05 || position.X > 0.95 || position.Y < 0.05 || position.Y > 0.95 {
			outOfBoundsCount++
		}
	}
	assert.Zero(t, outOfBoundsCount)
}
