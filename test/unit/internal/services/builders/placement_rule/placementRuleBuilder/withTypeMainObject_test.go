package placementRuleBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/placement_rule"
	"github.com/stretchr/testify/assert"
)

func TestWhenMainObjectTypeIsChosen_SetsMainObjectTypeOnBuiltRule(t *testing.T) {
	// Arrange
	builder := placement_rule.NewPlacementRuleBuilder()

	// Act
	rule := builder.WithTypeMainObject().Build()

	// Assert
	assert.Equal(t, entities.PlacementRule{Type: "MainObject"}, rule)
}
