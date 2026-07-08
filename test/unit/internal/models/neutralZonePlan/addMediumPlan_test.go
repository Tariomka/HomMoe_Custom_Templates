package neutralZonePlan_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenMediumPlanIsAdded_AppendsPlanWithMediumQuality(t *testing.T) {
	// Arrange
	label := gofakeit.LetterN(2)
	castleCount := gofakeit.Number(0, 4)
	plans := models.NeutralZonePlans{}
	expected := models.NeutralZonePlans{
		{Label: label, Quality: models.QualityMedium, CastleCount: castleCount},
	}

	// Act
	plans.AddMediumPlan(label, castleCount)

	// Assert
	assert.Equal(t, expected, plans)
}
