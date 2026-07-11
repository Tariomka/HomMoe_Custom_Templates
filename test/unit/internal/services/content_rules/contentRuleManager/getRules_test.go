package contentRuleManager_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
)

func TestWhenRulesAreListed_ReturnsEveryKnownRuleInManagerOrder(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := []string{
		"Distance to road",
		"Distance to town",
		"Guarded",
		"Variant",
		"Solo Encounter",
	}

	// Act
	rules := content_rules.GetRules()

	// Assert
	ruleNames := make([]string, 0, len(rules))
	for _, rule := range rules {
		ruleNames = append(ruleNames, rule.Name())
	}
	assert.Equal(t, expected, ruleNames)
}
