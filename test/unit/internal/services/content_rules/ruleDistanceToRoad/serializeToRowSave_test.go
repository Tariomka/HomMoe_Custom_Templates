package ruleDistanceToRoad_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
)

func TestWhenRuleIsSerialized_WritesNameAndDistanceName(t *testing.T) {
	t.Parallel()
	// Arrange
	distance := models.DistancePreset{Name: "Far", Min: 0.5, Max: 0.75}
	rule := content_rules.NewRuleDistanceToRoad(&distance)

	// Act
	saved := rule.SerializeToRowSave()

	// Assert
	assert.Equal(t, editor_state_model.ContentRuleRow{
		Name:         "Distance to road",
		DistanceName: "Far",
	}, saved)
}
