package contentRuleManager_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenRoadDistanceRuleIsApplied_AppendsRoadPlacementRule(t *testing.T) {
	// Arrange
	item := entities.MandatoryContentItem{SID: "x"}
	near := content_rules.DistanceNear

	// Act
	content_rules.ApplyRulesToItem(&item, []content_rules.ContentRule{
		content_rules.NewRuleDistanceToRoad(&near),
	})

	// Assert
	assert.Equal(t, []entities.PlacementRule{
		{Type: "Road", TargetMin: 0.1, TargetMax: 0.25, Weight: 1},
	}, item.Rules)
}

func TestWhenTownDistanceRuleIsApplied_AppendsMainObjectPlacementRule(t *testing.T) {
	// Arrange
	item := entities.MandatoryContentItem{SID: "x"}
	near := content_rules.DistanceNear

	// Act
	content_rules.ApplyRulesToItem(&item, []content_rules.ContentRule{
		content_rules.NewRuleDistanceToTown(&near),
	})

	// Assert
	assert.Equal(t, []entities.PlacementRule{
		{Type: "MainObject", Args: []any{"0"}, TargetMin: 0.1, TargetMax: 0.25, Weight: 1},
	}, item.Rules)
}

func TestWhenGuardedRuleIsApplied_SetsIsGuarded(t *testing.T) {
	// Arrange
	item := entities.MandatoryContentItem{SID: "x"}

	// Act
	content_rules.ApplyRulesToItem(&item, []content_rules.ContentRule{
		content_rules.NewRuleGuarded(true),
	})

	// Assert
	assert.True(t, item.IsGuarded)
}

func TestWhenGuardedRuleIsApplied_AddsNoPlacementRules(t *testing.T) {
	// Arrange
	item := entities.MandatoryContentItem{SID: "x"}

	// Act
	content_rules.ApplyRulesToItem(&item, []content_rules.ContentRule{
		content_rules.NewRuleGuarded(true),
	})

	// Assert
	assert.Empty(t, item.Rules)
}

func TestWhenSoloEncounterRuleIsApplied_SetsSoloEncounter(t *testing.T) {
	// Arrange
	item := entities.MandatoryContentItem{SID: "x"}

	// Act
	content_rules.ApplyRulesToItem(&item, []content_rules.ContentRule{
		content_rules.NewRuleSoloEncounter(true),
	})

	// Assert
	assert.True(t, item.SoloEncounter)
}

func TestWhenVariantRuleIsApplied_SetsVariantId(t *testing.T) {
	// Arrange
	item := entities.MandatoryContentItem{SID: "dragon_utopia"}
	variantId := 3
	variantRule, err := content_rules.NewRuleVariant(&content_rules.UtopiaVariants, &variantId)
	require.NoError(t, err)

	// Act
	content_rules.ApplyRulesToItem(&item, []content_rules.ContentRule{variantRule})

	// Assert
	require.NotNil(t, item.Variant)
	assert.Equal(t, 3, *item.Variant)
}

func TestWhenRuleListContainsNil_SkipsItWithoutPanicking(t *testing.T) {
	// Arrange
	item := entities.MandatoryContentItem{SID: "x"}

	// Act
	content_rules.ApplyRulesToItem(&item, []content_rules.ContentRule{
		nil,
		content_rules.NewRuleGuarded(true),
	})

	// Assert
	assert.True(t, item.IsGuarded)
}
