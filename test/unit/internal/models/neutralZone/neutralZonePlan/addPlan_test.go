package neutralZonePlan_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutralZone"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenPlanIsAdded_AppendsPlanWithGivenValues(t *testing.T) {
	t.Parallel()
	// Arrange
	label := gofakeit.LetterN(2)
	castleCount := gofakeit.Number(0, 4)
	plans := neutralZone.Plans{}
	expected := neutralZone.Plans{
		{Label: label, Quality: neutralZone.QualityHigh, CastleCount: castleCount},
	}

	// Act
	plans.AddPlan(label, neutralZone.QualityHigh, castleCount)

	// Assert
	assert.Equal(t, expected, plans)
}
