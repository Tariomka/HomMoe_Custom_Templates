package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
)

func TestWhenContentHasNoVariants_ReturnsBaseOptionsInDialogOrder(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := handlers.NewGuiHandler()
	content := models.SidMapping{Sid: "content_without_variants", Name: "Content"}
	expected := dtos.ContentRuleEditorOptionsDto{
		Rules: []dtos.ContentRuleOptionDto{
			{
				Key:         dtos.ContentRuleKeyDistanceToRoad,
				Name:        content_rules.RuleDistanceToRoadName,
				Description: content_rules.RuleDistanceToRoadDescription,
				Marker:      content_rules.RuleDistanceToRoadMarker,
				EditorKind:  dtos.ContentRuleEditorKindDistance,
				EditorLabel: "Distance",
			},
			{
				Key:         dtos.ContentRuleKeyDistanceToTown,
				Name:        content_rules.RuleDistanceToTownName,
				Description: content_rules.RuleDistanceToTownDescription,
				Marker:      content_rules.RuleDistanceToTownMarker,
				EditorKind:  dtos.ContentRuleEditorKindDistance,
				EditorLabel: "Distance",
			},
			{
				Key:         dtos.ContentRuleKeyGuarded,
				Name:        content_rules.RuleGuardedName,
				Description: content_rules.RuleGuardedDescription,
				Marker:      content_rules.RuleGuardedMarker,
				EditorKind:  dtos.ContentRuleEditorKindBoolean,
				EditorLabel: "Guarded",
			},
			{
				Key:         dtos.ContentRuleKeySoloEncounter,
				Name:        content_rules.RuleSoloEncounterName,
				Description: content_rules.RuleSoloEncounterDescription,
				Marker:      content_rules.RuleSoloEncounterMarker,
				EditorKind:  dtos.ContentRuleEditorKindBoolean,
				EditorLabel: "Solo encounter",
			},
		},
		Distances: []string{"Next To", "Near", "Medium", "Far", "Very Far"},
		Variants:  []dtos.ContentRuleVariantOptionDto{},
	}

	// Act
	result := handler.GetContentRuleEditorOptions(content)

	// Assert
	assert.Equal(t, expected, result)
}

func TestWhenContentHasVariants_AppendsVariantRule(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := handlers.NewGuiHandler()
	content := models.SidMapping{
		Sid:  registry.GetMapObjectT3GuardedResourceBankValues().DragonUtopia,
		Name: "Dragon Utopia",
	}
	expectedRule := dtos.ContentRuleOptionDto{
		Key:         dtos.ContentRuleKeyVariant,
		Name:        content_rules.RuleVariantName,
		Description: content_rules.RuleVariantDescription,
		Marker:      content_rules.RuleVariantMarker,
		EditorKind:  dtos.ContentRuleEditorKindVariant,
		EditorLabel: "Variant",
	}

	// Act
	result := handler.GetContentRuleEditorOptions(content)

	// Assert
	assert.Equal(t, expectedRule, result.Rules[len(result.Rules)-1])
}

func TestWhenContentHasVariants_ReturnsOptionsInVariantIdOrder(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := handlers.NewGuiHandler()
	content := models.SidMapping{
		Sid:  registry.GetMapObjectT3GuardedResourceBankValues().DragonUtopia,
		Name: "Dragon Utopia",
	}
	expectedVariants := []dtos.ContentRuleVariantOptionDto{
		{ID: 0, Label: "Small Guard"},
		{ID: 1, Label: "Medium Guard"},
		{ID: 2, Label: "Large Guard"},
		{ID: 3, Label: "Maximum Guard"},
	}

	// Act
	result := handler.GetContentRuleEditorOptions(content)

	// Assert
	assert.Equal(t, expectedVariants, result.Variants)
}
