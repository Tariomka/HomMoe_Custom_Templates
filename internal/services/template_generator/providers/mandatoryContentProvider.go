package providers

import (
	"fmt"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
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
	neutralZones models.NeutralZonePlans) []entities.MandatoryContent {
	var groups []entities.MandatoryContent
	footholdCount := 0
	if configuration.SpawnRemoteFootholds {
		footholdCount = configuration.RemoteFootholdCount
	}
	for _, letter := range playerLabels {
		groups = append(groups, entities.MandatoryContent{
			Name: "mandatory_content_side_" + letter,
			Content: this.createContentItemsWithFoothold(
				configuration.PlayerZoneMandatoryContent,
				footholdCount,
				configuration.ZoneConfiguration.PlayerZoneCastles),
		})
	}
	for _, neutralZone := range neutralZones {
		var content []entities.MandatoryContentItem
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
		groups = append(groups, entities.MandatoryContent{
			Name:    "mandatory_content_neutral_" + neutralZone.Label,
			Content: this.createContentItemsWithFoothold(content, footholdCount, neutralZone.CastleCount),
		})
	}
	return groups
}

func (this *MandatoryContentProvider) CreateContentItemsFrom(
	rows []models.ZoneContentRowSave) []entities.MandatoryContentItem {
	if len(rows) == 0 {
		return nil
	}
	var out []entities.MandatoryContentItem
	for _, raw := range rows {
		row := raw.Normalized()
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
	row models.ZoneContentRowSave) entities.MandatoryContentItem {
	item := entities.MandatoryContentItem{
		IsMine: row.IsMine,
	}
	if row.IsGroup {
		item.IncludeLists = []string{row.Sid}
	} else {
		item.SID = row.Sid
	}
	rules := content_rules.RestoreRulesFromRow(row, models.SidMapping{Sid: row.Sid})
	content_rules.ApplyRulesToItem(&item, rules)
	return item
}

func (this *MandatoryContentProvider) createContentItemsWithFoothold(
	rows []entities.MandatoryContentItem,
	footholdCount int,
	castleCount int) []entities.MandatoryContentItem {
	var content []entities.MandatoryContentItem
	for i := 1; i <= footholdCount; i++ {
		content = append(content, this.createFootholdContentItem(i, castleCount))
	}
	content = append(content, rows...)
	return content
}

func (this *MandatoryContentProvider) createFootholdContentItem(
	index int,
	castleCount int) entities.MandatoryContentItem {
	return mandatory_content.NewContentBuilder(nonContentObjects.RemoteFoothold).
		WithName(fmt.Sprintf("name_remote_foothold_%d", index)).
		WithSoloEncounter().
		WithRulesCallback(func() []entities.PlacementRule {
			rules := []entities.PlacementRule{
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
func stripNearCastleRules(items []entities.MandatoryContentItem) []entities.MandatoryContentItem {
	for i := range items {
		if len(items[i].Rules) == 0 {
			continue
		}
		kept := items[i].Rules[:0]
		for _, rule := range items[i].Rules {
			if rule.Type == ruleTypes.MainObject && len(rule.Args) > 0 {
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
