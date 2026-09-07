package content_rules

import (
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

type ContentRuleService struct {
	distanceCatalog       *DistanceCatalog
	variantMappingCatalog *VariantMappingCatalog
}

func NewContentRuleService() IContentRuleService {
	return &ContentRuleService{
		distanceCatalog:       NewDistanceCatalog(),
		variantMappingCatalog: NewVariantMappingCatalog(),
	}
}

func (this *ContentRuleService) GetRules() []IContentRule {
	defaultMapping := this.variantMappingCatalog.GetDefaultMapping()
	variantRule, _ := NewRuleVariant(&defaultMapping, nil)
	return []IContentRule{
		NewRuleDistanceToRoad(nil),
		NewRuleDistanceToTown(nil),
		NewRuleGuarded(false),
		variantRule,
		NewRuleSoloEncounter(false),
	}
}

func (this *ContentRuleService) ApplyRulesToItem(item *template_model.MandatoryContentItem, rules []IContentRule) {
	for _, rule := range rules {
		if rule != nil {
			rule.Apply(item)
		}
	}
}

func (this *ContentRuleService) CreateRuleFromSavedRule(
	saved editor_state_model.ContentRuleRow,
	content models.SidMapping) IContentRule {
	switch {
	case strings.EqualFold(saved.Name, RuleDistanceToRoadName):
		if variation, ok := this.distanceCatalog.GetByName(saved.DistanceName); ok {
			return NewRuleDistanceToRoad(&variation)
		}
	case strings.EqualFold(saved.Name, RuleDistanceToTownName):
		if variation, ok := this.distanceCatalog.GetByName(saved.DistanceName); ok {
			return NewRuleDistanceToTown(&variation)
		}
	case strings.EqualFold(saved.Name, RuleGuardedName):
		if saved.IsGuarded != nil {
			return NewRuleGuarded(*saved.IsGuarded)
		}
	case strings.EqualFold(saved.Name, RuleSoloEncounterName):
		if saved.IsSoloEncounter != nil {
			return NewRuleSoloEncounter(*saved.IsSoloEncounter)
		}
	case strings.EqualFold(saved.Name, RuleVariantName):
		if saved.VariantID == nil {
			return nil
		}

		mapping, ok := this.variantMappingCatalog.GetVariantForContentByID(content, *saved.VariantID)
		if !ok {
			return nil
		}

		rule, err := NewRuleVariant(&mapping, saved.VariantID)
		if err == nil {
			return rule
		}
	}

	return nil
}

func (this *ContentRuleService) RestoreRulesFromRow(
	row editor_state_model.ZoneContentRow,
	content models.SidMapping) []IContentRule {
	var result []IContentRule
	for _, savedRule := range row.Rules {
		if rule := this.CreateRuleFromSavedRule(savedRule, content); rule != nil {
			result = append(result, rule)
		}
	}
	return result
}

func (this *ContentRuleService) GetDistanceDisplayNames() []string {
	return this.distanceCatalog.GetDisplayNames()
}

func (this *ContentRuleService) GetVariantsForContent(content models.SidMapping) []models.VariantMapping {
	return this.variantMappingCatalog.GetVariantsForContent(content)
}

func (this *ContentRuleService) GetVariantForContentByID(
	content models.SidMapping,
	variantID int) (models.VariantMapping, bool) {
	return this.variantMappingCatalog.GetVariantForContentByID(content, variantID)
}
