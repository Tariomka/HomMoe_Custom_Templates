package services

import (
	"fmt"
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
)

var (
	distNextTo  = distancePreset{0.05, 0.1}
	distMedium  = distancePreset{0.25, 0.5}
	distFar     = distancePreset{0.5, 0.75}
	distVeryFar = distancePreset{0.75, 0.9}
)

// distanceForLabel resolves the road-distance label persisted in a
// ZoneContentRowSave to the matching distancePreset. An empty or unknown
// label means "Any" (no constraint added).
func distanceForLabel(label string) (distancePreset, bool) {
	switch strings.TrimSpace(label) {
	case "Next To":
		return distNextTo, true
	case "Near":
		return distNear, true
	case "Medium":
		return distMedium, true
	case "Far":
		return distFar, true
	case "Very Far":
		return distVeryFar, true
	}
	return distancePreset{}, false
}

// ── Rule presets ─────────────────────────────────────────────────────

func ruleRoadDistance(d distancePreset, weight int) template.PlacementRule {
	return template.PlacementRule{Type: "Road", Args: []any{}, TargetMin: d.Min, TargetMax: d.Max, Weight: weight}
}

func ruleCrossroadsDistance(d distancePreset, weight int) template.PlacementRule {
	return template.PlacementRule{Type: "Crossroads", Args: []any{}, TargetMin: d.Min, TargetMax: d.Max, Weight: weight}
}

func ruleNearCastle(weight int) template.PlacementRule {
	return template.PlacementRule{Type: "MainObject", Args: []any{"0"}, TargetMin: 0.1, TargetMax: 0.3, Weight: weight}
}

// ── Content builder (fluent) ─────────────────────────────────────────

type contentBuilder struct {
	item template.MandatoryContentItem
}

func newContentBuilder(sid string) *contentBuilder {
	return &contentBuilder{item: template.MandatoryContentItem{SID: sid}}
}

func (this *contentBuilder) withName(name string) *contentBuilder { this.item.Name = name; return this }
func (this *contentBuilder) guarded() *contentBuilder             { this.item.IsGuarded = true; return this }
func (this *contentBuilder) mine() *contentBuilder                { this.item.IsMine = true; return this }
func (this *contentBuilder) soloEncounter() *contentBuilder {
	this.item.SoloEncounter = true
	return this
}
func (this *contentBuilder) roadDistance(d distancePreset) *contentBuilder {
	return this.addRule(ruleRoadDistance(d, 1))
}
func (this *contentBuilder) addRule(r template.PlacementRule) *contentBuilder {
	this.item.Rules = append(this.item.Rules, r)
	return this
}
func (this *contentBuilder) addRules(rs []template.PlacementRule) *contentBuilder {
	this.item.Rules = append(this.item.Rules, rs...)
	return this
}
func (this *contentBuilder) build() template.MandatoryContentItem { return this.item }

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
	return newContentBuilder(constants.ContentIds.RemoteFoothold.Sid).
		withName("name_remote_foothold_1").
		soloEncounter().
		addRules(footholdRules(castleCount)).
		build()
}

// RowsToMandatoryContent expands a slice of save-rows to a flat list of
// MandatoryContentItem entries suitable for a MandatoryContent.Content.
func RowsToMandatoryContent(rows []models.ZoneContentRowSave) []template.MandatoryContentItem {
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
			out = append(out, rowToMandatoryItem(row))
		}
	}
	return out
}

func rowToMandatoryItem(row models.ZoneContentRowSave) template.MandatoryContentItem {
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
		item.Rules = append(item.Rules, ruleNearCastle(1))
	}
	if d, ok := distanceForLabel(row.RoadDistance); ok {
		item.Rules = append(item.Rules, ruleRoadDistance(d, 1))
	}
	return item
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
	settings *models.GeneratorSettings,
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
func BuildPlayerZoneMandatoryContent(settings *models.GeneratorSettings) []template.MandatoryContentItem {
	return contentWithFootholdAndRows(settings, settings.ZoneCfg.PlayerZoneCastles, settings.PlayerZoneMandatoryContent)
}

// BuildLowNeutralMandatoryContent returns mandatory content for a low-quality neutral zone.
func BuildLowNeutralMandatoryContent(settings *models.GeneratorSettings, castleCount int) []template.MandatoryContentItem {
	rows := settings.LowNeutralMandatoryContent
	if castleCount == 0 {
		rows = StripNearCastleRules(append([]template.MandatoryContentItem(nil), rows...))
	}
	return contentWithFootholdAndRows(settings, castleCount, rows)
}

// BuildMediumNeutralMandatoryContent returns mandatory content for a medium-quality neutral zone.
func BuildMediumNeutralMandatoryContent(settings *models.GeneratorSettings, castleCount int) []template.MandatoryContentItem {
	rows := settings.MediumNeutralMandatoryContent
	if castleCount == 0 {
		rows = StripNearCastleRules(append([]template.MandatoryContentItem(nil), rows...))
	}
	return contentWithFootholdAndRows(settings, castleCount, rows)
}

// BuildHighNeutralMandatoryContent returns mandatory content for a high-quality neutral zone.
func BuildHighNeutralMandatoryContent(settings *models.GeneratorSettings, castleCount int) []template.MandatoryContentItem {
	rows := settings.HighNeutralMandatoryContent
	if castleCount == 0 {
		rows = StripNearCastleRules(append([]template.MandatoryContentItem(nil), rows...))
	}
	return contentWithFootholdAndRows(settings, castleCount, rows)
}

// BuildHubZoneMandatoryContent returns mandatory content for the central
// hub zone of a Hub-and-Spoke layout
func BuildHubZoneMandatoryContent(settings *models.GeneratorSettings, castleCount int) []template.MandatoryContentItem {
	rows := settings.HubZoneMandatoryContent
	if castleCount == 0 {
		rows = StripNearCastleRules(append([]template.MandatoryContentItem(nil), rows...))
	}
	return contentWithFootholdAndRows(settings, castleCount, rows)
}

func BuildAllContentCountLimits(settings *models.GeneratorSettings) []template.ContentCountLimit {
	sidLimits := []template.ContentLimit{
		{SID: "black_tower", MaxCount: 0},
		{SID: constants.ContentIds.Fountain.Sid, MaxCount: 2},
		{SID: constants.ContentIds.Fountain2.Sid, MaxCount: 2},
		{SID: constants.ContentIds.ManaWell.Sid, MaxCount: 2},
		{SID: constants.ContentIds.BeerFountain.Sid, MaxCount: 2},
		{SID: constants.ContentIds.Market.Sid, MaxCount: 1},
		{SID: constants.ContentIds.Forge.Sid, MaxCount: 2},
		{SID: constants.ContentIds.Stables.Sid, MaxCount: 1},
		{SID: constants.ContentIds.Watchtower.Sid, MaxCount: 2},
		{SID: constants.ContentIds.WindRose.Sid, MaxCount: 1},
		{SID: constants.ContentIds.QuixsPath.Sid, MaxCount: 2},
		{SID: constants.ContentIds.CrystalTrail.Sid, MaxCount: 3},
		{SID: constants.ContentIds.MysteriousStone.Sid, MaxCount: 2},
		{SID: constants.ContentIds.University.Sid, MaxCount: 2},
		{SID: constants.ContentIds.WiseOwl.Sid, MaxCount: 4},
		{SID: constants.ContentIds.CelestialSphere.Sid, MaxCount: 2},
		{SID: constants.ContentIds.PileOfBooks.Sid, MaxCount: 2},
		{SID: constants.ContentIds.InsarasEye.Sid, MaxCount: 2},
		{SID: constants.ContentIds.TearOfTruth.Sid, MaxCount: 3},
		{SID: constants.ContentIds.TreeOfAbundance.Sid, MaxCount: 2},
		{SID: constants.ContentIds.HuntsmansCamp.Sid, MaxCount: 2},
		{SID: constants.ContentIds.ShadyDen.Sid, MaxCount: 2},
		{SID: constants.ContentIds.RandomHire1.Sid, MaxCount: 6},
		{SID: constants.ContentIds.RandomHire2.Sid, MaxCount: 6},
		{SID: constants.ContentIds.RandomHire3.Sid, MaxCount: 6},
		{SID: constants.ContentIds.RandomHire4.Sid, MaxCount: 6},
		{SID: constants.ContentIds.RandomHire5.Sid, MaxCount: 6},
		{SID: constants.ContentIds.RandomHire6.Sid, MaxCount: 6},
		{SID: constants.ContentIds.RandomHire7.Sid, MaxCount: 6},
		{SID: constants.ContentIds.Arena.Sid, MaxCount: 2},
		{SID: constants.ContentIds.SacrificialShrine.Sid, MaxCount: 2},
		{SID: constants.ContentIds.Chimerologist.Sid, MaxCount: 2},
		{SID: constants.ContentIds.Circus.Sid, MaxCount: 2},
		{SID: constants.ContentIds.InfernalCirque.Sid, MaxCount: 2},
		{SID: constants.ContentIds.FlatteringMirror.Sid, MaxCount: 2},
		{SID: constants.ContentIds.FickleShrine.Sid, MaxCount: 1},
		{SID: constants.ContentIds.PointOfBalance.Sid, MaxCount: 3},
		{SID: constants.ContentIds.PandoraBox.Sid, MaxCount: 4},
		{SID: constants.ContentIds.RitualPyre.Sid, MaxCount: 3},
		{SID: constants.ContentIds.BorealCall.Sid, MaxCount: 3},
		{SID: constants.ContentIds.JoustingRange.Sid, MaxCount: 1},
		{SID: constants.ContentIds.UnforgottenGrave.Sid, MaxCount: 1},
		{SID: constants.ContentIds.PetrifiedMemorial.Sid, MaxCount: 1},
		{SID: constants.ContentIds.TheGorge.Sid, MaxCount: 1},
	}

	// Lift limits when any mandatory-content list (player or neutral or
	// hub) requests more of a given SID than the default cap.
	sidCounts := map[string]int{}
	tally := func(items []template.MandatoryContentItem) {
		for _, item := range items {
			if item.SID != "" {
				sidCounts[strings.ToLower(item.SID)]++
			}
		}
	}
	tally(settings.PlayerZoneMandatoryContent)
	tally(settings.LowNeutralMandatoryContent)
	tally(settings.MediumNeutralMandatoryContent)
	tally(settings.HighNeutralMandatoryContent)
	tally(settings.HubZoneMandatoryContent)
	for i := range sidLimits {
		if count, ok := sidCounts[strings.ToLower(sidLimits[i].SID)]; ok {
			if count > sidLimits[i].MaxCount {
				sidLimits[i].MaxCount = count
			}
		}
	}

	var limits []template.ContentCountLimit
	limits = append(limits, template.ContentCountLimit{Name: "content_limits_side", Limits: sidLimits})
	limits = append(limits, template.ContentCountLimit{Name: "content_limits_side_0_0", Limits: sidLimits})
	for a := 1; a <= 5; a++ {
		for b := a + 1; b <= 6; b++ {
			limits = append(limits, template.ContentCountLimit{
				Name:   fmt.Sprintf("content_limits_side_%d_%d", a, b),
				Limits: sidLimits,
			})
		}
	}
	return limits
}
