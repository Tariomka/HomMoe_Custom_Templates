package placementRuleBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/placement_rule"
	"github.com/stretchr/testify/assert"
)

func TestWhenBuilderIsFreshlyCreated_ProducesEmptyRule(t *testing.T) {
	// Arrange & Act
	builder := placement_rule.NewPlacementRuleBuilder()

	// Assert
	assert.Equal(t, entities.PlacementRule{}, builder.Build())
}
