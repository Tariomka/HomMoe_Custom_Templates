package neutralZonePlan_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenPlanIsAdded_AppendsPlanWithGivenValues(t *testing.T) {
	t.Parallel()
	// Arrange
	label := gofakeit.LetterN(2)
	castleCount := gofakeit.Number(0, 4)
	plans := neutral_zone.Plans{}
	expected := neutral_zone.Plans{
		{Label: label, Quality: neutral_zone.QualityHigh, CastleCount: castleCount},
	}

	// Act
	plans.AddPlan(label, neutral_zone.QualityHigh, castleCount)

	// Assert
	assert.Equal(t, expected, plans)
}
