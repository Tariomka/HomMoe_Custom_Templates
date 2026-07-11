package ruleSoloEncounter_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
)

func TestWhenDescriptionIsQueried_ExplainsSoloEncounterBehavior(t *testing.T) {
	t.Parallel()
	// Arrange
	rule := content_rules.NewRuleSoloEncounter(true)

	// Act
	description := rule.Description()

	// Assert
	assert.Equal(
		t,
		"Solo encounter means that the content item will be spawned without any additional content items around, enforcing consistent guard strength. Setting to false will make it more likely to be spawned with other content items, but will not always guarantee it.",
		description,
	)
}
