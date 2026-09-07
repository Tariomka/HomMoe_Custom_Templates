package ruleSoloEncounter_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
)

func TestWhenSoloRuleIsSerialized_WritesNameAndTrueState(t *testing.T) {
	t.Parallel()
	// Arrange
	rule := content_rules.NewRuleSoloEncounter(true)
	isSoloEncounter := true

	// Act
	saved := rule.SerializeToRowSave()

	// Assert
	assert.Equal(t, editor_state_model.ContentRuleRow{Name: "Solo Encounter", IsSoloEncounter: &isSoloEncounter}, saved)
}

func TestWhenNotSoloRuleIsSerialized_WritesNameAndFalseState(t *testing.T) {
	t.Parallel()
	// Arrange
	rule := content_rules.NewRuleSoloEncounter(false)
	isSoloEncounter := false

	// Act
	saved := rule.SerializeToRowSave()

	// Assert
	assert.Equal(t, editor_state_model.ContentRuleRow{Name: "Solo Encounter", IsSoloEncounter: &isSoloEncounter}, saved)
}
