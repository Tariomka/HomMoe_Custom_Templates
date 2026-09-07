package ruleSoloEncounter_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
)

func TestWhenSoloRuleIsApplied_SetsItemSoloEncounter(t *testing.T) {
	t.Parallel()
	// Arrange
	rule := content_rules.NewRuleSoloEncounter(true)
	item := template_model.MandatoryContentItem{SID: "x"}

	// Act
	rule.Apply(&item)

	// Assert
	assert.True(t, item.SoloEncounter)
}

func TestWhenNotSoloRuleIsApplied_ClearsItemSoloEncounter(t *testing.T) {
	t.Parallel()
	// Arrange
	rule := content_rules.NewRuleSoloEncounter(false)
	item := template_model.MandatoryContentItem{SID: "x", SoloEncounter: true}

	// Act
	rule.Apply(&item)

	// Assert
	assert.False(t, item.SoloEncounter)
}
