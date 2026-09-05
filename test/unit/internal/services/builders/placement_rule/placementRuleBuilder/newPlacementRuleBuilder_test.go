package placementRuleBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/placement_rule"
	"github.com/stretchr/testify/assert"
)

func TestWhenBuilderIsFreshlyCreated_ProducesEmptyRule(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	builder := placement_rule.NewPlacementRuleBuilder()

	// Assert
	assert.Equal(t, template_model.PlacementRule{}, builder.Build())
}
