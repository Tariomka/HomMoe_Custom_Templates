package ruleDistanceToTown_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
)

func TestWhenRuleIsSerialized_WritesNameAndDistanceName(t *testing.T) {
	t.Parallel()
	// Arrange
	distance := models.DistancePreset{Name: "Near", Min: 0.1, Max: 0.25}
	rule := content_rules.NewRuleDistanceToTown(&distance)

	// Act
	saved := rule.SerializeToRowSave()

	// Assert
	assert.Equal(t, models.ContentRuleRow{
		Name:         "Distance to town",
		DistanceName: "Near",
	}, saved)
}
