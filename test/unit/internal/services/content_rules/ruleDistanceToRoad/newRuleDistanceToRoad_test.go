package ruleDistanceToRoad_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
)

func TestWhenDistanceIsNil_DefaultsToMedium(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	rule := content_rules.NewRuleDistanceToRoad(nil)

	// Assert
	assert.Equal(t, models.DistancePreset{Name: "Medium", Min: 0.25, Max: 0.5}, rule.Distance)
}

func TestWhenDistanceIsSupplied_UsesIt(t *testing.T) {
	t.Parallel()
	// Arrange
	distance := models.DistancePreset{Name: "Far", Min: 0.5, Max: 0.75}

	// Act
	rule := content_rules.NewRuleDistanceToRoad(&distance)

	// Assert
	assert.Equal(t, distance, rule.Distance)
}
