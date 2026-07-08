package mandatoryContentBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/mandatory_content"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/placement_rule"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenRoadDistanceIsProvided_AppendsRoadRuleWithWeightOneOnBuiltItem(t *testing.T) {
	// Arrange
	expectedSid := gofakeit.Word()
	expectedDistance := placement_rule.Distance{
		Min: gofakeit.Float64Range(0.01, 0.4),
		Max: gofakeit.Float64Range(0.5, 0.95),
	}
	builder := mandatory_content.NewContentBuilder(expectedSid)

	// Act
	item := builder.WithRoadDistance(expectedDistance).Build()

	// Assert
	assert.Equal(t, entities.MandatoryContentItem{
		SID: expectedSid,
		Rules: []entities.PlacementRule{{
			Type:      "Road",
			TargetMin: expectedDistance.Min,
			TargetMax: expectedDistance.Max,
			Weight:    1,
		}},
	}, item)
}
