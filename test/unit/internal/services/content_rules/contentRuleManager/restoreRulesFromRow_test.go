package contentRuleManager_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
)

func TestWhenRowHasSerializedRules_RestoresEachRule(t *testing.T) {
	// Arrange
	isGuarded := true
	row := models.ZoneContentRowSave{
		Sid: "x",
		Rules: []models.ContentRuleRowSave{
			{Name: "Guarded", IsGuarded: &isGuarded},
			{Name: "Distance to road", DistanceName: "Far"},
		},
	}

	// Act
	rules := content_rules.RestoreRulesFromRow(row, models.SidMapping{Sid: "x"})

	// Assert
	ruleNames := make([]string, 0, len(rules))
	for _, rule := range rules {
		ruleNames = append(ruleNames, rule.Name())
	}
	assert.Equal(t, []string{"Guarded", "Distance to road"}, ruleNames)
}

func TestWhenRowHasNoRules_ReturnsNoRules(t *testing.T) {
	// Arrange
	row := models.ZoneContentRowSave{Sid: "x"}

	// Act
	rules := content_rules.RestoreRulesFromRow(row, models.SidMapping{Sid: "x"})

	// Assert
	assert.Empty(t, rules)
}

func TestWhenSavedRuleIsInvalid_SkipsIt(t *testing.T) {
	// Arrange
	isGuarded := true
	row := models.ZoneContentRowSave{
		Sid: "x",
		Rules: []models.ContentRuleRowSave{
			{Name: "Nope"},
			{Name: "Guarded", IsGuarded: &isGuarded},
		},
	}

	// Act
	rules := content_rules.RestoreRulesFromRow(row, models.SidMapping{Sid: "x"})

	// Assert
	assert.Len(t, rules, 1)
}
