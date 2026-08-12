package handlers

import (
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
)

type contentRuleHandler struct {
	contentRuleService content_rules.IContentRuleService
}

func NewContentRuleHandler(
	contentRuleService content_rules.IContentRuleService) handler_interfaces.IContentRuleHandler {
	return &contentRuleHandler{contentRuleService: contentRuleService}
}

func (this *contentRuleHandler) GetContentRuleEditorOptions(
	content models.SidMapping) dtos.ContentRuleEditorOptionsDto {
	rules := []dtos.ContentRuleOptionDto{
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
	}

	variants := this.contentRuleVariantOptions(content)
	if len(variants) > 0 {
		rules = append(rules, dtos.ContentRuleOptionDto{
			Key:         dtos.ContentRuleKeyVariant,
			Name:        content_rules.RuleVariantName,
			Description: content_rules.RuleVariantDescription,
			Marker:      content_rules.RuleVariantMarker,
			EditorKind:  dtos.ContentRuleEditorKindVariant,
			EditorLabel: "Variant",
		})
	}

	return dtos.ContentRuleEditorOptionsDto{
		Rules:     rules,
		Distances: this.contentRuleService.GetDistanceDisplayNames(),
		Variants:  variants,
	}
}

func (this *contentRuleHandler) DescribeContentRule(
	content models.SidMapping,
	savedRule models.ContentRuleRowSave) dtos.ContentRuleDescriptionDto {
	description := dtos.ContentRuleDescriptionDto{
		Key:         contentRuleKeyFromName(savedRule.Name),
		DisplayText: savedRule.Name,
		SavedRule:   savedRule,
	}
	rule := this.contentRuleService.CreateRuleFromSavedRule(savedRule, content)
	if rule == nil {
		return description
	}

	description.DisplayText = rule.DisplayText()
	description.Marker = rule.Marker()
	description.Valid = true
	if savedRule.VariantID != nil {
		variant, ok := this.contentRuleService.GetVariantForContentByID(content, *savedRule.VariantID)
		if ok {
			description.VariantLabel, _ = variant.GetVariantByID(*savedRule.VariantID)
		}
	}
	return description
}

func (this *contentRuleHandler) contentRuleVariantOptions(
	content models.SidMapping) []dtos.ContentRuleVariantOptionDto {
	variants := this.contentRuleService.GetVariantsForContent(content)
	options := make([]dtos.ContentRuleVariantOptionDto, 0, len(variants))
	for _, variant := range variants {
		for _, tuple := range variant.Variants {
			options = append(options, dtos.ContentRuleVariantOptionDto{ID: tuple.Key, Label: tuple.Value})
		}
	}
	return options
}

func contentRuleKeyFromName(name string) dtos.ContentRuleKey {
	switch {
	case strings.EqualFold(name, content_rules.RuleDistanceToRoadName):
		return dtos.ContentRuleKeyDistanceToRoad
	case strings.EqualFold(name, content_rules.RuleDistanceToTownName):
		return dtos.ContentRuleKeyDistanceToTown
	case strings.EqualFold(name, content_rules.RuleGuardedName):
		return dtos.ContentRuleKeyGuarded
	case strings.EqualFold(name, content_rules.RuleSoloEncounterName):
		return dtos.ContentRuleKeySoloEncounter
	case strings.EqualFold(name, content_rules.RuleVariantName):
		return dtos.ContentRuleKeyVariant
	default:
		return ""
	}
}
