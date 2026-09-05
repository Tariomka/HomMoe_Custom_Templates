package placementRuleBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/placement_rule"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenDistanceIsProvided_SetsTargetRangeOnBuiltRule(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedDistance := models.DistancePreset{
		Min: gofakeit.Float64Range(0.01, 0.4),
		Max: gofakeit.Float64Range(0.5, 0.95),
	}
	builder := placement_rule.NewPlacementRuleBuilder()

	// Act
	rule := builder.WithDistance(expectedDistance).Build()

	// Assert
	assert.Equal(t, template_model.PlacementRule{
		TargetMin: expectedDistance.Min,
		TargetMax: expectedDistance.Max,
	}, rule)
}
