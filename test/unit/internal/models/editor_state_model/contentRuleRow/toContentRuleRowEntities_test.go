package contentRuleRow_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenThereAreNoRuleModels_ReturnsNil(t *testing.T) {
	t.Parallel()
	// Arrange
	var rules []editor_state_model.ContentRuleRow

	// Act
	entities := editor_state_model.ToContentRuleRowEntities(rules)

	// Assert
	assert.Nil(t, entities)
}

func TestWhenRuleModelsAreUnwrapped_TheEntitiesAreCarried(t *testing.T) {
	t.Parallel()
	// Arrange
	rules := []editor_state_model.ContentRuleRow{
		{Name: "Guarded", IsGuarded: new(true)},
		{Name: "Distance to road", DistanceName: "Next To"},
	}

	// Act
	entities := editor_state_model.ToContentRuleRowEntities(rules)

	// Assert
	assert.Equal(
		t,
		[]editor_state.ContentRuleRow{rules[0].ContentRuleRow, rules[1].ContentRuleRow},
		entities)
}
