package providers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/builders/mandatory_content"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/builders/placement_rule"
)

type MandatoryContentProvider struct{}

func NewMandatoryContentProvider() *MandatoryContentProvider {
	return &MandatoryContentProvider{}
}

func (this *MandatoryContentProvider) CreateContents(
	configuration config.GeneratorConfig,
	playerLabels []string,
	neutralZones models.NeutralZonePlans) []template.MandatoryContent {
	var groups []template.MandatoryContent
	for _, letter := range playerLabels {
		groups = append(groups, template.MandatoryContent{
			Name: "mandatory_content_side_" + letter,
			Content: this.createContentItemsWithFoothold(
				configuration.PlayerZoneMandatoryContent,
				configuration.SpawnRemoteFootholds,
				configuration.ZoneConfiguration.PlayerZoneCastles),
		})
	}
	for _, neutralZone := range neutralZones {
		var content []template.MandatoryContentItem
		switch neutralZone.Quality {
		case models.QualityLow:
			copy(content, configuration.LowNeutralMandatoryContent)
		case models.QualityMedium:
			copy(content, configuration.MediumNeutralMandatoryContent)
		case models.QualityHigh:
			copy(content, configuration.HighNeutralMandatoryContent)
		}
		if neutralZone.CastleCount == 0 {
			content = stripNearCastleRules(content)
		}
		groups = append(groups, template.MandatoryContent{
			Name:    "mandatory_content_neutral_" + neutralZone.Label,
			Content: this.createContentItemsWithFoothold(content, configuration.SpawnRemoteFootholds, neutralZone.CastleCount),
		})
	}
	return groups
}

func (this *MandatoryContentProvider) CreateContentItemsFrom(
	rows []models.ZoneContentRowSave) []template.MandatoryContentItem {
	if len(rows) == 0 {
		return nil
	}
	var out []template.MandatoryContentItem
	for _, raw := range rows {
		row := raw.Normalised()
		if row.Sid == "" {
			continue
		}
		for i := 0; i < row.Count; i++ {
			out = append(out, this.createContentItemFrom(row))
		}
	}
	return out
}

func (this *MandatoryContentProvider) createContentItemFrom(
	row models.ZoneContentRowSave) template.MandatoryContentItem {
	item := template.MandatoryContentItem{
		IsGuarded: row.IsGuarded,
		IsMine:    row.IsMine,
	}
	if row.IsGroup {
		item.IncludeLists = []string{row.Sid}
	} else {
		item.SID = row.Sid
	}
	if row.NearCastle {
		item.Rules = append(item.Rules, placement_rule.NewPlacementRuleBuilder().BuildNearCastleRule(1))
	}
	if distance, ok := placement_rule.TryGetDistanceFrom(row.RoadDistance); ok {
		item.Rules = append(item.Rules, placement_rule.NewPlacementRuleBuilder().BuildRoadRule(distance, 1))
	}
	return item
}

func (this *MandatoryContentProvider) createContentItemsWithFoothold(
	rows []template.MandatoryContentItem,
	addFoothold bool,
	castleCount int) []template.MandatoryContentItem {
	var content []template.MandatoryContentItem
	if addFoothold {
		content = append(content, this.createFootholdContentItem(castleCount))
	}
	content = append(content, rows...)
	return content
}

func (this *MandatoryContentProvider) createFootholdContentItem(
	castleCount int) template.MandatoryContentItem {
	return mandatory_content.NewContentBuilder(constants.ContentIds.RemoteFoothold.Sid).
		WithName("name_remote_foothold_1").
		WithSoloEncounter().
		WithRulesCallback(func() []template.PlacementRule {
			rules := []template.PlacementRule{
				placement_rule.NewPlacementRuleBuilder().
					BuildCrossroadsRule(placement_rule.Distance{Min: 0.2, Max: 0.3}, 0),
			}
			if castleCount > 0 {
				rules = append(rules,
					placement_rule.NewPlacementRuleBuilder().
						WithTypeMainObject().
						WithArgs("0").
						WithDistance(placement_rule.Distance{Min: 0.2, Max: 0.4}).
						WithWeight(0).
						Build())
			}
			if castleCount > 1 {
				rules = append(rules,
					placement_rule.NewPlacementRuleBuilder().
						WithTypeMainObject().
						WithArgs("1").
						WithDistance(placement_rule.Distance{Min: 0.5, Max: 0.5}).
						WithWeight(2).
						Build())
			}
			return rules
		}).
		Build()
}

// stripNearCastleRules removes placement rules that anchor an item near
// the zone's main castle. Used when a zone has no castle so the rule
// would never be satisfiable
func stripNearCastleRules(items []template.MandatoryContentItem) []template.MandatoryContentItem {
	for i := range items {
		if len(items[i].Rules) == 0 {
			continue
		}
		kept := items[i].Rules[:0]
		for _, rule := range items[i].Rules {
			if rule.Type == "MainObject" && len(rule.Args) > 0 {
				if arg, ok := rule.Args[0].(string); ok && arg == "0" {
					continue
				}
			}
			kept = append(kept, rule)
		}
		items[i].Rules = kept
	}
	return items
}
