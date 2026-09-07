package placementRuleBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/placement_rule"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenArgumentsAreProvidedTwice_AppendsAllArgumentsOnBuiltRule(t *testing.T) {
	t.Parallel()
	// Arrange
	firstArgument := gofakeit.Word()
	secondArgument := gofakeit.Word()
	thirdArgument := gofakeit.Number(1, 100)
	builder := placement_rule.NewPlacementRuleBuilder()

	// Act
	rule := builder.WithArgs(firstArgument, secondArgument).WithArgs(thirdArgument).Build()

	// Assert
	assert.Equal(t, template_model.PlacementRule{
		Args: []any{firstArgument, secondArgument, thirdArgument},
	}, rule)
}
