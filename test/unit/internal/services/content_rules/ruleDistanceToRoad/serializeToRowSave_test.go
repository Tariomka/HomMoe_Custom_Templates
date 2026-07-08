package ruleDistanceToRoad_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
)

func TestWhenRuleIsSerialized_WritesNameAndDistanceName(t *testing.T) {
	// Arrange
	distance := content_rules.DistanceFar
	rule := content_rules.NewRuleDistanceToRoad(&distance)

	// Act
	saved := rule.SerializeToRowSave()

	// Assert
	assert.Equal(t, models.ContentRuleRowSave{
		Name:         "Distance to road",
		DistanceName: "Far",
	}, saved)
}
