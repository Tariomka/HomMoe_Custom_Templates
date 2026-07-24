package ruleDistanceToTown_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
)

func TestWhenDistanceIsNil_DefaultsToMedium(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	rule := content_rules.NewRuleDistanceToTown(nil)

	// Assert
	assert.Equal(t, content_rules.DistanceVariation{Name: "Medium", Min: 0.25, Max: 0.5}, rule.Distance)
}

func TestWhenDistanceIsSupplied_UsesIt(t *testing.T) {
	t.Parallel()
	// Arrange
	distance := content_rules.DistanceVariation{Name: "Near", Min: 0.1, Max: 0.25}

	// Act
	rule := content_rules.NewRuleDistanceToTown(&distance)

	// Assert
	assert.Equal(t, distance, rule.Distance)
}
