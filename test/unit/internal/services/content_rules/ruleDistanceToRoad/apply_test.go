package ruleDistanceToRoad_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
)

func TestWhenRuleIsApplied_AppendsRoadPlacementRuleWithDistanceBounds(t *testing.T) {
	t.Parallel()
	// Arrange
	distance := models.DistancePreset{Name: "Near", Min: 0.1, Max: 0.25}
	rule := content_rules.NewRuleDistanceToRoad(&distance)
	item := template_model.MandatoryContentItem{SID: "x"}

	// Act
	rule.Apply(&item)

	// Assert
	assert.Equal(t, []template_model.PlacementRule{
		{Type: "Road", TargetMin: 0.1, TargetMax: 0.25, Weight: 1},
	}, item.Rules)
}

func TestWhenItemAlreadyHasRules_AppendsAfterExisting(t *testing.T) {
	t.Parallel()
	// Arrange
	rule := content_rules.NewRuleDistanceToRoad(nil)
	item := template_model.MandatoryContentItem{
		SID:   "x",
		Rules: []template_model.PlacementRule{{Type: "Crossroads", Weight: 2}},
	}

	// Act
	rule.Apply(&item)

	// Assert
	assert.Len(t, item.Rules, 2)
}
