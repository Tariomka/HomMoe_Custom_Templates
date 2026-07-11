package placementRuleBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/placement_rule"
	"github.com/stretchr/testify/assert"
)

func TestWhenCrossroadsTypeIsChosen_SetsCrossroadsTypeOnBuiltRule(t *testing.T) {
	t.Parallel()
	// Arrange
	builder := placement_rule.NewPlacementRuleBuilder()

	// Act
	rule := builder.WithTypeCrossroads().Build()

	// Assert
	assert.Equal(t, entities.PlacementRule{Type: "Crossroads"}, rule)
}
