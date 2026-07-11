package providers

import (
	"fmt"
	"slices"
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/mandatory_content"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/placement_rule"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
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
			content = cloneContentItems(configuration.LowNeutralMandatoryContent)
		case models.QualityMedium:
			content = cloneContentItems(configuration.MediumNeutralMandatoryContent)
		case models.QualityHigh:
			content = cloneContentItems(configuration.HighNeutralMandatoryContent)
		}
		if neutralZone.CastleCount == 0 {
			content = stripNearCastleRules(content)
		}
		groups = append(groups, entities.MandatoryContent{
			Name:    "mandatory_content_neutral_" + neutralZone.Label,
			Content: this.createContentItemsWithFoothold(content, footholdCount, neutralZone.CastleCount),
		})
	}
	if hub, ok := this.hubContentGroup(configuration); ok {
		groups = append(groups, hub)
	}
	return groups
}

// CreateContentsForZones rebuilds the mandatory-content groups from the final
// zones (after any manual edits) instead of from the original generation plan.
// A zone whose quality or castle count was changed in the manual zone editor no
// longer matches its plan, so keying content off the plan (as CreateContents
// does) would give a re-tiered zone the wrong content - e.g. a zone manually
// promoted to High would keep its original Medium mandatory content. Detecting
// the quality and castle count from the zone itself keeps the two in sync.
func (this *MandatoryContentProvider) CreateContentsForZones(
	configuration config.GeneratorConfig,
	zones []entities.Zone) []entities.MandatoryContent {
	footholdCount := 0
	if configuration.SpawnRemoteFootholds {
		footholdCount = configuration.RemoteFootholdCount
	}

	var groups []entities.MandatoryContent
	hubGroupAdded := false
	for _, zone := range zones {
		switch {
		case strings.HasPrefix(zone.Name, "Spawn-"):
			groups = append(groups, entities.MandatoryContent{
				Name: "mandatory_content_side_" + strings.TrimPrefix(zone.Name, "Spawn-"),
				Content: this.createContentItemsWithFoothold(
					cloneContentItems(configuration.PlayerZoneMandatoryContent),
					footholdCount,
					configuration.ZoneConfiguration.PlayerZoneCastles),
			})
		case strings.HasPrefix(zone.Name, "Neutral-"):
			castleCount := countCityMainObjects(zone)
			content := cloneContentItems(neutralRowsForQuality(configuration, detectZoneQuality(zone)))
			if castleCount == 0 {
				content = stripNearCastleRules(content)
			}
			groups = append(groups, entities.MandatoryContent{
				Name:    "mandatory_content_neutral_" + strings.TrimPrefix(zone.Name, "Neutral-"),
				Content: this.createContentItemsWithFoothold(content, footholdCount, castleCount),
			})
		case zone.Name == "Hub" || strings.HasPrefix(zone.Name, "Hub-"):
			// One shared hub group even if several hub zones exist (tournament).
			if hubGroupAdded || len(configuration.HubZoneMandatoryContent) == 0 {
				continue
			}
			content := cloneContentItems(configuration.HubZoneMandatoryContent)
			if countCityMainObjects(zone) == 0 {
				content = stripNearCastleRules(content)
			}
			groups = append(groups, entities.MandatoryContent{Name: "mandatory_content_hub", Content: content})
			hubGroupAdded = true
		}
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
		for range row.Count {
			out = append(out, this.createContentItemFrom(row))
		}
	}
	return out
}

// hubContentGroup builds the hub zone's mandatory-content group from the
// configured hub rows. It only exists for the Hub & Spoke topology and only
// when the user actually configured hub content, matching the parallel C#
// editor which references "mandatory_content_hub" only when hub rows are set.
// The hub zone has no remote-foothold roads, so no foothold item is added.
func (this *MandatoryContentProvider) hubContentGroup(
	configuration config.GeneratorConfig) (entities.MandatoryContent, bool) {
	if configuration.Topology != config.TopologyHubAndSpoke || len(configuration.HubZoneMandatoryContent) == 0 {
		return entities.MandatoryContent{}, false
	}

	content := cloneContentItems(configuration.HubZoneMandatoryContent)
	if configuration.ZoneConfiguration.HubZoneCastles == 0 {
		content = stripNearCastleRules(content)
	}
	return entities.MandatoryContent{Name: "mandatory_content_hub", Content: content}, true
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
// would never be satisfiable.
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

// cloneContentItems deep-copies the mandatory-content items (including each
// item's Rules slice) so callers like stripNearCastleRules can mutate the copy
// without corrupting the shared per-tier rows held on the configuration. The
// original CreateContents used copy() into a nil slice, which silently dropped
// every row; cloning preserves them while keeping the per-zone isolation that
// copy() was meant to provide.
func cloneContentItems(items []entities.MandatoryContentItem) []entities.MandatoryContentItem {
	if len(items) == 0 {
		return nil
	}
	out := make([]entities.MandatoryContentItem, len(items))
	for i, item := range items {
		item.Rules = slices.Clone(item.Rules)
		out[i] = item
	}
	return out
}

// neutralRowsForQuality returns the configured mandatory-content rows matching a
// neutral zone's quality tier.
func neutralRowsForQuality(
	configuration config.GeneratorConfig,
	quality models.NeutralZoneQuality) []entities.MandatoryContentItem {
	switch quality {
	case models.QualityLow:
		return configuration.LowNeutralMandatoryContent
	case models.QualityHigh:
		return configuration.HighNeutralMandatoryContent
	case models.QualityMedium:
		fallthrough
	default:
		return configuration.MediumNeutralMandatoryContent
	}
}

// detectZoneQuality infers a neutral zone's quality tier from its guarded
// content pool, mirroring connection_editor.QualityOfZone so manual re-tiering
// in the zone editor is reflected when rebuilding mandatory content.
func detectZoneQuality(zone entities.Zone) models.NeutralZoneQuality {
	pool := ""
	if len(zone.GuardedContentPool) > 0 {
		pool = zone.GuardedContentPool[0]
	}
	if strings.Contains(pool, "_t4_") || strings.Contains(pool, "_t5_") {
		return models.QualityHigh
	}
	if strings.Contains(pool, "_t1_") || strings.Contains(pool, "_t2_") {
		return models.QualityLow
	}
	return models.QualityMedium
}

// countCityMainObjects returns the number of City main objects (castles) in a
// zone, ignoring abandoned outposts and other object types.
func countCityMainObjects(zone entities.Zone) int {
	count := 0
	for _, mainObject := range zone.MainObjects {
		if strings.EqualFold(mainObject.Type, "City") {
			count++
		}
	}
	return count
}
