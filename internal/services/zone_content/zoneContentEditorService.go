package zone_content

import (
	"slices"
	"sort"
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/content_rule_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
)

// contentRuleMarkerSeparator joins the marker badges shown on a content row.
const contentRuleMarkerSeparator = " · "

type ZoneContentEditorService struct{}

func NewZoneContentEditorService() IZoneContentEditorService {
	return &ZoneContentEditorService{}
}

func (this *ZoneContentEditorService) ComposeContentRule(
	composition content_rule_model.ContentRuleComposition) (editor_state_model.ContentRuleRow, bool) {
	switch composition.Key {
	case content_rule_model.ContentRuleKeyDistanceToRoad, content_rule_model.ContentRuleKeyDistanceToTown:
		if composition.DistanceIndex < 0 || composition.DistanceIndex >= len(composition.DistanceNames) {
			return editor_state_model.ContentRuleRow{}, false
		}

		return editor_state_model.ContentRuleRow{
			Name:         composition.Name,
			DistanceName: composition.DistanceNames[composition.DistanceIndex],
		}, true
	case content_rule_model.ContentRuleKeyGuarded:
		guarded := composition.IsGuarded
		return editor_state_model.ContentRuleRow{Name: composition.Name, IsGuarded: &guarded}, true
	case content_rule_model.ContentRuleKeySoloEncounter:
		solo := composition.IsSoloEncounter
		return editor_state_model.ContentRuleRow{Name: composition.Name, IsSoloEncounter: &solo}, true
	case content_rule_model.ContentRuleKeyVariant:
		if composition.VariantIndex < 0 || composition.VariantIndex >= len(composition.VariantIDs) {
			return editor_state_model.ContentRuleRow{}, false
		}

		variantID := composition.VariantIDs[composition.VariantIndex]
		return editor_state_model.ContentRuleRow{Name: composition.Name, VariantID: &variantID}, true
	}

	return editor_state_model.ContentRuleRow{}, false
}

func (this *ZoneContentEditorService) UpsertContentRule(
	rules []editor_state_model.ContentRuleRow,
	rule editor_state_model.ContentRuleRow) []editor_state_model.ContentRuleRow {
	for index := range rules {
		if strings.EqualFold(rules[index].Name, rule.Name) {
			rules[index] = rule
			return rules
		}
	}

	return append(rules, rule)
}

func (this *ZoneContentEditorService) GetDefaultContentRules(
	options []content_rule_model.ContentRuleOption) []editor_state_model.ContentRuleRow {
	for _, option := range options {
		if option.Key == content_rule_model.ContentRuleKeyGuarded {
			guarded := true

			return []editor_state_model.ContentRuleRow{{Name: option.Name, IsGuarded: &guarded}}
		}
	}

	return nil
}

func (this *ZoneContentEditorService) GetContentRuleMarkers(
	descriptions []content_rule_model.ContentRuleDescription) string {
	markers := make([]string, 0, len(descriptions))
	for _, description := range descriptions {
		if description.Valid && description.Marker != "" {
			markers = append(markers, description.Marker)
		}
	}
	return strings.Join(markers, contentRuleMarkerSeparator)
}

func (this *ZoneContentEditorService) GetContentRowDisplayName(
	name string,
	descriptions []content_rule_model.ContentRuleDescription) string {
	for _, description := range descriptions {
		if description.Key == content_rule_model.ContentRuleKeyVariant && description.Valid {
			return name + " (" + description.VariantLabel + ")"
		}
	}

	return name
}

func (this *ZoneContentEditorService) SortContentItemsByName(items []models.SidMapping) []models.SidMapping {
	sorted := slices.Clone(items)
	sort.SliceStable(sorted, func(first int, second int) bool {
		return strings.ToLower(sorted[first].Name) < strings.ToLower(sorted[second].Name)
	})
	return sorted
}

func (this *ZoneContentEditorService) ClampContentCount(count int, maxCount int) int {
	if count < 1 {
		return 1
	}

	if count > maxCount {
		return maxCount
	}

	return count
}
