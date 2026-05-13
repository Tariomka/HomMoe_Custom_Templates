package services

import (
	"fmt"
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
)

// ── Distance presets ─────────────────────────────────────────────────

var (
	distNextTo = distancePreset{0.05, 0.1}
	// distMedium = distancePreset{0.25, 0.5}   // unused currently
	// distFar    = distancePreset{0.5, 0.75}
	// distVeryFar = distancePreset{0.75, 0.9}
)

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

func (b *contentBuilder) withName(name string) *contentBuilder { b.item.Name = name; return b }
func (b *contentBuilder) guarded() *contentBuilder             { b.item.IsGuarded = true; return b }
func (b *contentBuilder) mine() *contentBuilder                { b.item.IsMine = true; return b }
func (b *contentBuilder) soloEncounter() *contentBuilder       { b.item.SoloEncounter = true; return b }
func (b *contentBuilder) roadDistance(d distancePreset) *contentBuilder {
	return b.addRule(ruleRoadDistance(d, 1))
}
func (b *contentBuilder) addRule(r template.PlacementRule) *contentBuilder {
	b.item.Rules = append(b.item.Rules, r)
	return b
}
func (b *contentBuilder) addRules(rs []template.PlacementRule) *contentBuilder {
	b.item.Rules = append(b.item.Rules, rs...)
	return b
}
func (b *contentBuilder) build() template.MandatoryContentItem { return b.item }

// ── Content presets ──────────────────────────────────────────────────

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

func presetMineGoldNearCrossroads() template.MandatoryContentItem {
	return newContentBuilder(constants.ContentIds.MineGold.Sid).mine().
		addRule(ruleCrossroadsDistance(distNear, 1)).build()
}

func presetMineCrystalsNextToRoad() template.MandatoryContentItem {
	return newContentBuilder(constants.ContentIds.MineCrystals.Sid).withName("name_mine_crystals").mine().
		roadDistance(distNextTo).build()
}

func presetMineMercuryNextToRoad() template.MandatoryContentItem {
	return newContentBuilder(constants.ContentIds.MineMercury.Sid).withName("name_mine_mercury").mine().
		roadDistance(distNextTo).build()
}

func presetMineGemstonesNextToRoad() template.MandatoryContentItem {
	return newContentBuilder(constants.ContentIds.MineGemstones.Sid).withName("name_mine_gemstones").mine().
		roadDistance(distNextTo).build()
}

func presetAlchemyLabNearRoad() template.MandatoryContentItem {
	return newContentBuilder(constants.ContentIds.AlchemyLab.Sid).withName("name_alchemy_lab").mine().
		roadDistance(distNear).build()
}

// ── Zone content manager ─────────────────────────────────────────────

// BuildPlayerZoneMandatoryContent returns the mandatory content list for a player spawn zone.
func BuildPlayerZoneMandatoryContent(settings *models.GeneratorSettings) []template.MandatoryContentItem {
	var content []template.MandatoryContentItem

	if settings.SpawnRemoteFootholds {
		content = append(content, presetRemoteFoothold(settings.ZoneCfg.PlayerZoneCastles))
	}

	// Append user-configured items from the UI.
	content = append(content, settings.PlayerZoneMandatoryContent...)

	content = append(content,
		template.MandatoryContentItem{SID: "watchtower"},
		newContentBuilder(constants.ContentIds.Market.Sid).guarded().roadDistance(distNear).build(),
		newContentBuilder(constants.ContentIds.ManaWell.Sid).roadDistance(distNear).build(),
		template.MandatoryContentItem{IncludeLists: []string{"basic_content_list_building_hero_stats_and_skills_tier_2"}},
		template.MandatoryContentItem{IncludeLists: []string{"content_list_building_uncommon_hero_banks"}},
		template.MandatoryContentItem{IncludeLists: []string{"content_list_pickup_pandora_box_army_low_tier"}, IsGuarded: true},
	)
	return content
}

// BuildLowNeutralMandatoryContent returns mandatory content for a low-quality neutral zone.
func BuildLowNeutralMandatoryContent(castleCount int, spawnFootholds bool) []template.MandatoryContentItem {
	var content []template.MandatoryContentItem
	if spawnFootholds {
		content = append(content, presetRemoteFoothold(castleCount))
	}
	content = append(content,
		template.MandatoryContentItem{Name: "name_mine_by_biome_1", IncludeLists: []string{"basic_content_list_rare_mines_by_biome"}, IsMine: true},
		template.MandatoryContentItem{IncludeLists: []string{"basic_content_list_rare_mines"}, IsMine: true},
		newContentBuilder(constants.ContentIds.Market.Sid).guarded().addRule(ruleCrossroadsDistance(distNear, 1)).build(),
		template.MandatoryContentItem{IncludeLists: []string{"basic_content_list_vision_buildings_tier_1"}},
		template.MandatoryContentItem{IncludeLists: []string{"basic_content_list_building_hero_buff_tier_1"}},
		template.MandatoryContentItem{IncludeLists: []string{"basic_content_list_building_hero_buff_tier_1"}},
		template.MandatoryContentItem{IncludeLists: []string{"basic_content_list_building_hero_stats_and_skills_tier_1"}},
		template.MandatoryContentItem{IncludeLists: []string{"content_list_building_random_hires_low_tier"}},
		template.MandatoryContentItem{IncludeLists: []string{"content_list_building_random_hires_low_tier"}},
		newContentBuilder(constants.ContentIds.PandoraBox.Sid).soloEncounter().build(),
		template.MandatoryContentItem{IncludeLists: []string{"content_list_pickup_pandora_box_army_low_tier"}},
		template.MandatoryContentItem{IncludeLists: []string{"basic_content_list_pickup_random_items"}},
		template.MandatoryContentItem{IncludeLists: []string{"basic_content_list_building_magic_tier_1"}},
	)
	return content
}

// BuildMediumNeutralMandatoryContent returns mandatory content for a medium-quality neutral zone.
func BuildMediumNeutralMandatoryContent(castleCount int, spawnFootholds bool) []template.MandatoryContentItem {
	var content []template.MandatoryContentItem
	if spawnFootholds {
		content = append(content, presetRemoteFoothold(castleCount))
	}
	content = append(content,
		presetMineCrystalsNextToRoad(),
		presetMineMercuryNextToRoad(),
		presetMineGemstonesNextToRoad(),
		presetAlchemyLabNearRoad(),
		newContentBuilder(constants.ContentIds.MineGold.Sid).mine().roadDistance(distNear).build(),
		newContentBuilder(constants.ContentIds.Watchtower.Sid).guarded().build(),
		template.MandatoryContentItem{IncludeLists: []string{"basic_content_list_vision_buildings_tier_1"}},
		template.MandatoryContentItem{IncludeLists: []string{"basic_content_list_building_hero_buff_tier_1"}},
		template.MandatoryContentItem{IncludeLists: []string{"basic_content_list_building_hero_stats_and_skills_tier_1"}},
		template.MandatoryContentItem{IncludeLists: []string{"basic_content_list_building_hero_stats_and_skills_tier_2"}},
		template.MandatoryContentItem{IncludeLists: []string{"basic_content_list_building_magic_tier_1"}},
		template.MandatoryContentItem{IncludeLists: []string{"basic_content_list_building_magic_tier_2"}},
		template.MandatoryContentItem{IncludeLists: []string{"content_list_building_random_hires_low_tier"}},
		template.MandatoryContentItem{IncludeLists: []string{"content_list_building_random_hires_high_tier"}},
		template.MandatoryContentItem{IncludeLists: []string{"basic_content_list_building_guarded_units_banks_only_biome_restriction"}},
		template.MandatoryContentItem{IncludeLists: []string{"basic_content_list_building_guarded_resource_banks_tier_2"}},
		newContentBuilder(constants.ContentIds.RandomItemEpic.Sid).soloEncounter().build(),
		template.MandatoryContentItem{SID: constants.ContentIds.RandomItemEpic.Sid},
		newContentBuilder(constants.ContentIds.PandoraBox.Sid).soloEncounter().build(),
		template.MandatoryContentItem{SID: constants.ContentIds.PandoraBox.Sid},
		template.MandatoryContentItem{IncludeLists: []string{"content_list_pickup_pandora_box_army_low_tier"}},
	)
	return content
}

// BuildHighNeutralMandatoryContent returns mandatory content for a high-quality neutral zone.
func BuildHighNeutralMandatoryContent(castleCount int, spawnFootholds bool) []template.MandatoryContentItem {
	var content []template.MandatoryContentItem
	if spawnFootholds {
		content = append(content, presetRemoteFoothold(castleCount))
	}
	content = append(content,
		template.MandatoryContentItem{IncludeLists: []string{"content_list_building_utopia"}},
		template.MandatoryContentItem{IncludeLists: []string{"content_list_building_utopia"}},
		template.MandatoryContentItem{IncludeLists: []string{"content_list_building_epic_guarded_resource_banks"}},
		template.MandatoryContentItem{IncludeLists: []string{"content_list_building_epic_guarded_resource_banks"}},
		template.MandatoryContentItem{IncludeLists: []string{"basic_content_list_vision_buildings_tier_1"}},
		template.MandatoryContentItem{IncludeLists: []string{"basic_content_list_building_hero_buff_tier_1"}},
		template.MandatoryContentItem{IncludeLists: []string{"basic_content_list_building_hero_stats_and_skills_tier_2"}},
		template.MandatoryContentItem{IncludeLists: []string{"basic_content_list_building_hero_stats_and_skills_tier_3"}},
		template.MandatoryContentItem{IncludeLists: []string{"basic_content_list_building_hero_stats_and_skills_tier_3"}},
		template.MandatoryContentItem{IncludeLists: []string{"basic_content_list_building_magic_tier_2"}},
		template.MandatoryContentItem{IncludeLists: []string{"basic_content_list_building_magic_tier_2"}},
		template.MandatoryContentItem{IncludeLists: []string{"content_list_building_random_hires_high_tier"}},
		template.MandatoryContentItem{IncludeLists: []string{"content_list_building_random_hires_high_tier"}},
		template.MandatoryContentItem{IncludeLists: []string{"basic_content_list_building_random_hires"}},
		template.MandatoryContentItem{IncludeLists: []string{"basic_content_list_building_guarded_units_banks_only_biome_restriction"}},
		template.MandatoryContentItem{IncludeLists: []string{"basic_content_list_building_guarded_units_banks_no_biome_restriction"}},
		template.MandatoryContentItem{IncludeLists: []string{"basic_content_list_building_guarded_units_banks_no_biome_restriction"}},
		template.MandatoryContentItem{IncludeLists: []string{"basic_content_list_building_guarded_resource_banks_tier_2"}},
		template.MandatoryContentItem{IncludeLists: []string{"basic_content_list_building_guarded_resource_banks_tier_3"}},
		template.MandatoryContentItem{IncludeLists: []string{"basic_content_list_pickup_mythic_scroll_box"}},
		template.MandatoryContentItem{IncludeLists: []string{"basic_content_list_pickup_mythic_scroll_box"}},
		newContentBuilder(constants.ContentIds.RandomItemLegendary.Sid).soloEncounter().build(),
		template.MandatoryContentItem{SID: constants.ContentIds.RandomItemLegendary.Sid},
		template.MandatoryContentItem{SID: constants.ContentIds.RandomItemEpic.Sid},
		newContentBuilder(constants.ContentIds.PandoraBox.Sid).soloEncounter().build(),
		template.MandatoryContentItem{SID: constants.ContentIds.PandoraBox.Sid},
		template.MandatoryContentItem{SID: constants.ContentIds.PandoraBox.Sid},
		template.MandatoryContentItem{IncludeLists: []string{"content_list_pickup_pandora_box_army_high_tier"}},
		template.MandatoryContentItem{IncludeLists: []string{"content_list_pickup_pandora_box_army_high_tier"}},
		presetMineGoldNearCrossroads(),
		template.MandatoryContentItem{SID: constants.ContentIds.MineGold.Sid, IsMine: true},
		template.MandatoryContentItem{SID: constants.ContentIds.MineGold.Sid, IsMine: true},
		template.MandatoryContentItem{SID: constants.ContentIds.MineCrystals.Sid, IsMine: true},
		template.MandatoryContentItem{SID: constants.ContentIds.MineMercury.Sid, IsMine: true},
		template.MandatoryContentItem{SID: constants.ContentIds.MineGemstones.Sid, IsMine: true},
		template.MandatoryContentItem{SID: constants.ContentIds.AlchemyLab.Sid, IsMine: true},
		template.MandatoryContentItem{SID: constants.ContentIds.AlchemyLab.Sid, IsMine: true},
	)
	return content
}

// ── Content count limits ─────────────────────────────────────────────

// BuildAllContentCountLimits returns the full set of content count limits.
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

	// Lift limits when player mandatory content has more items than the default.
	sidCounts := map[string]int{}
	for _, item := range settings.PlayerZoneMandatoryContent {
		if item.SID != "" {
			sidCounts[strings.ToLower(item.SID)]++
		}
	}
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
