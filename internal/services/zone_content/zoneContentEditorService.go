package zone_content

import (
	"slices"
	"sort"
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

// contentRuleMarkerSeparator joins the marker badges shown on a content row.
const contentRuleMarkerSeparator = " · "

type ZoneContentEditorService struct{}

func NewZoneContentEditorService() IZoneContentEditorService {
	return &ZoneContentEditorService{}
}

func (this *ZoneContentEditorService) ComposeContentRule(
	request dtos.ContentRuleCompositionRequestDto) dtos.ContentRuleCompositionResultDto {
	switch request.Option.Key {
	case dtos.ContentRuleKeyDistanceToRoad, dtos.ContentRuleKeyDistanceToTown:
		if request.DistanceIndex < 0 || request.DistanceIndex >= len(request.DistanceNames) {
			return dtos.ContentRuleCompositionResultDto{}
		}

		return validRule(models.ContentRuleRow{
			Name:         request.Option.Name,
			DistanceName: request.DistanceNames[request.DistanceIndex],
		})
	case dtos.ContentRuleKeyGuarded:
		guarded := request.IsGuarded
		return validRule(models.ContentRuleRow{Name: request.Option.Name, IsGuarded: &guarded})
	case dtos.ContentRuleKeySoloEncounter:
		solo := request.IsSoloEncounter
		return validRule(models.ContentRuleRow{Name: request.Option.Name, IsSoloEncounter: &solo})
	case dtos.ContentRuleKeyVariant:
		if request.VariantIndex < 0 || request.VariantIndex >= len(request.VariantIDs) {
			return dtos.ContentRuleCompositionResultDto{}
		}

		variantID := request.VariantIDs[request.VariantIndex]
		return validRule(models.ContentRuleRow{Name: request.Option.Name, VariantID: &variantID})
	}

	return dtos.ContentRuleCompositionResultDto{}
}

func (this *ZoneContentEditorService) UpsertContentRule(
	rules []models.ContentRuleRow,
	rule models.ContentRuleRow) []models.ContentRuleRow {
	for index := range rules {
		if strings.EqualFold(rules[index].Name, rule.Name) {
			rules[index] = rule
			return rules
		}
	}

	return append(rules, rule)
}

func (this *ZoneContentEditorService) GetDefaultContentRules(
	options dtos.ContentRuleEditorOptionsDto) []models.ContentRuleRow {
	for _, option := range options.Rules {
		if option.Key == dtos.ContentRuleKeyGuarded {
			guarded := true

			return []models.ContentRuleRow{{Name: option.Name, IsGuarded: &guarded}}
		}
	}

	return nil
}

func (this *ZoneContentEditorService) GetContentRuleMarkers(descriptions []dtos.ContentRuleDescriptionDto) string {
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
	descriptions []dtos.ContentRuleDescriptionDto) string {
	for _, description := range descriptions {
		if description.Key == dtos.ContentRuleKeyVariant && description.Valid {
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

func validRule(rule models.ContentRuleRow) dtos.ContentRuleCompositionResultDto {
	return dtos.ContentRuleCompositionResultDto{Rule: rule, Valid: true}
}
