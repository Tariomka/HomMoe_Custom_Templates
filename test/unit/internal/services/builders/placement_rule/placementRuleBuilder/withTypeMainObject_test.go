package placementRuleBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/placement_rule"
	"github.com/stretchr/testify/assert"
)

func TestWhenMainObjectTypeIsChosen_SetsMainObjectTypeOnBuiltRule(t *testing.T) {
	t.Parallel()
	// Arrange
	builder := placement_rule.NewPlacementRuleBuilder()

	// Act
	rule := builder.WithTypeMainObject().Build()

	// Assert
	assert.Equal(t, template_model.PlacementRule{Type: "MainObject"}, rule)
}
