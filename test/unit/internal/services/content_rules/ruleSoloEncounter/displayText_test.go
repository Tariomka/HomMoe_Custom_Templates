package ruleSoloEncounter_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
)

func TestWhenSolo_ShowsTrueState(t *testing.T) {
	t.Parallel()
	// Arrange
	rule := content_rules.NewRuleSoloEncounter(true)

	// Act
	displayText := rule.DisplayText()

	// Assert
	assert.Equal(t, "Solo Encounter: true", displayText)
}

func TestWhenExplicitlyNotSolo_ShowsFalseState(t *testing.T) {
	t.Parallel()
	// Arrange
	rule := content_rules.NewRuleSoloEncounter(false)

	// Act
	displayText := rule.DisplayText()

	// Assert
	assert.Equal(t, "Solo Encounter: false", displayText)
}
