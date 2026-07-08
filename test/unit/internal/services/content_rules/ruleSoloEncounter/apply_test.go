package ruleSoloEncounter_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
)

func TestWhenSoloRuleIsApplied_SetsItemSoloEncounter(t *testing.T) {
	// Arrange
	rule := content_rules.NewRuleSoloEncounter(true)
	item := entities.MandatoryContentItem{SID: "x"}

	// Act
	rule.Apply(&item)

	// Assert
	assert.True(t, item.SoloEncounter)
}

func TestWhenNotSoloRuleIsApplied_ClearsItemSoloEncounter(t *testing.T) {
	// Arrange
	rule := content_rules.NewRuleSoloEncounter(false)
	item := entities.MandatoryContentItem{SID: "x", SoloEncounter: true}

	// Act
	rule.Apply(&item)

	// Assert
	assert.False(t, item.SoloEncounter)
}
