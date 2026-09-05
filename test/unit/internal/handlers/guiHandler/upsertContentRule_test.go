package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenARuleIsUpsertedIntoAnEmptyList_ItIsAppended(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	rule := editor_state_model.ContentRuleRow{Name: gofakeit.Word()}

	// Act
	rules := handler.UpsertContentRule(nil, rule)

	// Assert
	assert.Equal(t, []editor_state_model.ContentRuleRow{rule}, rules)
}
