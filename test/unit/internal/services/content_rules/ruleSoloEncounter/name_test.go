package ruleSoloEncounter_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
)

func TestWhenNameIsQueried_ReturnsSoloEncounter(t *testing.T) {
	// Arrange
	rule := content_rules.NewRuleSoloEncounter(true)

	// Act
	name := rule.Name()

	// Assert
	assert.Equal(t, "Solo Encounter", name)
}
