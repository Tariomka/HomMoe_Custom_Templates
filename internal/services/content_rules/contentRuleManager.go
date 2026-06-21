package content_rules

import (
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

// GetRules returns one prototype instance of every known content rule, in the
// same order as the C# ContentRuleManager static fields.
func GetRules() []ContentRule {
	variantRule, _ := NewRuleVariant(nil, nil)
	return []ContentRule{
		NewRuleDistanceToRoad(nil),
		NewRuleDistanceToTown(nil),
		NewRuleGuarded(false),
		variantRule,
		NewRuleSoloEncounter(false),
	}
}

// ApplyRulesToItem applies every rule, in order, to the final content item.
// Distance rules append placement rules; guarded/variant/solo rules set fields.
func ApplyRulesToItem(item *entities.MandatoryContentItem, rules []ContentRule) {
	for _, rule := range rules {
		if rule != nil {
			rule.Apply(item)
		}
	}
}

// CreateRuleFromSavedRule reconstructs a concrete rule from its persisted form.
// It returns nil when the saved data does not map to a known, valid rule.
func CreateRuleFromSavedRule(saved models.ContentRuleRowSave, content models.SidMapping) ContentRule {
	switch {
	case strings.EqualFold(saved.Name, RuleDistanceToRoadName):
		if variation, ok := GetDistanceVariationByName(saved.DistanceName); ok {
			return NewRuleDistanceToRoad(&variation)
		}
		return nil
	case strings.EqualFold(saved.Name, RuleDistanceToTownName):
		if variation, ok := GetDistanceVariationByName(saved.DistanceName); ok {
			return NewRuleDistanceToTown(&variation)
		}
		return nil
	case strings.EqualFold(saved.Name, RuleGuardedName):
		if saved.IsGuarded == nil {
			return nil
		}
		return NewRuleGuarded(*saved.IsGuarded)
	case strings.EqualFold(saved.Name, RuleSoloEncounterName):
		if saved.IsSoloEncounter == nil {
			return nil
		}
		return NewRuleSoloEncounter(*saved.IsSoloEncounter)
	case strings.EqualFold(saved.Name, RuleVariantName):
		if saved.VariantId == nil {
			return nil
		}
		mapping, ok := GetVariantForContentById(content, *saved.VariantId)
		if !ok {
			return nil
		}
		rule, err := NewRuleVariant(&mapping, saved.VariantId)
		if err != nil {
			return nil
		}
		return rule
	}
	return nil
}

// RestoreRulesFromRow reconstructs the rule list for a content row. New-format
// rows use their serialized Rules; legacy rows are migrated from the deprecated
// flat fields, matching the C# RestoreContentRulesFromRow fallback exactly.
func RestoreRulesFromRow(row models.ZoneContentRowSave, content models.SidMapping) []ContentRule {
	var result []ContentRule

	if len(row.Rules) > 0 {
		for _, saved := range row.Rules {
			if rule := CreateRuleFromSavedRule(saved, content); rule != nil {
				result = append(result, rule)
			}
		}
	}

	return result
}
