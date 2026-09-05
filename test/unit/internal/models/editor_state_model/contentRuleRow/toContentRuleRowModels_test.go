package contentRuleRow_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenThereAreNoRules_ReturnsNil(t *testing.T) {
	t.Parallel()
	// Arrange
	var rules []editor_state.ContentRuleRow

	// Act
	models := editor_state_model.ToContentRuleRowModels(rules)

	// Assert
	assert.Nil(t, models)
}

func TestWhenRulesArePersisted_EachOneIsWrapped(t *testing.T) {
	t.Parallel()
	// Arrange
	rules := []editor_state.ContentRuleRow{
		{Name: "Guarded", IsGuarded: new(true)},
		{Name: "Distance to town", DistanceName: "Near"},
	}

	// Act
	models := editor_state_model.ToContentRuleRowModels(rules)

	// Assert
	assert.Equal(
		t,
		[]editor_state_model.ContentRuleRow{{ContentRuleRow: rules[0]}, {ContentRuleRow: rules[1]}},
		models)
}
