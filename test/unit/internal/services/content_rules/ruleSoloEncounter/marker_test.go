package ruleSoloEncounter_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
)

func TestWhenSolo_ReturnsS(t *testing.T) {
	t.Parallel()
	// Arrange
	rule := content_rules.NewRuleSoloEncounter(true)

	// Act
	marker := rule.Marker()

	// Assert
	assert.Equal(t, "S", marker)
}

func TestWhenExplicitlyNotSolo_ReturnsNegatedS(t *testing.T) {
	t.Parallel()
	// Arrange
	rule := content_rules.NewRuleSoloEncounter(false)

	// Act
	marker := rule.Marker()

	// Assert
	assert.Equal(t, "!S", marker)
}
