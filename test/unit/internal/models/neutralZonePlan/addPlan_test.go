package neutralZonePlan_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenPlanIsAdded_AppendsPlanWithGivenValues(t *testing.T) {
	// Arrange
	label := gofakeit.LetterN(2)
	castleCount := gofakeit.Number(0, 4)
	plans := models.NeutralZonePlans{}
	expected := models.NeutralZonePlans{
		{Label: label, Quality: models.QualityHigh, CastleCount: castleCount},
	}

	// Act
	plans.AddPlan(label, models.QualityHigh, castleCount)

	// Assert
	assert.Equal(t, expected, plans)
}
