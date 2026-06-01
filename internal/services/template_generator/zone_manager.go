package template_generator

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/builders/mandatory_content"
)

var (
	distNextTo  = distancePreset{0.05, 0.1}
	distMedium  = distancePreset{0.25, 0.5}
	distFar     = distancePreset{0.5, 0.75}
	distVeryFar = distancePreset{0.75, 0.9}
)

// ── Foothold preset (still used by every zone type) ──────────────────

func footholdRules(castleCount int) []template.PlacementRule {
	rules := []template.PlacementRule{
		{Type: "Crossroads", Args: []any{}, TargetMin: 0.2, TargetMax: 0.3, Weight: 0},
	}
	if castleCount > 0 {
		rules = append(rules, template.PlacementRule{Type: "MainObject", Args: []any{"0"}, TargetMin: 0.2, TargetMax: 0.4, Weight: 0})
	}
	if castleCount > 1 {
		rules = append(rules, template.PlacementRule{Type: "MainObject", Args: []any{"1"}, TargetMin: 0.5, TargetMax: 0.5, Weight: 2})
	}
	return rules
}

func presetRemoteFoothold(castleCount int) template.MandatoryContentItem {
	return mandatory_content.NewContentBuilder(constants.ContentIds.RemoteFoothold.Sid).
		WithName("name_remote_foothold_1").
		WithSoloEncounter().
		WithRules(footholdRules(castleCount)...).
		Build()
}

// StripNearCastleRules removes placement rules that anchor an item near
// the zone's main castle. Used when a zone has no castle so the rule
// would never be satisfiable
func StripNearCastleRules(items []template.MandatoryContentItem) []template.MandatoryContentItem {
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

func contentWithFootholdAndRows(
	settings *config.GeneratorConfig,
	castleCount int,
	rows []template.MandatoryContentItem,
) []template.MandatoryContentItem {
	var content []template.MandatoryContentItem
	if settings.SpawnRemoteFootholds {
		content = append(content, presetRemoteFoothold(castleCount))
	}
	content = append(content, rows...)
	return content
}

// BuildPlayerZoneMandatoryContent returns the mandatory content list for a player spawn zone.
func BuildPlayerZoneMandatoryContent(settings *config.GeneratorConfig) []template.MandatoryContentItem {
	return contentWithFootholdAndRows(settings, settings.ZoneConfiguration.PlayerZoneCastles, settings.PlayerZoneMandatoryContent)
}

// BuildLowNeutralMandatoryContent returns mandatory content for a low-quality neutral zone.
func BuildLowNeutralMandatoryContent(settings *config.GeneratorConfig, castleCount int) []template.MandatoryContentItem {
	rows := settings.LowNeutralMandatoryContent
	if castleCount == 0 {
		rows = StripNearCastleRules(append([]template.MandatoryContentItem(nil), rows...))
	}
	return contentWithFootholdAndRows(settings, castleCount, rows)
}

// BuildMediumNeutralMandatoryContent returns mandatory content for a medium-quality neutral zone.
func BuildMediumNeutralMandatoryContent(settings *config.GeneratorConfig, castleCount int) []template.MandatoryContentItem {
	rows := settings.MediumNeutralMandatoryContent
	if castleCount == 0 {
		rows = StripNearCastleRules(append([]template.MandatoryContentItem(nil), rows...))
	}
	return contentWithFootholdAndRows(settings, castleCount, rows)
}

// BuildHighNeutralMandatoryContent returns mandatory content for a high-quality neutral zone.
func BuildHighNeutralMandatoryContent(settings *config.GeneratorConfig, castleCount int) []template.MandatoryContentItem {
	rows := settings.HighNeutralMandatoryContent
	if castleCount == 0 {
		rows = StripNearCastleRules(append([]template.MandatoryContentItem(nil), rows...))
	}
	return contentWithFootholdAndRows(settings, castleCount, rows)
}

// BuildHubZoneMandatoryContent returns mandatory content for the central
// hub zone of a Hub-and-Spoke layout
func BuildHubZoneMandatoryContent(settings *config.GeneratorConfig, castleCount int) []template.MandatoryContentItem {
	rows := settings.HubZoneMandatoryContent
	if castleCount == 0 {
		rows = StripNearCastleRules(append([]template.MandatoryContentItem(nil), rows...))
	}
	return contentWithFootholdAndRows(settings, castleCount, rows)
}
