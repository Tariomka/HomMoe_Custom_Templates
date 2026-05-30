package generator

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
)

// TODO: Make generator a class
// TODO: Use strategy pattern for building topology specific layouts
// TODO: Split this shit and organize everything

const (
	defaultGuardRandomization = 0.05
	spawnLayoutName           = "zone_layout_spawns"
	sideLayoutName            = "zone_layout_sides"
	treasureLayoutName        = "zone_layout_treasure_zone"
	centerLayoutName          = "zone_layout_center"
)

// distance presets
type distancePreset struct{ Min, Max float64 }

var distNear = distancePreset{0.0, 0.35}

func roundToDP(v float64, dp int) float64 {
	mul := math.Pow(10, float64(dp))
	return math.Round(v*mul) / mul
}

func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// ZoneLetters are the labels used to name zones (max 32).
var ZoneLetters = []string{
	"A", "B", "C", "D", "E", "F", "G", "H",
	"I", "J", "K", "L", "M", "N", "O", "P",
	"Q", "R", "S", "T", "U", "V", "W", "X",
	"Y", "Z", "AA", "AB", "AC", "AD", "AE", "AF",
}

// Generate produces an RmgTemplate from the given settings.
func Generate(generatorSettings *config.GeneratorConfig) (*template.RmgTemplateModel, error) {
	if generatorSettings.TemplateName == "" {
		return nil, fmt.Errorf("template name is required")
	}
	playerLetters := make([]string, generatorSettings.PlayerCount)
	copy(playerLetters, ZoneLetters[:generatorSettings.PlayerCount])

	neutralZones := buildNeutralZonePlan(generatorSettings)
	useCityHold := generatorSettings.GameEndConditions != nil &&
		(generatorSettings.GameEndConditions.CityHold || generatorSettings.GameEndConditions.VictoryCondition == "win_condition_5")
	var holdCityNeutralLetter string
	if useCityHold && generatorSettings.Topology != config.TopologyHubAndSpoke {
		adj := buildTopologyAdjacency(generatorSettings, playerLetters, neutralZones)
		holdCityNeutralLetter = pickHoldCityNeutralLetter(neutralZones, playerLetters, adj)
	}

	totalZones := generatorSettings.PlayerCount + len(neutralZones)
	tuning := generationTuning{
		ContentScale:                   computeContentScale(generatorSettings.MapSize, totalZones),
		ResourceDensityMultiplier:      float64(generatorSettings.ZoneConfiguration.ResourceDensityPercent) / 200.0,
		StructureDensityMultiplier:     float64(generatorSettings.ZoneConfiguration.StructureDensityPercent) / 100.0,
		NeutralStackStrengthMultiplier: float64(generatorSettings.ZoneConfiguration.NeutralStackStrengthPercent) / 100.0,
		BorderGuardStrengthMultiplier:  float64(generatorSettings.ZoneConfiguration.BorderGuardStrengthPercent) / 100.0,
		GuardRandomization:             effectiveGuardRandomization(generatorSettings),
	}

	effectiveVC := "win_condition_1"
	if generatorSettings.GameEndConditions != nil {
		effectiveVC = generatorSettings.GameEndConditions.VictoryCondition
	}

	tmpl := &template.RmgTemplateModel{
		Name:                generatorSettings.TemplateName,
		GameMode:            generatorSettings.GameMode,
		Description:         buildTemplateDescription(generatorSettings, len(neutralZones)),
		DisplayWinCondition: effectiveVC,
		SizeX:               generatorSettings.MapSize,
		SizeZ:               generatorSettings.MapSize,
		GameRules:           buildGameRules(generatorSettings, effectiveVC),
		Variants:            []template.Variant{buildVariant(generatorSettings, playerLetters, neutralZones, tuning, holdCityNeutralLetter, useCityHold && generatorSettings.Topology == config.TopologyHubAndSpoke)},
		ZoneLayouts:         buildZoneLayouts(),
		MandatoryContent:    buildAllMandatoryContent(playerLetters, neutralZones, generatorSettings),
		ContentCountLimits:  BuildAllContentCountLimits(generatorSettings),
		ContentPools:        []template.ContentPool{},
		ContentLists:        []template.ContentList{},
	}
	return tmpl, nil
}

// ── description ──────────────────────────────────────────────────────

func buildTemplateDescription(settings *config.GeneratorConfig, neutralCount int) string {
	parts := []string{
		constants.GetTopologyDescriptor(settings.Topology).Label + " layout",
		countPhrase(neutralCount, "neutral zone", "neutral zones"),
		countPhrase(settings.ZoneConfiguration.PlayerZoneCastles, "castle", "castles") + " per player zone",
	}
	if neutralCount > 0 {
		if settings.ZoneConfiguration.Advanced.Enabled {
			parts = append(parts, "mixed neutral zone tiers")
		} else {
			parts = append(parts, countPhrase(settings.ZoneConfiguration.NeutralZoneCastles, "castle", "castles")+" per neutral zone")
		}
	}
	var options []string
	if settings.NoDirectPlayerConnections {
		options = append(options, "isolated player starts")
	}
	if settings.RandomPortals {
		options = append(options, "random portals")
	}
	if !settings.SpawnRemoteFootholds {
		options = append(options, "no remote footholds")
	}
	if !settings.GenerateRoads {
		options = append(options, "roads disabled")
	}
	if len(options) > 0 {
		parts = append(parts, "options: "+strings.Join(options, ", "))
	}
	return "Generated with Olden Era Template Generator: " + strings.Join(parts, ", ") + "."
}

func countPhrase(count int, singular, plural string) string {
	if count == 0 {
		return "no " + plural
	}
	word := singular
	if count != 1 {
		word = plural
	}
	return fmt.Sprintf("%d %s", count, word)
}

type neutralZonePlan struct {
	Letter      string
	Quality     constants.NeutralZoneQuality
	CastleCount int
}

func buildNeutralZonePlan(generatorSettings *config.GeneratorConfig) []neutralZonePlan {
	var plans []neutralZonePlan
	maxNeutral := max(0, len(ZoneLetters)-generatorSettings.PlayerCount)
	castleZoneCastleCount := clampInt(generatorSettings.ZoneConfiguration.NeutralZoneCastles, 1, 4)

	add := func(requested int, quality constants.NeutralZoneQuality, castleCount int) {
		count := clampInt(requested, 0, 30)
		for i := 0; i < count && len(plans) < maxNeutral; i++ {
			letter := ZoneLetters[generatorSettings.PlayerCount+len(plans)]
			plans = append(plans, neutralZonePlan{letter, quality, castleCount})
		}
	}

	advanced := generatorSettings.ZoneConfiguration.Advanced
	advancedTotal := advanced.NeutralLowNoCastleCount + advanced.NeutralLowCastleCount +
		advanced.NeutralMediumNoCastleCount + advanced.NeutralMediumCastleCount +
		advanced.NeutralHighNoCastleCount + advanced.NeutralHighCastleCount
	if advancedTotal > 0 {
		add(advanced.NeutralLowNoCastleCount, constants.QualityLow, 0)
		add(advanced.NeutralLowCastleCount, constants.QualityLow, castleZoneCastleCount)
		add(advanced.NeutralMediumNoCastleCount, constants.QualityMedium, 0)
		add(advanced.NeutralMediumCastleCount, constants.QualityMedium, castleZoneCastleCount)
		add(advanced.NeutralHighNoCastleCount, constants.QualityHigh, 0)
		add(advanced.NeutralHighCastleCount, constants.QualityHigh, castleZoneCastleCount)
	} else {
		cc := clampInt(generatorSettings.ZoneConfiguration.NeutralZoneCastles, 0, 4)
		add(generatorSettings.ZoneConfiguration.NeutralZoneCount, constants.QualityMedium, cc)
	}
	if generatorSettings.Topology == config.TopologySharedWeb && len(plans) == 0 && maxNeutral > 0 {
		letter := ZoneLetters[generatorSettings.PlayerCount]
		cc := clampInt(generatorSettings.ZoneConfiguration.NeutralZoneCastles, 0, 4)
		plans = append(plans, neutralZonePlan{letter, constants.QualityMedium, cc})
	}
	return plans
}

type generationTuning struct {
	ContentScale                   float64
	ResourceDensityMultiplier      float64
	StructureDensityMultiplier     float64
	NeutralStackStrengthMultiplier float64
	BorderGuardStrengthMultiplier  float64
	GuardRandomization             float64
}

func scaleValue(value, multiplier float64) int {
	return max(0, int(value*multiplier))
}
func scaleStructureValue(value float64, t generationTuning) int {
	return scaleValue(value, t.StructureDensityMultiplier)
}
func scaleResourceValue(value float64, t generationTuning) int {
	return scaleValue(value, t.ResourceDensityMultiplier)
}
func scaleNeutralGuardValue(value int, t generationTuning) int {
	return scaleValue(float64(value), t.NeutralStackStrengthMultiplier)
}
func scaleBorderGuardValue(value int, t generationTuning) int {
	return scaleValue(float64(value), t.BorderGuardStrengthMultiplier)
}
func scaleGuardMultiplier(value float64, t generationTuning) float64 {
	return roundToDP(value*t.NeutralStackStrengthMultiplier, 3)
}

// borderGuardValue returns the base border-guard value (already scaled
// by the BorderGuardStrengthMultiplier tuning) for a connection between
// two zone letters
func borderGuardValue(letterA, letterB string, playerLetters []string, neutralByLetter map[string]neutralZonePlan, t generationTuning) int {
	aIsPlayer := contains(playerLetters, letterA)
	bIsPlayer := contains(playerLetters, letterB)
	if aIsPlayer && bIsPlayer {
		return scaleBorderGuardValue(30000, t)
	}
	if !aIsPlayer && !bIsPlayer {
		qa := neutralQualityOf(neutralByLetter, letterA)
		qb := neutralQualityOf(neutralByLetter, letterB)
		higher := qa
		if int(qb) > int(qa) {
			higher = qb
		}
		return scaleBorderGuardValue(qualityGuardBase(higher), t)
	}
	neutralLetter := letterB
	if !aIsPlayer {
		neutralLetter = letterA
	}
	return scaleBorderGuardValue(qualityGuardBase(neutralQualityOf(neutralByLetter, neutralLetter)), t)
}

func neutralQualityOf(neutralByLetter map[string]neutralZonePlan, letter string) constants.NeutralZoneQuality {
	if neutralByLetter == nil {
		return constants.QualityMedium
	}
	plan, ok := neutralByLetter[letter]
	if !ok {
		return constants.QualityMedium
	}
	return plan.Quality
}

func qualityGuardBase(q constants.NeutralZoneQuality) int {
	switch q {
	case constants.QualityLow:
		return 15000
	case constants.QualityHigh:
		return 25000
	default:
		return 20000
	}
}

func effectiveGuardRandomization(settings *config.GeneratorConfig) float64 {
	if !settings.ZoneConfiguration.Advanced.Enabled {
		return defaultGuardRandomization
	}
	v := settings.ZoneConfiguration.Advanced.GuardRandomization
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return defaultGuardRandomization
	}
	return roundToDP(math.Max(0, math.Min(v, 0.5)), 3)
}

func computeContentScale(mapSize, totalZones int) float64 {
	const referenceArea = 160.0 * 160.0 / 4.0
	zoneArea := float64(mapSize) * float64(mapSize) / math.Max(1, float64(totalZones))
	return math.Max(0.5, math.Min(2.5, math.Sqrt(zoneArea/referenceArea)))
}

func normalizeZoneSize(zoneSize float64) float64 {
	if math.IsNaN(zoneSize) || math.IsInf(zoneSize, 0) {
		return 1.0
	}
	return roundToDP(math.Max(0.1, math.Min(zoneSize, 2.0)), 2)
}

func percentToModifier(percent int) float64 {
	return roundToDP(float64(clampInt(percent, 25, 200))/100.0, 2)
}

func buildGameRules(settings *config.GeneratorConfig, effectiveVC string) template.GameRules {
	hs := settings.HeroSettings
	isSingleHero := settings.GameMode == "SingleHero"
	if isSingleHero {
		hs.HeroCountMin = 1
		hs.HeroCountMax = 1
		hs.HeroCountIncrement = 1
	}
	return template.GameRules{
		HeroCountMin:           hs.HeroCountMin,
		HeroCountMax:           hs.HeroCountMax,
		HeroCountIncrement:     hs.HeroCountIncrement,
		HeroHireBan:            isSingleHero,
		EncounterHoles:         false,
		FactionLawsExpModifier: percentToModifier(settings.FactionLawsExpPercent),
		AstrologyExpModifier:   percentToModifier(settings.AstrologyExpPercent),
		Bonuses: template.BonusList{
			{SID: "add_bonus_hero_stat", ReceiverSide: -1, ReceiverFilter: "all_heroes", Parameters: []string{"movementBonus", "0"}},
		},
		WinConditions: buildAdvancedWinConditions(settings, effectiveVC),
	}
}

func buildAdvancedWinConditions(generatorSettings *config.GeneratorConfig, effectiveVC string) template.WinConditions {
	ge := generatorSettings.GameEndConditions
	if ge == nil {
		ge = &config.GameEndConditions{VictoryCondition: "win_condition_1", LostStartCityDay: 3, CityHoldDays: 6}
	}
	gr := generatorSettings.GladiatorArenaRules
	if gr == nil {
		gr = &config.GladiatorArenaRules{}
	}
	tr := generatorSettings.TournamentRules
	if tr == nil {
		tr = &config.TournamentRules{}
	}

	useLostStartCity := ge.LostStartCity || effectiveVC == "win_condition_3"
	useCityHold := ge.CityHold || effectiveVC == "win_condition_5"
	useGladiator := gr.Enabled || effectiveVC == "win_condition_4"
	useTournament := tr.Enabled || effectiveVC == "win_condition_6"

	wc := template.WinConditions{
		Classic:          true,
		Desertion:        true,
		DesertionDay:     3,
		DesertionValue:   3000,
		HeroLighting:     true,
		HeroLightingDay:  1,
		LostStartCity:    useLostStartCity,
		LostStartCityDay: clampInt(ge.LostStartCityDay, 1, 30),
		LostStartHero:    ge.LostStartHero || useGladiator || generatorSettings.GameMode == "SingleHero",
		CityHold:         useCityHold,
		CityHoldDays:     clampInt(ge.CityHoldDays, 1, 30),
	}
	if useGladiator {
		wc.GladiatorArena = true
		wc.GladiatorArenaRegistrationStartFight = true
		wc.GladiatorArenaDaysDelayStart = clampInt(gr.DaysDelayStart, 1, 60)
		wc.GladiatorArenaCountDay = clampInt(gr.CountDay, 1, 30)
		wc.ChampionSelectRule = "StartHero"
	}
	if useTournament {
		firstDay := clampInt(tr.FirstTournamentDay, 3, 60)
		interval := clampInt(tr.Interval, 3, 30)
		pointsToWin := clampInt(tr.PointsToWin, 1, 10)
		roundCount := pointsToWin*2 - 1
		wc.ChampionSelectRule = "StartHero"
		wc.Tournament = true
		wc.TournamentSaveArmy = true
		wc.TournamentPointsToWin = pointsToWin

		var announceDays, battleOffsets []int
		prevBattle := 0
		for i := range roundCount {
			announce := 1
			if i > 0 {
				announce = prevBattle + 1
			}
			offset := firstDay - 1
			if i > 0 {
				offset = interval - 1
			}
			announceDays = append(announceDays, announce)
			battleOffsets = append(battleOffsets, offset)
			prevBattle = announce + offset
		}
		wc.TournamentAnnounceDays = announceDays
		wc.TournamentDays = battleOffsets
	}
	return wc
}

func buildVariant(generatorSettings *config.GeneratorConfig, playerLetters []string, neutralZones []neutralZonePlan, tuning generationTuning, holdCityNeutralLetter string, hubIsHoldCity bool) template.Variant {
	// Shuffle player letters
	pl := make([]string, len(playerLetters))
	copy(pl, playerLetters)
	rand.Shuffle(len(pl), func(i, j int) { pl[i], pl[j] = pl[j], pl[i] })

	isTournament := (generatorSettings.TournamentRules != nil && generatorSettings.TournamentRules.Enabled) ||
		(generatorSettings.GameEndConditions != nil && generatorSettings.GameEndConditions.VictoryCondition == "win_condition_6")
	if isTournament && len(pl) == 2 {
		return buildVariantTournament(generatorSettings, pl, neutralZones, tuning)
	}

	switch generatorSettings.Topology {
	case config.TopologyHubAndSpoke:
		return buildVariantHubAndSpoke(generatorSettings, pl, neutralZones, tuning, hubIsHoldCity)
	case config.TopologyChain:
		return buildVariantChain(generatorSettings, pl, neutralZones, tuning, holdCityNeutralLetter)
	case config.TopologySharedWeb:
		return buildVariantSharedWeb(generatorSettings, pl, neutralZones, tuning, holdCityNeutralLetter)
	case config.TopologyRandom:
		return buildVariantRandom(generatorSettings, pl, neutralZones, tuning, holdCityNeutralLetter)
	case config.TopologyBalanced:
		return buildVariantBalanced(generatorSettings, pl, neutralZones, tuning, holdCityNeutralLetter)
	default:
		return buildVariantDefault(generatorSettings, pl, neutralZones, tuning, holdCityNeutralLetter)
	}
}

var (
	t2Guarded   = []string{"classic_template_pool_random_t2_item", "classic_template_pool_random_t2_pandora", "classic_template_pool_random_t2_hire", "classic_template_pool_random_t2_unit_bank", "classic_template_pool_random_t2_res_bank", "classic_template_pool_random_t2_stat", "classic_template_pool_random_t2_magic"}
	t2Unguarded = []string{"classic_template_pool_random_unguarded_t2_item", "classic_template_pool_random_unguarded_t2_pandora", "classic_template_pool_random_unguarded_t2_hire", "classic_template_pool_random_unguarded_t2_unit_bank", "classic_template_pool_random_unguarded_t2_res_bank", "classic_template_pool_random_unguarded_t2_stat", "classic_template_pool_random_unguarded_t2_magic"}
	t3Guarded   = []string{"classic_template_pool_random_t3_item", "classic_template_pool_random_t3_pandora", "classic_template_pool_random_t3_hire", "classic_template_pool_random_t3_unit_bank", "classic_template_pool_random_t3_res_bank", "classic_template_pool_random_t3_stat", "classic_template_pool_random_t3_magic"}
	t3Unguarded = []string{"classic_template_pool_random_unguarded_t3_item", "classic_template_pool_random_unguarded_t3_pandora", "classic_template_pool_random_unguarded_t3_hire", "classic_template_pool_random_unguarded_t3_unit_bank", "classic_template_pool_random_unguarded_t3_res_bank", "classic_template_pool_random_unguarded_t3_stat", "classic_template_pool_random_unguarded_t3_magic"}
	t4Guarded   = []string{"classic_template_pool_random_t4_item", "classic_template_pool_random_t4_pandora", "classic_template_pool_random_t4_hire", "classic_template_pool_random_t4_unit_bank", "classic_template_pool_random_t4_res_bank", "classic_template_pool_random_t4_stat", "classic_template_pool_random_t4_magic"}
	t4Unguarded = []string{"classic_template_pool_random_unguarded_t4_item", "classic_template_pool_random_unguarded_t4_pandora", "classic_template_pool_random_unguarded_t4_hire", "classic_template_pool_random_unguarded_t4_unit_bank", "classic_template_pool_random_unguarded_t4_res_bank", "classic_template_pool_random_unguarded_t4_stat", "classic_template_pool_random_unguarded_t4_magic"}
	t5Guarded   = []string{"classic_template_pool_random_t5_item", "classic_template_pool_random_t5_pandora", "classic_template_pool_random_t5_hire", "classic_template_pool_random_t5_unit_bank", "classic_template_pool_random_t5_res_bank", "classic_template_pool_random_t5_stat", "classic_template_pool_random_t5_magic"}
	t5Unguarded = []string{"classic_template_pool_random_unguarded_t5_item", "classic_template_pool_random_unguarded_t5_pandora", "classic_template_pool_random_unguarded_t5_hire", "classic_template_pool_random_unguarded_t5_unit_bank", "classic_template_pool_random_unguarded_t5_res_bank", "classic_template_pool_random_unguarded_t5_stat", "classic_template_pool_random_unguarded_t5_magic"}

	generalResourcesPoor   = []string{"content_pool_general_resources_start_zone_poor"}
	generalResourcesMedium = []string{"content_pool_general_resources_start_zone_medium"}
	generalResourcesRich   = []string{"content_pool_general_resources_start_zone_rich"}
)

type neutralZoneProfile struct {
	Layout                       string
	GuardMultiplier              float64
	GuardedContentPool           []string
	UnguardedContentPool         []string
	ResourcesContentPool         []string
	GuardedContentValue          int
	GuardedContentValuePerArea   int
	UnguardedContentValue        int
	UnguardedContentValuePerArea int
	ResourcesValue               int
	ResourcesValuePerArea        int
	PrimaryCityGuardValue        int
	ExtraCityGuardValue          int
	PrimaryBuildingsCSid         string
	ExtraBuildingsCSid           string
}

func getNeutralZoneProfile(quality constants.NeutralZoneQuality) neutralZoneProfile {
	switch quality {
	case constants.QualityLow:
		return neutralZoneProfile{sideLayoutName, 1.1, cp(t2Guarded), cp(t2Unguarded), cp(generalResourcesPoor), 120000, 1000, 25000, 200, 30000, 240, 4000, 2000, "poor_buildings_construction", "poor_buildings_construction"}
	case constants.QualityHigh:
		return neutralZoneProfile{treasureLayoutName, 1.8, append(cp(t4Guarded), t5Guarded...), append(cp(t4Unguarded), t5Unguarded...), cp(generalResourcesRich), 480000, 3000, 80000, 620, 90000, 580, 16000, 8000, "rich_buildings_construction", "rich_buildings_construction"}
	default: // Medium
		return neutralZoneProfile{treasureLayoutName, 1.4, cp(t3Guarded), cp(t3Unguarded), cp(generalResourcesMedium), 240000, 2000, 38000, 300, 55000, 420, 8000, 4000, "rich_buildings_construction", "poor_buildings_construction"}
	}
}

func cp(s []string) []string { r := make([]string, len(s)); copy(r, s); return r }

func buildSpawnZone(letter, player string, ringConns []string, castleCount int, matchFactions bool, zoneSize float64, spawnFootholds, generateRoads bool, tuning generationTuning) template.Zone {
	mainObjects := []template.MainObject{
		{
			Type: "Spawn", Spawn: player, RemoveGuardIfHasOwner: true,
			GuardChance: 1, GuardValue: scaleNeutralGuardValue(5000, tuning), GuardWeeklyIncrement: 0.10,
			BuildingsConstructionSid: "default_buildings_construction",
			Placement:                "Uniform", PlacementArgs: []string{"true", "0.7", "0"},
		},
	}
	for i := 1; i < castleCount; i++ {
		// TODO: add player own castles
		mo := template.MainObject{
			Type: "City", GuardChance: 1, GuardValue: scaleNeutralGuardValue(2500, tuning), GuardWeeklyIncrement: 0.10,
			BuildingsConstructionSid: "poor_buildings_construction",
			Placement:                "Uniform", PlacementArgs: []string{"false", "-0.8", "3"},
		}
		if matchFactions {
			mo.Faction = &template.TypedRef{Type: "Match", Args: []string{"0"}}
		} else {
			mo.Faction = &template.TypedRef{Type: "Random", Args: []string{}}
		}
		mainObjects = append(mainObjects, mo)
	}
	cp0 := 0
	biome := template.TypedRef{Type: "MatchMainObject", Args: []string{"0"}}
	var roads []template.Road
	if castleCount > 0 {
		roads = buildOuterZoneRoads(ringConns, castleCount, spawnFootholds, generateRoads)
	} else {
		roads = buildConnectorZoneRoads(ringConns, generateRoads)
	}
	return template.Zone{
		Name: "Spawn-" + letter, Size: normalizeZoneSize(zoneSize), Layout: spawnLayoutName,
		GuardCutoffValue: 2000, GuardRandomization: tuning.GuardRandomization,
		GuardMultiplier: scaleGuardMultiplier(1.0, tuning), GuardWeeklyIncrement: 0.20,
		GuardReactionDistribution: []int{60, 20, 10, 10, 2, 0}, DiplomacyModifier: -0.5,
		GuardedContentPool: cp(t2Guarded), UnguardedContentPool: cp(t2Unguarded), ResourcesContentPool: cp(generalResourcesPoor),
		MandatoryContent:             template.StringList{"mandatory_content_side_" + letter},
		ContentCountLimits:           buildSideContentLimits(),
		GuardedContentValue:          scaleStructureValue(200000*tuning.ContentScale, tuning),
		GuardedContentValuePerArea:   scaleStructureValue(2000*math.Sqrt(tuning.ContentScale), tuning),
		UnguardedContentValue:        scaleStructureValue(50000*tuning.ContentScale, tuning),
		UnguardedContentValuePerArea: scaleStructureValue(400*math.Sqrt(tuning.ContentScale), tuning),
		ResourcesValue:               scaleResourceValue(80000*tuning.ContentScale, tuning),
		ResourcesValuePerArea:        scaleResourceValue(600*math.Sqrt(tuning.ContentScale), tuning),
		MainObjects:                  mainObjects, ZoneBiome: biome, ContentBiome: biome, MetaObjectsBiome: biome,
		CrossroadsPosition: &cp0, Roads: roads,
	}
}

func buildNeutralZone(plan neutralZonePlan, ringConns []string, zoneSize float64, spawnFootholds, generateRoads bool, tuning generationTuning, isHoldCity bool) template.Zone {
	letter := plan.Letter
	castleCount := plan.CastleCount
	if isHoldCity && castleCount < 1 {
		castleCount = 1
	}
	profile := getNeutralZoneProfile(plan.Quality)

	var mainObjects []template.MainObject
	if castleCount > 0 {
		guardVal := profile.PrimaryCityGuardValue
		bcsid := profile.PrimaryBuildingsCSid
		placement := "Uniform"
		var placementArgs []string
		if isHoldCity {
			if guardVal < 20000 {
				guardVal = 20000
			}
			bcsid = "ultra_rich_buildings_construction"
			placement = "Center"
		} else {
			placementArgs = []string{"true", "0.8", "2"}
		}
		mo := template.MainObject{
			Type: "City", GuardChance: 1, GuardValue: scaleNeutralGuardValue(guardVal, tuning), GuardWeeklyIncrement: 0.10,
			BuildingsConstructionSid: bcsid,
			Faction:                  &template.TypedRef{Type: "FromList", Args: []string{}},
			Placement:                placement, PlacementArgs: placementArgs, HoldCityWinCon: isHoldCity,
		}
		mainObjects = append(mainObjects, mo)
	}
	for i := 1; i < castleCount; i++ {
		mainObjects = append(mainObjects, template.MainObject{
			Type: "City", GuardChance: 1, GuardValue: scaleNeutralGuardValue(profile.ExtraCityGuardValue, tuning), GuardWeeklyIncrement: 0.10,
			BuildingsConstructionSid: profile.ExtraBuildingsCSid,
			Faction:                  &template.TypedRef{Type: "FromList", Args: []string{}},
			Placement:                "Uniform", PlacementArgs: []string{"false", "-0.8", "3"},
		})
	}

	reaction := []int{0, 10, 10, 10, 10, 0}
	if plan.Quality == constants.QualityHigh {
		reaction = []int{0, 10, 10, 20, 10, 0}
	}

	biome := template.TypedRef{Type: "MatchZone", Args: []string{}}
	if castleCount > 0 {
		biome = template.TypedRef{Type: "MatchMainObject", Args: []string{"0"}}
	}
	cp0 := 0
	var roads []template.Road
	if castleCount > 0 {
		roads = buildOuterZoneRoads(ringConns, castleCount, spawnFootholds, generateRoads)
	} else {
		roads = buildConnectorZoneRoads(ringConns, generateRoads)
	}
	return template.Zone{
		Name: "Neutral-" + letter, Size: normalizeZoneSize(zoneSize), Layout: profile.Layout,
		GuardCutoffValue: 2000, GuardRandomization: tuning.GuardRandomization,
		GuardMultiplier: scaleGuardMultiplier(profile.GuardMultiplier, tuning), GuardWeeklyIncrement: 0.20,
		GuardReactionDistribution: reaction, DiplomacyModifier: -0.5,
		GuardedContentPool: profile.GuardedContentPool, UnguardedContentPool: profile.UnguardedContentPool, ResourcesContentPool: profile.ResourcesContentPool,
		MandatoryContent:             template.StringList{"mandatory_content_neutral_" + letter},
		ContentCountLimits:           buildSideContentLimits(),
		GuardedContentValue:          scaleStructureValue(float64(profile.GuardedContentValue)*tuning.ContentScale, tuning),
		GuardedContentValuePerArea:   scaleStructureValue(float64(profile.GuardedContentValuePerArea)*math.Sqrt(tuning.ContentScale), tuning),
		UnguardedContentValue:        scaleStructureValue(float64(profile.UnguardedContentValue)*tuning.ContentScale, tuning),
		UnguardedContentValuePerArea: scaleStructureValue(float64(profile.UnguardedContentValuePerArea)*math.Sqrt(tuning.ContentScale), tuning),
		ResourcesValue:               scaleResourceValue(float64(profile.ResourcesValue)*tuning.ContentScale, tuning),
		ResourcesValuePerArea:        scaleResourceValue(float64(profile.ResourcesValuePerArea)*math.Sqrt(tuning.ContentScale), tuning),
		MainObjects:                  mainObjects, ZoneBiome: biome, ContentBiome: biome, MetaObjectsBiome: biome,
		CrossroadsPosition: &cp0, Roads: roads,
	}
}

func buildHubZone(spokeConns []string, tuning generationTuning, isHoldCity bool, size float64, castleCount int, generateRoads bool) template.Zone {
	effectiveCastleCount := castleCount
	if isHoldCity && effectiveCastleCount < 1 {
		effectiveCastleCount = 1
	}
	var mainObjects []template.MainObject
	for i := 0; i < effectiveCastleCount; i++ {
		isHoldSlot := isHoldCity && i == 0
		gv := 16000
		if isHoldSlot {
			gv = 25000
		}
		gc := 0.5
		if isHoldSlot {
			gc = 1.0
		}
		bcsid := "rich_buildings_construction"
		if isHoldSlot {
			bcsid = "ultra_rich_buildings_construction"
		}
		placement := "Uniform"
		var placementArgs []string
		if isHoldSlot {
			placement = "Center"
		} else {
			placementArgs = []string{"true", "0.8", "2"}
		}
		mainObjects = append(mainObjects, template.MainObject{
			Type: "City", GuardChance: gc, GuardValue: scaleNeutralGuardValue(gv, tuning), GuardWeeklyIncrement: 0.10,
			BuildingsConstructionSid: bcsid,
			Faction:                  &template.TypedRef{Type: "FromList", Args: []string{}},
			Placement:                placement, PlacementArgs: placementArgs, HoldCityWinCon: isHoldSlot,
		})
	}
	biome := template.TypedRef{Type: "MatchZone", Args: []string{}}
	if effectiveCastleCount > 0 {
		biome = template.TypedRef{Type: "MatchMainObject", Args: []string{"0"}}
	}
	cp0 := 0
	var roads []template.Road
	if effectiveCastleCount > 0 {
		roads = buildOuterZoneRoads(spokeConns, effectiveCastleCount, false, generateRoads)
	} else {
		roads = buildConnectorZoneRoads(spokeConns, generateRoads)
	}
	return template.Zone{
		Name: "Hub", Size: size, Layout: centerLayoutName,
		GuardCutoffValue: 2000, GuardRandomization: 0.05,
		GuardMultiplier: scaleGuardMultiplier(1.5, tuning), GuardWeeklyIncrement: 0.20,
		GuardReactionDistribution: []int{0, 10, 10, 20, 10, 0}, DiplomacyModifier: -0.5,
		GuardedContentPool: cp(t3Guarded), UnguardedContentPool: cp(t3Unguarded), ResourcesContentPool: cp(generalResourcesMedium),
		MandatoryContent: template.StringList{}, ContentCountLimits: buildSideContentLimits(),
		GuardedContentValue:          scaleStructureValue(300000*tuning.ContentScale, tuning),
		GuardedContentValuePerArea:   scaleStructureValue(2400*math.Sqrt(tuning.ContentScale), tuning),
		UnguardedContentValue:        scaleStructureValue(50000*tuning.ContentScale, tuning),
		UnguardedContentValuePerArea: scaleStructureValue(600*math.Sqrt(tuning.ContentScale), tuning),
		ResourcesValue:               scaleResourceValue(80000*tuning.ContentScale, tuning),
		ResourcesValuePerArea:        scaleResourceValue(600*math.Sqrt(tuning.ContentScale), tuning),
		MainObjects:                  mainObjects, ZoneBiome: biome, ContentBiome: biome, MetaObjectsBiome: biome,
		CrossroadsPosition: &cp0, Roads: roads,
	}
}

func mainObjectEndpoint(index string) template.TypedRef {
	return template.TypedRef{Type: "MainObject", Args: []string{index}}
}
func connectionEndpoint(name string) template.TypedRef {
	return template.TypedRef{Type: "Connection", Args: []string{name}}
}
func mandatoryContentEndpoint(name string) template.TypedRef {
	return template.TypedRef{Type: "MandatoryContent", Args: []string{name}}
}
func plainRoad(from, to template.TypedRef) template.Road { return template.Road{From: from, To: to} }

func buildOuterZoneRoads(ringConns []string, castleCount int, includeFoothold, generateRoads bool) []template.Road {
	if !generateRoads || castleCount == 0 {
		return nil
	}
	var roads []template.Road
	for i := 1; i < castleCount; i++ {
		roads = append(roads, plainRoad(mainObjectEndpoint("0"), mainObjectEndpoint(fmt.Sprintf("%d", i))))
	}
	if includeFoothold {
		roads = append(roads, plainRoad(mainObjectEndpoint("0"), mandatoryContentEndpoint("name_remote_foothold_1")))
	}
	for _, rc := range ringConns {
		roads = append(roads, plainRoad(mainObjectEndpoint("0"), connectionEndpoint(rc)))
	}
	return roads
}

func buildConnectorZoneRoads(connectionNames []string, generateRoads bool) []template.Road {
	if !generateRoads {
		return nil
	}
	var distinct []string
	seen := map[string]bool{}
	for _, n := range connectionNames {
		if n != "" && !seen[n] {
			seen[n] = true
			distinct = append(distinct, n)
		}
	}
	if len(distinct) == 1 {
		return []template.Road{plainRoad(connectionEndpoint(distinct[0]), connectionEndpoint(distinct[0]))}
	}
	if len(distinct) == 0 {
		return nil
	}
	var roads []template.Road
	for _, n := range distinct[1:] {
		roads = append(roads, plainRoad(connectionEndpoint(distinct[0]), connectionEndpoint(n)))
	}
	return roads
}

func buildSideContentLimits() template.StringList {
	var limits []string
	for a := 1; a <= 5; a++ {
		for b := a + 1; b <= 6; b++ {
			limits = append(limits, fmt.Sprintf("content_limits_side_%d_%d", a, b))
		}
	}
	return limits
}

func makeVariant(playerLetters []string, firstLetter string, totalZones int, zones []template.Zone, connections []template.Connection) template.Variant {
	zeroZone := "Neutral-" + firstLetter
	if contains(playerLetters, firstLetter) {
		zeroZone = "Spawn-" + firstLetter
	}
	step := 0
	if totalZones > 0 {
		step = 360 / totalZones
	}
	return template.Variant{
		Orientation: template.Orientation{
			ZeroAngleZone: zeroZone, BaseAngleMin: 45, BaseAngleMax: 45,
			RandomAngleAmplitude: 360, RandomAngleStep: step,
		},
		Border: template.Border{
			CornerRadius: 0, ObstaclesWidth: 3,
			ObstaclesNoise: []template.Noise{{Amplitude: 1, Frequency: 12}},
			WaterWidth:     0, WaterNoise: []template.Noise{{Amplitude: 1, Frequency: 12}},
			WaterType: "water grass",
		},
		Zones: zones, Connections: connections,
	}
}

func buildZoneLayouts() []template.ZoneLayoutDef {
	return []template.ZoneLayoutDef{
		buildZoneLayout(spawnLayoutName, 0.24, 0.48, 0.30, 16, 0.16, 160, -0.30, 0.4, []int{20, 2, 1}),
		buildZoneLayout(sideLayoutName, 0.36, 0.50, 0.25, 16, 0.128, 128, -0.30, 0.3, []int{20, 2, 1}),
		buildZoneLayout(treasureLayoutName, 0.50, 0.50, 0.45, 12, 0.12, 96, -0.30, 0.3, []int{12, 3, 1}),
		buildZoneLayout(centerLayoutName, 0.56, 0.60, 0.30, 10, 0.128, 96, -0.25, 0.3, []int{12, 4, 1}),
	}
}

func buildZoneLayout(name string, obsFill, obsFillVoid, lakesFill float64, minLake int, elevScale float64, roadCluster int, roadAttraction, ambientNoise float64, groupWeights []int) template.ZoneLayoutDef {
	return template.ZoneLayoutDef{
		Name: name, ObstaclesFill: obsFill, ObstaclesFillVoid: obsFillVoid,
		LakesFill: lakesFill, MinLakeArea: minLake, ElevationClusterScale: elevScale,
		ElevationModes: []template.ElevationMode{
			{Weight: 2, MinElevatedFraction: 0.2, MaxElevatedFraction: 0.4},
			{Weight: 1, MinElevatedFraction: 0.6, MaxElevatedFraction: 0.8},
		},
		RoadClusterArea:                   roadCluster,
		GuardedEncounterResourceFractions: template.GuardedEncounterResourceFractions{CountBounds: []int{}, Fractions: []float64{0.66}},
		AmbientPickupDistribution: template.AmbientPickupDistribution{
			Repulsion: 1.0, Noise: ambientNoise, RoadAttraction: roadAttraction,
			ObstacleAttraction: 0, GroupSizeWeights: groupWeights,
		},
	}
}

func buildAllMandatoryContent(playerLetters []string, neutralZones []neutralZonePlan, settings *config.GeneratorConfig) []template.MandatoryContent {
	var groups []template.MandatoryContent
	for _, letter := range playerLetters {
		groups = append(groups, template.MandatoryContent{
			Name:    "mandatory_content_side_" + letter,
			Content: BuildPlayerZoneMandatoryContent(settings),
		})
	}
	for _, nz := range neutralZones {
		var content []template.MandatoryContentItem
		switch nz.Quality {
		case constants.QualityLow:
			content = BuildLowNeutralMandatoryContent(settings, nz.CastleCount)
		case constants.QualityHigh:
			content = BuildHighNeutralMandatoryContent(settings, nz.CastleCount)
		default:
			content = BuildMediumNeutralMandatoryContent(settings, nz.CastleCount)
		}
		groups = append(groups, template.MandatoryContent{
			Name:    "mandatory_content_neutral_" + nz.Letter,
			Content: content,
		})
	}
	return groups
}

func buildRingConnections(playerLetters, orderedLetters []string, tuning generationTuning, isolate bool, neutralByLetter map[string]neutralZonePlan) []template.Connection {
	count := len(orderedLetters)
	if count < 2 {
		return nil
	}
	var conns []template.Connection
	for i := 0; i < count; i++ {
		next := (i + 1) % count
		from := orderedLetters[i]
		to := orderedLetters[next]
		if isolate && contains(playerLetters, from) && contains(playerLetters, to) {
			continue
		}
		fromZone := zoneName(from, playerLetters)
		toZone := zoneName(to, playerLetters)
		conns = append(conns, template.Connection{
			Name: fmt.Sprintf("Ring-%s-%s", from, to), From: fromZone, To: toZone,
			ConnectionType: "Direct", GuardZone: fromZone, SimTurnSquad: true,
			GuardValue: borderGuardValue(from, to, playerLetters, neutralByLetter, tuning), GuardWeeklyIncrement: 0.15,
			GuardMatchGroup: fmt.Sprintf("ring_guard_%s_%s", from, to),
		})
	}
	return conns
}

func buildRandomPortalConnections(playerLetters, orderedLetters []string, tuning generationTuning, maxCount int) []template.Connection {
	count := len(orderedLetters)
	if count < 2 {
		return nil
	}
	dest := buildNonAdjacentDerangement(count)
	indices := make([]int, count)
	for i := range indices {
		indices[i] = i
	}
	rand.Shuffle(len(indices), func(i, j int) { indices[i], indices[j] = indices[j], indices[i] })

	limit := min(count, maxCount)
	trueVal := true
	rule := template.PlacementRule{Type: "Crossroads", Args: []any{}, TargetMin: distNear.Min, TargetMax: distNear.Max, Weight: 2}
	var conns []template.Connection
	for i := range limit {
		idx := indices[i]
		from := orderedLetters[idx]
		to := orderedLetters[dest[idx]]
		fromZone := zoneName(from, playerLetters)
		toZone := zoneName(to, playerLetters)
		conns = append(conns, template.Connection{
			Name: fmt.Sprintf("Portal-%s-%s", from, to), From: fromZone, To: toZone,
			ConnectionType:           "Portal",
			PortalPlacementRulesFrom: []template.PlacementRule{rule},
			PortalPlacementRulesTo:   []template.PlacementRule{rule},
			Road:                     &trueVal, GuardValue: scaleBorderGuardValue(25000, tuning), GuardWeeklyIncrement: 0.15,
		})
	}
	return conns
}

func buildNonAdjacentDerangement(count int) []int {
	dest := make([]int, count)
	for range 100 {
		candidates := make([]int, count)
		for i := range candidates {
			candidates[i] = i
		}
		rand.Shuffle(len(candidates), func(i, j int) { candidates[i], candidates[j] = candidates[j], candidates[i] })
		valid := true
		used := make([]bool, count)
		for i := range count {
			found := -1
			for j := range candidates {
				if used[candidates[j]] {
					continue
				}
				c := candidates[j]
				if c != i && c != (i+1)%count && c != (i-1+count)%count {
					found = j
					break
				}
			}
			if found < 0 {
				for j := range candidates {
					if !used[candidates[j]] && candidates[j] != i {
						found = j
						break
					}
				}
			}
			if found < 0 {
				valid = false
				break
			}
			dest[i] = candidates[found]
			used[candidates[found]] = true
		}
		if valid {
			return dest
		}
	}
	shift := max(1, count/2)
	for i := range count {
		dest[i] = (i + shift) % count
	}
	return dest
}

func buildVariantBalanced(settings *config.GeneratorConfig, playerLetters []string, neutralZones []neutralZonePlan, tuning generationTuning, holdCityNeutralLetter string) template.Variant {
	return buildVariantRandom(settings, playerLetters, neutralZones, tuning, holdCityNeutralLetter)
}

func buildVariantDefault(settings *config.GeneratorConfig, playerLetters []string, neutralZones []neutralZonePlan, tuning generationTuning, holdCityNeutralLetter string) template.Variant {
	neutralByLetter := mapNeutralByLetter(neutralZones)
	ordered := buildOrderedLetters(settings, playerLetters, neutralZones, true)
	n := len(ordered)
	isolate := settings.NoDirectPlayerConnections && len(playerLetters) > 1

	ringConnRight := make([]string, n)
	ringConnLeft := make([]string, n)
	for i := range n {
		next := (i + 1) % n
		if isolate && contains(playerLetters, ordered[i]) && contains(playerLetters, ordered[next]) {
			continue
		}
		name := fmt.Sprintf("Ring-%s-%s", ordered[i], ordered[next])
		ringConnRight[i] = name
		ringConnLeft[next] = name
	}

	var zones []template.Zone
	for i := range n {
		letter := ordered[i]
		var myConns []string
		if ringConnLeft[i] != "" {
			myConns = append(myConns, ringConnLeft[i])
		}
		if ringConnRight[i] != "" && ringConnRight[i] != ringConnLeft[i] {
			myConns = append(myConns, ringConnRight[i])
		}
		if pi := indexOf(playerLetters, letter); pi >= 0 {
			zones = append(zones, buildSpawnZone(letter, fmt.Sprintf("Player%d", pi+1), myConns, settings.ZoneConfiguration.PlayerZoneCastles, settings.MatchPlayerCastleFactions, settings.ZoneConfiguration.Advanced.PlayerZoneSize, settings.SpawnRemoteFootholds, settings.GenerateRoads, tuning))
		} else {
			zones = append(zones, buildNeutralZone(neutralByLetter[letter], myConns, settings.ZoneConfiguration.Advanced.NeutralZoneSize, settings.SpawnRemoteFootholds, settings.GenerateRoads, tuning, letter == holdCityNeutralLetter))
		}
	}

	conns := buildRingConnections(playerLetters, ordered, tuning, isolate, neutralByLetter)
	if settings.RandomPortals {
		conns = append(conns, buildRandomPortalConnections(playerLetters, ordered, tuning, settings.MaxPortalConnections)...)
	}
	if isolate {
		ensurePlayerZonesConnected(playerLetters, zones, &conns, tuning)
	}
	return makeVariant(playerLetters, ordered[0], n, zones, conns)
}

// ── topology: Chain ──────────────────────────────────────────────────

func buildVariantChain(settings *config.GeneratorConfig, playerLetters []string, neutralZones []neutralZonePlan, tuning generationTuning, holdCityNeutralLetter string) template.Variant {
	neutralByLetter := mapNeutralByLetter(neutralZones)
	ordered := buildOrderedLetters(settings, playerLetters, neutralZones, false)
	n := len(ordered)
	isolate := settings.NoDirectPlayerConnections && len(playerLetters) > 1

	connNames := make([]string, n-1)
	for i := 0; i < n-1; i++ {
		if isolate && contains(playerLetters, ordered[i]) && contains(playerLetters, ordered[i+1]) {
			continue
		}
		connNames[i] = fmt.Sprintf("Chain-%s-%s", ordered[i], ordered[i+1])
	}

	var zones []template.Zone
	for i := range n {
		letter := ordered[i]
		var myConns []string
		if i > 0 && connNames[i-1] != "" {
			myConns = append(myConns, connNames[i-1])
		}
		if i < n-1 && connNames[i] != "" {
			myConns = append(myConns, connNames[i])
		}
		if pi := indexOf(playerLetters, letter); pi >= 0 {
			zones = append(zones, buildSpawnZone(letter, fmt.Sprintf("Player%d", pi+1), myConns, settings.ZoneConfiguration.PlayerZoneCastles, settings.MatchPlayerCastleFactions, settings.ZoneConfiguration.Advanced.PlayerZoneSize, settings.SpawnRemoteFootholds, settings.GenerateRoads, tuning))
		} else {
			zones = append(zones, buildNeutralZone(neutralByLetter[letter], myConns, settings.ZoneConfiguration.Advanced.NeutralZoneSize, settings.SpawnRemoteFootholds, settings.GenerateRoads, tuning, letter == holdCityNeutralLetter))
		}
	}

	var conns []template.Connection
	for i := 0; i < n-1; i++ {
		if connNames[i] == "" {
			continue
		}
		from := ordered[i]
		to := ordered[i+1]
		fromZone := zoneName(from, playerLetters)
		toZone := zoneName(to, playerLetters)
		conns = append(conns, template.Connection{
			Name: connNames[i], From: fromZone, To: toZone,
			ConnectionType: "Direct", GuardZone: fromZone, SimTurnSquad: true,
			GuardValue: borderGuardValue(from, to, playerLetters, neutralByLetter, tuning), GuardWeeklyIncrement: 0.15,
			GuardMatchGroup: fmt.Sprintf("chain_guard_%s_%s", from, to),
		})
	}
	if settings.RandomPortals {
		conns = append(conns, buildRandomPortalConnections(playerLetters, ordered, tuning, settings.MaxPortalConnections)...)
	}
	if isolate {
		ensurePlayerZonesConnected(playerLetters, zones, &conns, tuning)
	}
	return makeVariant(playerLetters, ordered[0], n, zones, conns)
}

// ── topology: Hub & Spoke ────────────────────────────────────────────

func buildVariantHubAndSpoke(settings *config.GeneratorConfig, playerLetters []string, neutralZones []neutralZonePlan, tuning generationTuning, hubIsHoldCity bool) template.Variant {
	neutralByLetter := mapNeutralByLetter(neutralZones)
	neutralLetters := make([]string, len(neutralZones))
	for i, nz := range neutralZones {
		neutralLetters[i] = nz.Letter
	}

	var outerLetters []string
	if settings.Topology == config.TopologyBalanced {
		sep := 0
		if settings.MinNeutralZonesBetweenPlayers > 0 && canHonorNeutralSeparation(settings, len(neutralZones)) {
			sep = settings.MinNeutralZonesBetweenPlayers
		}
		outerLetters = buildBalancedRingLetters(playerLetters, neutralZones, sep)
	} else {
		outerLetters = append(append([]string{}, playerLetters...), neutralLetters...)
	}

	var zones []template.Zone
	var conns []template.Connection

	hubConns := make([]string, len(outerLetters))
	for i, l := range outerLetters {
		hubConns[i] = "Hub-" + l
	}
	zones = append(zones, buildHubZone(hubConns, tuning, hubIsHoldCity, settings.ZoneConfiguration.HubZoneSize, settings.ZoneConfiguration.HubZoneCastles, settings.GenerateRoads))

	for i, letter := range outerLetters {
		spokeConns := []string{"Hub-" + letter}
		if pi := indexOf(playerLetters, letter); pi >= 0 {
			zones = append(zones, buildSpawnZone(letter, fmt.Sprintf("Player%d", pi+1), spokeConns, settings.ZoneConfiguration.PlayerZoneCastles, settings.MatchPlayerCastleFactions, settings.ZoneConfiguration.Advanced.PlayerZoneSize, settings.SpawnRemoteFootholds, settings.GenerateRoads, tuning))
		} else {
			zones = append(zones, buildNeutralZone(neutralByLetter[letter], spokeConns, settings.ZoneConfiguration.Advanced.NeutralZoneSize, settings.SpawnRemoteFootholds, settings.GenerateRoads, tuning, false))
		}
		_ = i
	}

	for _, letter := range outerLetters {
		outerZone := zoneName(letter, playerLetters)
		hubAnchor := letter
		if len(playerLetters) > 0 {
			hubAnchor = playerLetters[0]
		}
		hubGuard := borderGuardValue(hubAnchor, letter, playerLetters, neutralByLetter, tuning)
		conns = append(conns,
			template.Connection{
				Name: "Hub-" + letter, From: "Hub", To: outerZone,
				ConnectionType: "Direct", GuardZone: "Hub", SimTurnSquad: true,
				GuardValue: hubGuard, GuardWeeklyIncrement: 0.15,
				GuardMatchGroup: "hub_guard_" + letter,
			},
			template.Connection{
				From: "Hub", To: outerZone, ConnectionType: "Direct",
				GuardZone: "Hub", SimTurnSquad: true,
				GuardValue: hubGuard, GuardWeeklyIncrement: 0.15,
				GuardMatchGroup: fmt.Sprintf("hub_guard_%s_%d", letter, 1),
			})
	}

	// Proximity ring
	for i := 0; i < len(outerLetters); i++ {
		next := (i + 1) % len(outerLetters)
		from := outerLetters[i]
		to := outerLetters[next]
		if settings.NoDirectPlayerConnections && contains(playerLetters, from) && contains(playerLetters, to) {
			continue
		}
		conns = append(conns, template.Connection{
			Name: fmt.Sprintf("Pseudo-%s-%s", from, to),
			From: zoneName(from, playerLetters), To: zoneName(to, playerLetters),
			ConnectionType: "Proximity",
		})
	}

	if settings.RandomPortals {
		conns = append(conns, buildRandomPortalConnections(playerLetters, outerLetters, tuning, settings.MaxPortalConnections)...)
	}
	return makeVariant(playerLetters, outerLetters[0], len(outerLetters)+1, zones, conns)
}

// ── topology: Shared Web ─────────────────────────────────────────────

func buildVariantSharedWeb(settings *config.GeneratorConfig, playerLetters []string, neutralZones []neutralZonePlan, tuning generationTuning, holdCityNeutralLetter string) template.Variant {
	neutralByLetter := mapNeutralByLetter(neutralZones)
	var neutrals []string
	if settings.Topology == config.TopologyBalanced {
		neutrals = buildBalancedNeutralRing(neutralZones, len(playerLetters))
	} else {
		for _, nz := range neutralZones {
			neutrals = append(neutrals, nz.Letter)
		}
	}

	p := len(playerLetters)
	nn := len(neutrals)

	spokeByPlayer := map[string][]string{}
	spokeByNeutral := map[string][]string{}
	for _, l := range playerLetters {
		spokeByPlayer[l] = nil
	}
	for _, l := range neutrals {
		spokeByNeutral[l] = nil
	}
	addSpoke := func(pl, nl string) {
		cn := fmt.Sprintf("Web-%s-%s", pl, nl)
		spokeByPlayer[pl] = append(spokeByPlayer[pl], cn)
		spokeByNeutral[nl] = append(spokeByNeutral[nl], cn)
	}
	for i := 0; i < p; i++ {
		n1 := (i * nn / p) % nn
		n2 := ((i * nn / p) + 1) % nn
		addSpoke(playerLetters[i], neutrals[n1])
		if n1 != n2 {
			addSpoke(playerLetters[i], neutrals[n2])
		}
	}

	var zones []template.Zone
	var conns []template.Connection

	neutralRingConns := make([]string, nn)
	for i := 0; i < nn; i++ {
		next := (i + 1) % nn
		neutralRingConns[i] = fmt.Sprintf("NRing-%s-%s", neutrals[i], neutrals[next])
	}

	for i := 0; i < nn; i++ {
		prev := (i - 1 + nn) % nn
		var nConns []string
		if nn > 1 {
			nConns = append(nConns, neutralRingConns[prev], neutralRingConns[i])
		}
		nConns = append(nConns, spokeByNeutral[neutrals[i]]...)
		nConns = uniqueStrings(nConns)
		z := buildNeutralZone(neutralByLetter[neutrals[i]], nConns, settings.ZoneConfiguration.Advanced.NeutralZoneSize, settings.SpawnRemoteFootholds, settings.GenerateRoads, tuning, neutrals[i] == holdCityNeutralLetter)
		if neutralByLetter[neutrals[i]].CastleCount == 0 {
			z.Roads = buildConnectorZoneRoads(nConns, settings.GenerateRoads)
		}
		zones = append(zones, z)
	}

	for i := 0; i < p; i++ {
		pl := playerLetters[i]
		sc := spokeByPlayer[pl]
		zones = append(zones, buildSpawnZone(pl, fmt.Sprintf("Player%d", i+1), sc, settings.ZoneConfiguration.PlayerZoneCastles, settings.MatchPlayerCastleFactions, settings.ZoneConfiguration.Advanced.PlayerZoneSize, settings.SpawnRemoteFootholds, settings.GenerateRoads, tuning))

		for _, cn := range sc {
			parts := strings.Split(cn, "-")
			nl := parts[2]
			conns = append(conns, template.Connection{
				Name: cn, From: "Spawn-" + pl, To: "Neutral-" + nl,
				ConnectionType: "Direct", GuardZone: "Neutral-" + nl, SimTurnSquad: true,
				GuardValue: borderGuardValue(pl, nl, playerLetters, neutralByLetter, tuning), GuardWeeklyIncrement: 0.15,
				GuardMatchGroup: fmt.Sprintf("web_guard_%s_%s", pl, nl),
			})
		}
	}

	if nn > 1 {
		for i := 0; i < nn; i++ {
			next := (i + 1) % nn
			conns = append(conns, template.Connection{
				Name: neutralRingConns[i], From: "Neutral-" + neutrals[i], To: "Neutral-" + neutrals[next],
				ConnectionType: "Direct", GuardZone: "Neutral-" + neutrals[i], SimTurnSquad: true,
				GuardValue: borderGuardValue(neutrals[i], neutrals[next], playerLetters, neutralByLetter, tuning), GuardWeeklyIncrement: 0.15,
				GuardMatchGroup: fmt.Sprintf("nring_guard_%s_%s", neutrals[i], neutrals[next]),
			})
		}
	}

	if settings.RandomPortals {
		all := append(append([]string{}, playerLetters...), neutrals...)
		conns = append(conns, buildRandomPortalConnections(playerLetters, all, tuning, settings.MaxPortalConnections)...)
	}
	if settings.NoDirectPlayerConnections && len(playerLetters) > 1 {
		ensurePlayerZonesConnected(playerLetters, zones, &conns, tuning)
	}
	return makeVariant(playerLetters, playerLetters[0], len(zones), zones, conns)
}

// ── topology: Random ─────────────────────────────────────────────────

func buildVariantRandom(settings *config.GeneratorConfig, playerLetters []string, neutralZones []neutralZonePlan, tuning generationTuning, holdCityNeutralLetter string) template.Variant {
	neutralByLetter := mapNeutralByLetter(neutralZones)
	neutralLetters := make([]string, len(neutralZones))
	for i, nz := range neutralZones {
		neutralLetters[i] = nz.Letter
	}
	isolate := settings.NoDirectPlayerConnections && len(playerLetters) > 1

	var allLetters []string
	if settings.Topology == config.TopologyBalanced {
		allLetters = buildBalancedRingLetters(playerLetters, neutralZones, 0)
	} else {
		allLetters = append(append([]string{}, playerLetters...), neutralLetters...)
		rand.Shuffle(len(allLetters), func(i, j int) { allLetters[i], allLetters[j] = allLetters[j], allLetters[i] })
	}
	count := len(allLetters)

	var pos [][2]float64
	if settings.Topology == config.TopologyBalanced {
		pos = buildBalancedRandomPositions(allLetters, playerLetters, neutralByLetter)
	} else {
		for i := 0; i < count; i++ {
			pos = append(pos, [2]float64{rand.Float64()*0.9 + 0.05, rand.Float64()*0.9 + 0.05})
		}
	}

	pairs := delaunayEdges(pos)

	if settings.Topology == config.TopologyBalanced {
		presentTiers := map[int]bool{}
		for _, l := range allLetters {
			presentTiers[zoneTierRank(l, playerLetters, neutralByLetter)] = true
		}
		var filtered [][2]int
		for _, p := range pairs {
			ta := zoneTierRank(allLetters[p[0]], playerLetters, neutralByLetter)
			tb := zoneTierRank(allLetters[p[1]], playerLetters, neutralByLetter)
			lo, hi := ta, tb
			if lo > hi {
				lo, hi = hi, lo
			}
			if hi-lo <= 1 {
				filtered = append(filtered, p)
				continue
			}
			skip := false
			for t := lo + 1; t < hi; t++ {
				if presentTiers[t] {
					skip = true
					break
				}
			}
			if !skip {
				filtered = append(filtered, p)
			}
		}
		pairs = filtered
	}

	connsByZone := make(map[int][]string, count)
	var conns []template.Connection
	for _, p := range pairs {
		a, b := p[0], p[1]
		from := allLetters[a]
		to := allLetters[b]
		if isolate && contains(playerLetters, from) && contains(playerLetters, to) {
			continue
		}
		cn := fmt.Sprintf("Rnd-%s-%s", from, to)
		connsByZone[a] = append(connsByZone[a], cn)
		connsByZone[b] = append(connsByZone[b], cn)
		conns = append(conns, template.Connection{
			Name: cn, From: zoneName(from, playerLetters), To: zoneName(to, playerLetters),
			ConnectionType: "Direct", GuardZone: zoneName(from, playerLetters), SimTurnSquad: true,
			GuardValue: borderGuardValue(from, to, playerLetters, neutralByLetter, tuning), GuardWeeklyIncrement: 0.15,
			GuardMatchGroup: fmt.Sprintf("rnd_guard_%s_%s", from, to),
		})
	}

	var zones []template.Zone
	for i := 0; i < count; i++ {
		letter := allLetters[i]
		myConns := connsByZone[i]
		if pi := indexOf(playerLetters, letter); pi >= 0 {
			zones = append(zones, buildSpawnZone(letter, fmt.Sprintf("Player%d", pi+1), myConns, settings.ZoneConfiguration.PlayerZoneCastles, settings.MatchPlayerCastleFactions, settings.ZoneConfiguration.Advanced.PlayerZoneSize, settings.SpawnRemoteFootholds, settings.GenerateRoads, tuning))
		} else {
			zones = append(zones, buildNeutralZone(neutralByLetter[letter], myConns, settings.ZoneConfiguration.Advanced.NeutralZoneSize, settings.SpawnRemoteFootholds, settings.GenerateRoads, tuning, letter == holdCityNeutralLetter))
		}
	}

	// Stamp generator-driven positions onto the freshly built zones so the
	// preview renderer can reproduce the exact geometry used to derive the
	// Delaunay connections. Balanced layouts also stamp the concentric ring
	// index so the preview can snap zones to clean rings
	for i := range zones {
		p := pos[i]
		zones[i].GeneratorPosition = &[2]float64{p[0], p[1]}
		if settings.Topology == config.TopologyBalanced {
			r := zoneTierRank(allLetters[i], playerLetters, neutralByLetter)
			zones[i].GeneratorRing = &r
		}
	}

	if settings.RandomPortals {
		conns = append(conns, buildRandomPortalConnections(playerLetters, allLetters, tuning, settings.MaxPortalConnections)...)
	}
	if isolate {
		ensurePlayerZonesConnected(playerLetters, zones, &conns, tuning)
	}
	ensureFullConnectivity(playerLetters, allLetters, pos, zones, &conns, tuning, neutralByLetter)
	return makeVariant(playerLetters, allLetters[0], count, zones, conns)
}

// ── topology: Tournament ─────────────────────────────────────────────

func buildVariantTournament(settings *config.GeneratorConfig, playerLetters []string, neutralZones []neutralZonePlan, tuning generationTuning) template.Variant {
	neutralByLetter := mapNeutralByLetter(neutralZones)

	// Distribute neutrals in a balanced round-robin: sort by descending quality
	// (then castle count, then letter) so that quality tiers are split evenly
	// across the two players (e91e79f / v0.7 ordering).
	sorted := make([]neutralZonePlan, len(neutralZones))
	copy(sorted, neutralZones)
	sort.SliceStable(sorted, func(i, j int) bool {
		if sorted[i].Quality != sorted[j].Quality {
			return sorted[i].Quality > sorted[j].Quality
		}
		if sorted[i].CastleCount != sorted[j].CastleCount {
			return sorted[i].CastleCount > sorted[j].CastleCount
		}
		return sorted[i].Letter < sorted[j].Letter
	})
	neutralsForPlayer := [2][]neutralZonePlan{}
	for i, nz := range sorted {
		neutralsForPlayer[i%2] = append(neutralsForPlayer[i%2], nz)
	}

	for p := range 2 {
		sort.SliceStable(neutralsForPlayer[p], func(i, j int) bool {
			ai, aj := neutralsForPlayer[p][i], neutralsForPlayer[p][j]
			si, sj := neutralZoneBalanceScore(ai), neutralZoneBalanceScore(aj)
			if si != sj {
				return si < sj
			}
			return ai.Letter < aj.Letter
		})
	}

	var zones []template.Zone
	var conns []template.Connection

	switch settings.Topology {
	case config.TopologyHubAndSpoke:
		for p := range 2 {
			buildTournamentHubCluster(p, playerLetters[p], neutralsForPlayer[p], neutralByLetter, settings, tuning, &zones, &conns)
		}
	case config.TopologyBalanced:
		for p := range 2 {
			buildTournamentBalancedCluster(p, playerLetters[p], neutralsForPlayer[p], neutralByLetter, settings, tuning, &zones, &conns)
		}
	case config.TopologyDefault:
		for p := range 2 {
			buildTournamentRingCluster(p, playerLetters[p], neutralsForPlayer[p], neutralByLetter, settings, tuning, &zones, &conns)
		}
	default:
		// Chain, SharedWeb, Random → chain-per-cluster fallback.
		for p := range 2 {
			buildTournamentChainCluster(p, playerLetters[p], neutralsForPlayer[p], neutralByLetter, settings, tuning, &zones, &conns)
		}
	}

	// c20b40d: per-cluster portals so they never cross the isolation boundary.
	if settings.RandomPortals {
		for p := range 2 {
			clusterLetters := []string{playerLetters[p]}
			for _, n := range neutralsForPlayer[p] {
				clusterLetters = append(clusterLetters, n.Letter)
			}
			conns = append(conns, buildRandomPortalConnections(playerLetters, clusterLetters, tuning, settings.MaxPortalConnections)...)
		}
	}

	return makeVariant(playerLetters, playerLetters[0], len(zones), zones, conns)
}

func buildTournamentChainCluster(playerIndex int, playerLetter string, myNeutrals []neutralZonePlan, neutralByLetter map[string]neutralZonePlan, settings *config.GeneratorConfig, tuning generationTuning, zones *[]template.Zone, connections *[]template.Connection) {
	chain := []string{playerLetter}
	for _, n := range myNeutrals {
		chain = append(chain, n.Letter)
	}
	connNames := make([]string, len(chain)-1)
	for i := range connNames {
		connNames[i] = fmt.Sprintf("Tourney-%s-%s", chain[i], chain[i+1])
	}
	for i, letter := range chain {
		var myConns []string
		if i > 0 {
			myConns = append(myConns, connNames[i-1])
		}
		if i < len(connNames) {
			myConns = append(myConns, connNames[i])
		}
		if i == 0 {
			*zones = append(*zones, buildSpawnZone(letter, fmt.Sprintf("Player%d", playerIndex+1), myConns, settings.ZoneConfiguration.PlayerZoneCastles, settings.MatchPlayerCastleFactions, settings.ZoneConfiguration.Advanced.PlayerZoneSize, settings.SpawnRemoteFootholds, settings.GenerateRoads, tuning))
		} else {
			*zones = append(*zones, buildNeutralZone(neutralByLetter[letter], myConns, settings.ZoneConfiguration.Advanced.NeutralZoneSize, settings.SpawnRemoteFootholds, settings.GenerateRoads, tuning, false))
		}
	}
	for i := range connNames {
		from := chain[i]
		to := chain[i+1]
		fromZone := "Spawn-" + from
		if i > 0 {
			fromZone = "Neutral-" + from
		}
		*connections = append(*connections, template.Connection{
			Name: connNames[i], From: fromZone, To: "Neutral-" + to,
			ConnectionType: "Direct", GuardZone: fromZone, SimTurnSquad: true,
			GuardValue: borderGuardValue(from, to, []string{playerLetter}, neutralByLetter, tuning), GuardWeeklyIncrement: 0.15,
			GuardMatchGroup: fmt.Sprintf("tourney_guard_%s_%s", from, to),
		})
	}
}

// buildTournamentRingCluster — one player's isolated cluster as a ring:
// player → low → … → high → … → low → player. Sorts neutrals by balance
// score, then fills outward-from-player so the highest-quality zones sit
// at the ring midpoint
func buildTournamentRingCluster(playerIndex int, playerLetter string, myNeutrals []neutralZonePlan, neutralByLetter map[string]neutralZonePlan, settings *config.GeneratorConfig, tuning generationTuning, zones *[]template.Zone, connections *[]template.Connection) {
	sortedNeutrals := make([]neutralZonePlan, len(myNeutrals))
	copy(sortedNeutrals, myNeutrals)
	sort.SliceStable(sortedNeutrals, func(i, j int) bool {
		si, sj := neutralZoneBalanceScore(sortedNeutrals[i]), neutralZoneBalanceScore(sortedNeutrals[j])
		if si != sj {
			return si < sj
		}
		return sortedNeutrals[i].Letter < sortedNeutrals[j].Letter
	})

	n := len(sortedNeutrals)
	orderedNeutrals := make([]neutralZonePlan, n)
	lo, hi := 0, n-1
	for i := range n {
		if i%2 == 0 {
			orderedNeutrals[lo] = sortedNeutrals[i]
			lo++
		} else {
			orderedNeutrals[hi] = sortedNeutrals[i]
			hi--
		}
	}

	ringLetters := []string{playerLetter}
	for _, nz := range orderedNeutrals {
		ringLetters = append(ringLetters, nz.Letter)
	}
	count := len(ringLetters)
	if count < 2 {
		// Lone player zone — no ring edges possible.
		*zones = append(*zones, buildSpawnZone(playerLetter, fmt.Sprintf("Player%d", playerIndex+1), nil, settings.ZoneConfiguration.PlayerZoneCastles, settings.MatchPlayerCastleFactions, settings.ZoneConfiguration.Advanced.PlayerZoneSize, settings.SpawnRemoteFootholds, settings.GenerateRoads, tuning))
		return
	}

	connNames := make([]string, count)
	for i := range count {
		next := (i + 1) % count
		connNames[i] = fmt.Sprintf("TRing-%s-%s", ringLetters[i], ringLetters[next])
	}

	for i, letter := range ringLetters {
		prev := (i - 1 + count) % count
		seen := map[string]bool{}
		var myConns []string
		for _, name := range []string{connNames[prev], connNames[i]} {
			if !seen[name] {
				seen[name] = true
				myConns = append(myConns, name)
			}
		}
		if i == 0 {
			*zones = append(*zones, buildSpawnZone(letter, fmt.Sprintf("Player%d", playerIndex+1), myConns, settings.ZoneConfiguration.PlayerZoneCastles, settings.MatchPlayerCastleFactions, settings.ZoneConfiguration.Advanced.PlayerZoneSize, settings.SpawnRemoteFootholds, settings.GenerateRoads, tuning))
		} else {
			*zones = append(*zones, buildNeutralZone(neutralByLetter[letter], myConns, settings.ZoneConfiguration.Advanced.NeutralZoneSize, settings.SpawnRemoteFootholds, settings.GenerateRoads, tuning, false))
		}
	}

	for i := range count {
		next := (i + 1) % count
		from := ringLetters[i]
		to := ringLetters[next]
		fromZone := "Spawn-" + from
		if i != 0 {
			fromZone = "Neutral-" + from
		}
		toZone := "Spawn-" + to
		if next != 0 {
			toZone = "Neutral-" + to
		}
		*connections = append(*connections, template.Connection{
			Name: connNames[i], From: fromZone, To: toZone,
			ConnectionType: "Direct", GuardZone: fromZone, SimTurnSquad: true,
			GuardValue: borderGuardValue(from, to, []string{playerLetter}, neutralByLetter, tuning), GuardWeeklyIncrement: 0.15,
			GuardMatchGroup: fmt.Sprintf("tourney_ring_guard_%s_%s", from, to),
		})
	}
}

// buildTournamentHubCluster — one player's isolated cluster as a private
// hub-and-spoke layout. A dedicated mini-hub zone "Hub-{playerLetter}" sits
// at the centre and connects directly to the player spawn and all of their
// exclusive neutrals
func buildTournamentHubCluster(playerIndex int, playerLetter string, myNeutrals []neutralZonePlan, neutralByLetter map[string]neutralZonePlan, settings *config.GeneratorConfig, tuning generationTuning, zones *[]template.Zone, connections *[]template.Connection) {
	hubName := "Hub-" + playerLetter

	spokeLetters := []string{playerLetter}
	for _, nz := range myNeutrals {
		spokeLetters = append(spokeLetters, nz.Letter)
	}

	spokeConnNames := make([]string, len(spokeLetters))
	for i, l := range spokeLetters {
		spokeConnNames[i] = fmt.Sprintf("THubSpoke-%s-%s", playerLetter, l)
	}

	hubZone := buildHubZone(spokeConnNames, tuning, false, settings.ZoneConfiguration.HubZoneSize, settings.ZoneConfiguration.HubZoneCastles, settings.GenerateRoads)
	hubZone.Name = hubName
	*zones = append(*zones, hubZone)

	for i, letter := range spokeLetters {
		conn := []string{spokeConnNames[i]}
		if i == 0 {
			*zones = append(*zones, buildSpawnZone(letter, fmt.Sprintf("Player%d", playerIndex+1), conn, settings.ZoneConfiguration.PlayerZoneCastles, settings.MatchPlayerCastleFactions, settings.ZoneConfiguration.Advanced.PlayerZoneSize, settings.SpawnRemoteFootholds, settings.GenerateRoads, tuning))
		} else {
			*zones = append(*zones, buildNeutralZone(neutralByLetter[letter], conn, settings.ZoneConfiguration.Advanced.NeutralZoneSize, settings.SpawnRemoteFootholds, settings.GenerateRoads, tuning, false))
		}
	}

	for i, spokeLetter := range spokeLetters {
		spokeZone := "Spawn-" + spokeLetter
		if i != 0 {
			spokeZone = "Neutral-" + spokeLetter
		}
		*connections = append(*connections, template.Connection{
			Name: spokeConnNames[i], From: hubName, To: spokeZone,
			ConnectionType: "Direct", GuardZone: hubName, SimTurnSquad: true,
			GuardValue: borderGuardValue(playerLetter, spokeLetter, []string{playerLetter}, neutralByLetter, tuning), GuardWeeklyIncrement: 0.15,
			GuardMatchGroup: fmt.Sprintf("tourney_hub_guard_%s_%s", playerLetter, spokeLetter),
		})
	}

	// Proximity ring around spokes so the engine arranges them sensibly.
	for i := 0; i < len(spokeLetters); i++ {
		next := (i + 1) % len(spokeLetters)
		from := spokeLetters[i]
		to := spokeLetters[next]
		fromZone := "Spawn-" + from
		if i != 0 {
			fromZone = "Neutral-" + from
		}
		toZone := "Spawn-" + to
		if next != 0 {
			toZone = "Neutral-" + to
		}
		*connections = append(*connections, template.Connection{
			Name:           fmt.Sprintf("TPseudo-%s-%s-%s", playerLetter, from, to),
			From:           fromZone,
			To:             toZone,
			ConnectionType: "Proximity",
		})
	}
}

func buildTournamentBalancedCluster(playerIndex int, playerLetter string, myNeutrals []neutralZonePlan, neutralByLetter map[string]neutralZonePlan, settings *config.GeneratorConfig, tuning generationTuning, zones *[]template.Zone, connections *[]template.Connection) {
	singlePlayerList := []string{playerLetter}
	orderedLetters := buildBalancedRingLetters(singlePlayerList, myNeutrals, 0)
	rawPos := buildBalancedRandomPositions(orderedLetters, singlePlayerList, neutralByLetter)

	// Position remap to the player's canvas half (kept for parity even
	// though Go's template.Zone does not yet have a GeneratorPosition
	// field; preview will pick this up in Phase 7).
	rawXMin, rawXMax, rawYMin, rawYMax := 0.05, 0.95, 0.05, 0.95
	if len(rawPos) > 0 {
		rawXMin, rawXMax = rawPos[0][0], rawPos[0][0]
		rawYMin, rawYMax = rawPos[0][1], rawPos[0][1]
		for _, p := range rawPos[1:] {
			if p[0] < rawXMin {
				rawXMin = p[0]
			}
			if p[0] > rawXMax {
				rawXMax = p[0]
			}
			if p[1] < rawYMin {
				rawYMin = p[1]
			}
			if p[1] > rawYMax {
				rawYMax = p[1]
			}
		}
	}
	spanX := math.Max(rawXMax-rawXMin, 0.001)
	spanY := math.Max(rawYMax-rawYMin, 0.001)
	xMin, xMax := 0.03, 0.43
	if playerIndex != 0 {
		xMin, xMax = 0.57, 0.97
	}
	targetW := xMax - xMin
	const targetH = 0.90
	scale := math.Min(targetW/spanX, targetH/spanY)
	xCentre := (xMin + xMax) / 2.0
	const yCentre = 0.5
	pos := make([][2]float64, len(rawPos))
	for i, pt := range rawPos {
		pos[i] = [2]float64{
			xCentre + (pt[0]-(rawXMin+rawXMax)/2.0)*scale,
			yCentre + (pt[1]-(rawYMin+rawYMax)/2.0)*scale,
		}
	}

	// Build connections from pure ring structure (a040c98).
	angDist := func(a, b float64) float64 {
		d := math.Mod(math.Abs(a-b), 2*math.Pi)
		if d > math.Pi {
			d = 2*math.Pi - d
		}
		return d
	}
	tierIndices := map[int][]int{}
	for i, l := range orderedLetters {
		t := zoneTierRank(l, singlePlayerList, neutralByLetter)
		tierIndices[t] = append(tierIndices[t], i)
	}
	tierKeys := make([]int, 0, len(tierIndices))
	for k := range tierIndices {
		tierKeys = append(tierKeys, k)
	}
	sort.Ints(tierKeys)

	tierSorted := map[int][]int{}
	tierAngles := map[int][]float64{}
	for tier, idx := range tierIndices {
		s := make([]int, len(idx))
		copy(s, idx)
		sort.SliceStable(s, func(i, j int) bool {
			return math.Atan2(rawPos[s[i]][1]-0.5, rawPos[s[i]][0]-0.5) <
				math.Atan2(rawPos[s[j]][1]-0.5, rawPos[s[j]][0]-0.5)
		})
		tierSorted[tier] = s
		ang := make([]float64, len(s))
		for j, ii := range s {
			ang[j] = math.Atan2(rawPos[ii][1]-0.5, rawPos[ii][0]-0.5)
		}
		tierAngles[tier] = ang
	}

	pairSet := map[[2]int]bool{}
	addPair := func(a, b int) {
		if a > b {
			a, b = b, a
		}
		pairSet[[2]int{a, b}] = true
	}

	// Same-ring: circle-neighbors only; skip degenerate < 3 rings.
	for _, sorted := range tierSorted {
		nn := len(sorted)
		if nn < 3 {
			continue
		}
		for j := 0; j < nn; j++ {
			addPair(sorted[j], sorted[(j+1)%nn])
		}
	}

	// Cross-ring: bidirectional nearest-neighbor between adjacent tiers.
	for ti := 0; ti+1 < len(tierKeys); ti++ {
		outerSorted := tierSorted[tierKeys[ti]]
		innerSorted := tierSorted[tierKeys[ti+1]]
		outerAngles := tierAngles[tierKeys[ti]]
		innerAngles := tierAngles[tierKeys[ti+1]]

		for oi := 0; oi < len(outerSorted); oi++ {
			best, bestD := 0, math.MaxFloat64
			for ii := 0; ii < len(innerSorted); ii++ {
				if d := angDist(outerAngles[oi], innerAngles[ii]); d < bestD {
					bestD = d
					best = ii
				}
			}
			if len(innerSorted) > 0 {
				addPair(outerSorted[oi], innerSorted[best])
			}
		}

		const epsilon = 1e-9
		for ii := 0; ii < len(innerSorted); ii++ {
			bestD := math.MaxFloat64
			for oi := 0; oi < len(outerSorted); oi++ {
				if d := angDist(innerAngles[ii], outerAngles[oi]); d < bestD {
					bestD = d
				}
			}
			for oi := 0; oi < len(outerSorted); oi++ {
				if angDist(innerAngles[ii], outerAngles[oi]) <= bestD+epsilon {
					addPair(innerSorted[ii], outerSorted[oi])
				}
			}
		}
	}

	count := len(orderedLetters)
	connsByZone := make([][]string, count)
	for _, p := range sortedPairs(pairSet) {
		from := orderedLetters[p[0]]
		to := orderedLetters[p[1]]
		connName := fmt.Sprintf("TBal-%s-%s", from, to)
		connsByZone[p[0]] = append(connsByZone[p[0]], connName)
		connsByZone[p[1]] = append(connsByZone[p[1]], connName)

		fromZone := "Spawn-" + from
		if from != playerLetter {
			fromZone = "Neutral-" + from
		}
		toZone := "Spawn-" + to
		if to != playerLetter {
			toZone = "Neutral-" + to
		}
		*connections = append(*connections, template.Connection{
			Name: connName, From: fromZone, To: toZone,
			ConnectionType: "Direct", GuardZone: fromZone, SimTurnSquad: true,
			GuardValue: borderGuardValue(from, to, []string{playerLetter}, neutralByLetter, tuning), GuardWeeklyIncrement: 0.15,
			GuardMatchGroup: fmt.Sprintf("tourney_bal_guard_%s_%s", from, to),
		})
	}

	clusterStart := len(*zones)
	for i, letter := range orderedLetters {
		myConns := connsByZone[i]
		if letter == playerLetter {
			*zones = append(*zones, buildSpawnZone(letter, fmt.Sprintf("Player%d", playerIndex+1), myConns, settings.ZoneConfiguration.PlayerZoneCastles, settings.MatchPlayerCastleFactions, settings.ZoneConfiguration.Advanced.PlayerZoneSize, settings.SpawnRemoteFootholds, settings.GenerateRoads, tuning))
		} else {
			*zones = append(*zones, buildNeutralZone(neutralByLetter[letter], myConns, settings.ZoneConfiguration.Advanced.NeutralZoneSize, settings.SpawnRemoteFootholds, settings.GenerateRoads, tuning, false))
		}
	}

	// Stamp generator positions and ring indices onto the freshly built cluster
	// zones so the preview renderer can reproduce the tournament-balanced
	// geometry without re-deriving it from connections.
	for i := 0; i < len(orderedLetters); i++ {
		p := pos[i]
		(*zones)[clusterStart+i].GeneratorPosition = &[2]float64{p[0], p[1]}
		r := zoneTierRank(orderedLetters[i], singlePlayerList, neutralByLetter)
		(*zones)[clusterStart+i].GeneratorRing = &r
	}

	// Ensure the cluster is fully connected (same guarantee as the standard
	// balanced variant). Operate on the slice header of the cluster's zones.
	clusterZones := (*zones)[clusterStart:]
	ensureFullConnectivity(singlePlayerList, orderedLetters, pos, clusterZones, connections, tuning, neutralByLetter)
}

// sortedPairs returns the keys of a [2]int→bool set in deterministic order
// (by element 0 then element 1) so that connection ordering does not depend
// on Go's map iteration randomness.
func sortedPairs(set map[[2]int]bool) [][2]int {
	out := make([][2]int, 0, len(set))
	for p := range set {
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i][0] != out[j][0] {
			return out[i][0] < out[j][0]
		}
		return out[i][1] < out[j][1]
	})
	return out
}

// ── isolation / connectivity failsafes ───────────────────────────────

func ensurePlayerZonesConnected(playerLetters []string, zones []template.Zone, connections *[]template.Connection, tuning generationTuning) {
	if len(playerLetters) < 2 {
		return
	}
	connNames := map[string]bool{}
	for _, c := range *connections {
		if c.Name != "" {
			connNames[c.Name] = true
		}
	}
	for _, letter := range playerLetters {
		zn := "Spawn-" + letter
		z := findZone(zones, zn)
		if z == nil {
			continue
		}
		hasConn := false
		for _, r := range z.Roads {
			if r.To.Type == "Connection" && len(r.To.Args) > 0 && connNames[r.To.Args[0]] {
				hasConn = true
				break
			}
		}
		if hasConn {
			continue
		}
		var partner string
		for _, pl := range playerLetters {
			if pl != letter {
				partner = pl
				break
			}
		}
		if partner == "" {
			continue
		}
		a, b := letter, partner
		if a > b {
			a, b = b, a
		}
		fn := "Fallback-" + a + "-" + b
		if connNames[fn] {
			continue
		}
		*connections = append(*connections, template.Connection{
			Name: fn, From: "Spawn-" + letter, To: "Spawn-" + partner,
			ConnectionType: "Direct", GuardZone: "Spawn-" + letter, SimTurnSquad: true,
			GuardValue: borderGuardValue(letter, partner, playerLetters, nil, tuning), GuardWeeklyIncrement: 0.15,
			GuardMatchGroup: "fallback_guard_" + fn,
		})
		connNames[fn] = true
		for _, pl := range []string{letter, partner} {
			pz := findZone(zones, "Spawn-"+pl)
			if pz != nil {
				pz.Roads = append(pz.Roads, plainRoad(mainObjectEndpoint("0"), connectionEndpoint(fn)))
			}
		}
	}
}

func ensureFullConnectivity(playerLetters, allLetters []string, pos [][2]float64, zones []template.Zone, connections *[]template.Connection, tuning generationTuning, neutralByLetter map[string]neutralZonePlan) {
	if len(allLetters) <= 1 {
		return
	}
	zoneNameToIdx := map[string]int{}
	for i, l := range allLetters {
		zoneNameToIdx[zoneName(l, playerLetters)] = i
	}
	adj := make(map[int]map[int]bool, len(allLetters))
	for i := range allLetters {
		adj[i] = map[int]bool{}
	}
	for _, c := range *connections {
		if c.ConnectionType != "Direct" && c.ConnectionType != "Portal" {
			continue
		}
		a, okA := zoneNameToIdx[c.From]
		b, okB := zoneNameToIdx[c.To]
		if !okA || !okB {
			continue
		}
		adj[a][b] = true
		adj[b][a] = true
	}

	connNameSet := map[string]bool{}
	for _, c := range *connections {
		if c.Name != "" {
			connNameSet[c.Name] = true
		}
	}

	for {
		components := findComponents(adj, len(allLetters))
		if len(components) <= 1 {
			break
		}
		mainComp := map[int]bool{}
		for _, idx := range components[0] {
			mainComp[idx] = true
		}
		bestA, bestB := -1, -1
		bestDist := math.MaxFloat64
		for _, a := range components[0] {
			for ci := 1; ci < len(components); ci++ {
				for _, b := range components[ci] {
					dx := pos[a][0] - pos[b][0]
					dy := pos[a][1] - pos[b][1]
					d := dx*dx + dy*dy
					if d < bestDist {
						bestDist = d
						bestA, bestB = a, b
					}
				}
			}
		}
		if bestA < 0 {
			break
		}
		la, lb := allLetters[bestA], allLetters[bestB]
		if la > lb {
			la, lb = lb, la
		}
		bridgeName := fmt.Sprintf("Bridge-%s-%s", la, lb)
		if !connNameSet[bridgeName] {
			za := zoneName(allLetters[bestA], playerLetters)
			zb := zoneName(allLetters[bestB], playerLetters)
			*connections = append(*connections, template.Connection{
				Name: bridgeName, From: za, To: zb,
				ConnectionType: "Direct", GuardZone: za, SimTurnSquad: true,
				GuardValue: borderGuardValue(la, lb, playerLetters, neutralByLetter, tuning), GuardWeeklyIncrement: 0.15,
				GuardMatchGroup: fmt.Sprintf("bridge_guard_%s-%s", la, lb),
			})
			connNameSet[bridgeName] = true
			for _, zn := range []string{za, zb} {
				z := findZone(zones, zn)
				if z != nil {
					if len(z.MainObjects) > 0 {
						z.Roads = append(z.Roads, plainRoad(mainObjectEndpoint("0"), connectionEndpoint(bridgeName)))
					} else if len(z.Roads) > 0 {
						existingConn := ""
						for _, r := range z.Roads {
							if r.From.Type == "Connection" && len(r.From.Args) > 0 {
								existingConn = r.From.Args[0]
								break
							}
							if r.To.Type == "Connection" && len(r.To.Args) > 0 {
								existingConn = r.To.Args[0]
								break
							}
						}
						if existingConn != "" {
							z.Roads = append(z.Roads, plainRoad(connectionEndpoint(existingConn), connectionEndpoint(bridgeName)))
						} else {
							z.Roads = append(z.Roads, plainRoad(connectionEndpoint(bridgeName), connectionEndpoint(bridgeName)))
						}
					} else {
						z.Roads = append(z.Roads, plainRoad(connectionEndpoint(bridgeName), connectionEndpoint(bridgeName)))
					}
				}
			}
		}
		adj[bestA][bestB] = true
		adj[bestB][bestA] = true
	}
}

func findComponents(adj map[int]map[int]bool, nodeCount int) [][]int {
	visited := make([]bool, nodeCount)
	var components [][]int
	for start := 0; start < nodeCount; start++ {
		if visited[start] {
			continue
		}
		var comp []int
		queue := []int{start}
		visited[start] = true
		for len(queue) > 0 {
			cur := queue[0]
			queue = queue[1:]
			comp = append(comp, cur)
			for nb := range adj[cur] {
				if !visited[nb] {
					visited[nb] = true
					queue = append(queue, nb)
				}
			}
		}
		components = append(components, comp)
	}
	return components
}

// ── Delaunay triangulation ───────────────────────────────────────────

func delaunayEdges(pts [][2]float64) [][2]int {
	n := len(pts)
	if n <= 1 {
		return nil
	}
	if n == 2 {
		return [][2]int{{0, 1}}
	}
	minX, minY := pts[0][0], pts[0][1]
	maxX, maxY := minX, minY
	for _, p := range pts[1:] {
		if p[0] < minX {
			minX = p[0]
		}
		if p[1] < minY {
			minY = p[1]
		}
		if p[0] > maxX {
			maxX = p[0]
		}
		if p[1] > maxY {
			maxY = p[1]
		}
	}
	dx, dy := maxX-minX, maxY-minY
	delta := math.Max(dx, dy) * 10
	superPts := make([][2]float64, len(pts)+3)
	copy(superPts, pts)
	superPts[n] = [2]float64{minX - delta, minY - delta*3}
	superPts[n+1] = [2]float64{minX + delta*3, minY - delta}
	superPts[n+2] = [2]float64{minX, minY + delta*3}

	type tri struct{ i0, i1, i2 int }
	triangles := []tri{{n, n + 1, n + 2}}

	for p := 0; p < n; p++ {
		px, py := superPts[p][0], superPts[p][1]
		var bad []tri
		for _, t := range triangles {
			if inCircumcircle(superPts, t.i0, t.i1, t.i2, px, py) {
				bad = append(bad, t)
			}
		}
		type edge struct{ a, b int }
		var boundary []edge
		for _, t := range bad {
			edges := [3]edge{{t.i0, t.i1}, {t.i1, t.i2}, {t.i2, t.i0}}
			for _, e := range edges {
				shared := false
				for _, o := range bad {
					if o == t {
						continue
					}
					if (o.i0 == e.a && o.i1 == e.b) || (o.i1 == e.a && o.i0 == e.b) ||
						(o.i1 == e.a && o.i2 == e.b) || (o.i2 == e.a && o.i1 == e.b) ||
						(o.i2 == e.a && o.i0 == e.b) || (o.i0 == e.a && o.i2 == e.b) {
						shared = true
						break
					}
				}
				if !shared {
					boundary = append(boundary, e)
				}
			}
		}
		badSet := map[tri]bool{}
		for _, t := range bad {
			badSet[t] = true
		}
		var newTris []tri
		for _, t := range triangles {
			if !badSet[t] {
				newTris = append(newTris, t)
			}
		}
		for _, e := range boundary {
			newTris = append(newTris, tri{e.a, e.b, p})
		}
		triangles = newTris
	}

	var realTris []tri
	for _, t := range triangles {
		if t.i0 < n && t.i1 < n && t.i2 < n {
			realTris = append(realTris, t)
		}
	}
	edgeSet := map[[2]int]bool{}
	for _, t := range realTris {
		addEdge := func(a, b int) {
			if a > b {
				a, b = b, a
			}
			edgeSet[[2]int{a, b}] = true
		}
		addEdge(t.i0, t.i1)
		addEdge(t.i1, t.i2)
		addEdge(t.i2, t.i0)
	}
	result := make([][2]int, 0, len(edgeSet))
	for e := range edgeSet {
		result = append(result, e)
	}
	return result
}

func inCircumcircle(pts [][2]float64, i0, i1, i2 int, px, py float64) bool {
	ax, ay := pts[i0][0]-px, pts[i0][1]-py
	bx, by := pts[i1][0]-px, pts[i1][1]-py
	cx, cy := pts[i2][0]-px, pts[i2][1]-py
	det := ax*(by*(cx*cx+cy*cy)-cy*(bx*bx+by*by)) -
		ay*(bx*(cx*cx+cy*cy)-cx*(bx*bx+by*by)) +
		(ax*ax+ay*ay)*(bx*cy-by*cx)
	return det > 0
}

// ── balanced positions ───────────────────────────────────────────────

func zoneTierRank(letter string, playerLetters []string, neutralByLetter map[string]neutralZonePlan) int {
	if contains(playerLetters, letter) {
		return 0
	}
	plan, ok := neutralByLetter[letter]
	if !ok {
		return 1
	}
	switch plan.Quality {
	case constants.QualityHigh:
		return 3
	case constants.QualityMedium:
		return 2
	default:
		return 1
	}
}

func buildBalancedRandomPositions(orderedLetters, playerLetters []string, neutralByLetter map[string]neutralZonePlan) [][2]float64 {
	count := len(orderedLetters)
	if count == 0 {
		return nil
	}

	tierRadius := func(tier int) float64 {
		switch tier {
		case 0:
			return 0.38
		case 1:
			return 0.27
		case 2:
			return 0.16
		default:
			return 0.06
		}
	}

	byTier := map[int][]int{}
	for i, l := range orderedLetters {
		t := zoneTierRank(l, playerLetters, neutralByLetter)
		byTier[t] = append(byTier[t], i)
	}

	positions := make([][2]float64, count)
	for tier, indices := range byTier {
		radius := tierRadius(tier)
		nn := len(indices)
		offset := float64(tier) * math.Pi / math.Max(1, float64(nn))
		for j, idx := range indices {
			angle := 2*math.Pi*float64(j)/float64(nn) + offset
			jitter := float64(j%3-1) * 0.008
			positions[idx] = [2]float64{
				math.Max(0.05, math.Min(0.95, 0.5+math.Cos(angle+jitter)*radius)),
				math.Max(0.05, math.Min(0.95, 0.5+math.Sin(angle+jitter)*radius)),
			}
		}
	}
	return positions
}

// ── ordered letters / balanced placement ─────────────────────────────

func buildOrderedLetters(settings *config.GeneratorConfig, playerLetters []string, neutralZones []neutralZonePlan, isRing bool) []string {
	neutralLetters := make([]string, len(neutralZones))
	for i, nz := range neutralZones {
		neutralLetters[i] = nz.Letter
	}
	if settings.Topology == config.TopologyBalanced {
		sep := 0
		if settings.MinNeutralZonesBetweenPlayers > 0 && canHonorNeutralSeparation(settings, len(neutralLetters)) {
			sep = settings.MinNeutralZonesBetweenPlayers
		}
		if isRing {
			return buildBalancedRingLetters(playerLetters, neutralZones, sep)
		}
		return buildBalancedChainLetters(playerLetters, neutralZones, sep)
	}
	minSep := settings.MinNeutralZonesBetweenPlayers
	if minSep <= 0 || settings.RandomPortals || !canHonorNeutralSeparation(settings, len(neutralLetters)) {
		return append(append([]string{}, playerLetters...), neutralLetters...)
	}
	var ordered []string
	queue := make([]string, len(neutralLetters))
	copy(queue, neutralLetters)
	qi := 0
	for i, pl := range playerLetters {
		ordered = append(ordered, pl)
		needsSep := isRing || i < len(playerLetters)-1
		if !needsSep {
			continue
		}
		for j := 0; j < minSep && qi < len(queue); j++ {
			ordered = append(ordered, queue[qi])
			qi++
		}
	}
	for qi < len(queue) {
		ordered = append(ordered, queue[qi])
		qi++
	}
	if len(ordered) == 0 {
		return append(append([]string{}, playerLetters...), neutralLetters...)
	}
	return ordered
}

func canHonorNeutralSeparation(settings *config.GeneratorConfig, neutralCount int) bool {
	min := settings.MinNeutralZonesBetweenPlayers
	if min <= 0 {
		return true
	}
	if settings.RandomPortals {
		return false
	}
	switch settings.Topology {
	case config.TopologyDefault, config.TopologyBalanced:
		return neutralCount >= settings.PlayerCount*min
	case config.TopologyChain:
		return neutralCount >= (settings.PlayerCount-1)*min
	case config.TopologyHubAndSpoke:
		return min <= 1
	case config.TopologySharedWeb:
		return min <= 1 && neutralCount >= 1
	default:
		return false
	}
}

func buildBalancedRingLetters(playerLetters []string, neutralZones []neutralZonePlan, minSep int) []string {
	if len(playerLetters) == 0 {
		return buildBalancedNeutralRing(neutralZones, 1)
	}
	if len(neutralZones) == 0 {
		return append([]string{}, playerLetters...)
	}
	caps := buildEvenGapCapacities(len(playerLetters), len(neutralZones), minSep)
	gaps := assignNeutralZonesToGaps(neutralZones, caps, false)
	var ordered []string
	for i, pl := range playerLetters {
		ordered = append(ordered, pl)
		for _, nz := range orderNeutralsWithinGap(gaps[i]) {
			ordered = append(ordered, nz.Letter)
		}
	}
	return ordered
}

func buildBalancedChainLetters(playerLetters []string, neutralZones []neutralZonePlan, minSep int) []string {
	if len(playerLetters) == 0 {
		letters := make([]string, len(neutralZones))
		for i, nz := range neutralZones {
			letters[i] = nz.Letter
		}
		return letters
	}
	gapCount := len(playerLetters) + 1
	capacities := make([]int, gapCount)
	remaining := len(neutralZones)
	reqInterior := max(0, len(playerLetters)-1) * minSep
	if minSep > 0 && len(neutralZones) >= reqInterior {
		for i := 1; i < gapCount-1; i++ {
			capacities[i] = minSep
		}
		remaining -= reqInterior
	}
	// Distribute extra neutrals only into interior gaps so that the first
	// and last positions of the chain are always player zones. Degenerate cases (0 or 1
	// player) fall back to even distribution across every gap
	interiorGapCount := max(0, gapCount-2)
	if interiorGapCount > 0 {
		extras := buildEvenGapCapacities(interiorGapCount, remaining, 0)
		for i := 1; i < gapCount-1; i++ {
			capacities[i] += extras[i-1]
		}
	} else {
		extras := buildEvenGapCapacities(gapCount, remaining, 0)
		for i := 0; i < gapCount; i++ {
			capacities[i] += extras[i]
		}
	}
	gaps := assignNeutralZonesToGaps(neutralZones, capacities, true)
	var ordered []string
	for _, nz := range orderEdgeGap(gaps[0], true) {
		ordered = append(ordered, nz.Letter)
	}
	for i, pl := range playerLetters {
		ordered = append(ordered, pl)
		gap := gaps[i+1]
		trailing := i == len(playerLetters)-1
		var g []neutralZonePlan
		if trailing {
			g = orderEdgeGap(gap, false)
		} else {
			g = orderNeutralsWithinGap(gap)
		}
		for _, nz := range g {
			ordered = append(ordered, nz.Letter)
		}
	}
	if len(ordered) == 0 {
		nl := make([]string, len(neutralZones))
		for i, nz := range neutralZones {
			nl[i] = nz.Letter
		}
		return append(append([]string{}, playerLetters...), nl...)
	}
	return ordered
}

func buildBalancedNeutralRing(neutralZones []neutralZonePlan, playerCount int) []string {
	if len(neutralZones) <= 1 {
		r := make([]string, len(neutralZones))
		for i, nz := range neutralZones {
			r[i] = nz.Letter
		}
		return r
	}
	gc := max(1, playerCount)
	caps := buildEvenGapCapacities(gc, len(neutralZones), 0)
	gaps := assignNeutralZonesToGaps(neutralZones, caps, false)
	var result []string
	for _, gap := range gaps {
		for _, nz := range orderNeutralsWithinGap(gap) {
			result = append(result, nz.Letter)
		}
	}
	return result
}

func buildEvenGapCapacities(gapCount, itemCount, minimumPerGap int) []int {
	if gapCount <= 0 {
		return nil
	}
	capacities := make([]int, gapCount)
	if itemCount <= 0 {
		return capacities
	}
	minimum := max(0, minimumPerGap)
	reserved := minimum * gapCount
	remaining := itemCount
	if minimum > 0 && itemCount >= reserved {
		for i := range capacities {
			capacities[i] = minimum
		}
		remaining -= reserved
	}
	for i := 0; i < remaining; i++ {
		gap := int(math.Floor((float64(i) + 0.5) * float64(gapCount) / float64(remaining)))
		capacities[clampInt(gap, 0, gapCount-1)]++
	}
	return capacities
}

func neutralZoneBalanceScore(zone neutralZonePlan) float64 {
	q := 1.0
	switch zone.Quality {
	case constants.QualityHigh:
		q = 3.0
	case constants.QualityMedium:
		q = 2.0
	}
	return q + math.Min(float64(zone.CastleCount), 4)*0.15
}

func assignNeutralZonesToGaps(neutralZones []neutralZonePlan, caps []int, preferInterior bool) [][]neutralZonePlan {
	gaps := make([][]neutralZonePlan, len(caps))
	loads := make([]float64, len(caps))
	sorted := make([]neutralZonePlan, len(neutralZones))
	copy(sorted, neutralZones)
	sort.SliceStable(sorted, func(i, j int) bool {
		si, sj := neutralZoneBalanceScore(sorted[i]), neutralZoneBalanceScore(sorted[j])
		if si != sj {
			return si > sj
		}
		return sorted[i].Letter < sorted[j].Letter
	})
	for _, nz := range sorted {
		var candidates []int
		for i := range caps {
			if len(gaps[i]) < caps[i] {
				candidates = append(candidates, i)
			}
		}
		if len(candidates) == 0 {
			break
		}
		if preferInterior {
			var interior []int
			for _, c := range candidates {
				if c > 0 && c < len(caps)-1 {
					interior = append(interior, c)
				}
			}
			if len(interior) > 0 {
				candidates = interior
			}
		}
		best := candidates[0]
		for _, c := range candidates[1:] {
			if loads[c] < loads[best] || (loads[c] == loads[best] && len(gaps[c]) < len(gaps[best])) || (loads[c] == loads[best] && len(gaps[c]) == len(gaps[best]) && c < best) {
				best = c
			}
		}
		gaps[best] = append(gaps[best], nz)
		loads[best] += neutralZoneBalanceScore(nz)
	}
	return gaps
}

func orderNeutralsWithinGap(neutralZones []neutralZonePlan) []neutralZonePlan {
	if len(neutralZones) <= 1 {
		return append([]neutralZonePlan{}, neutralZones...)
	}
	sorted := make([]neutralZonePlan, len(neutralZones))
	copy(sorted, neutralZones)
	sort.SliceStable(sorted, func(i, j int) bool {
		si, sj := neutralZoneBalanceScore(sorted[i]), neutralZoneBalanceScore(sorted[j])
		if si != sj {
			return si < sj
		}
		return sorted[i].Letter < sorted[j].Letter
	})
	slots := make([]neutralZonePlan, len(sorted))
	lo, hi := 0, len(sorted)-1
	for i, nz := range sorted {
		if i%2 == 0 {
			slots[lo] = nz
			lo++
		} else {
			slots[hi] = nz
			hi--
		}
	}
	return slots
}

func orderEdgeGap(neutralZones []neutralZonePlan, playerAtEnd bool) []neutralZonePlan {
	sorted := make([]neutralZonePlan, len(neutralZones))
	copy(sorted, neutralZones)
	sort.SliceStable(sorted, func(i, j int) bool {
		si, sj := neutralZoneBalanceScore(sorted[i]), neutralZoneBalanceScore(sorted[j])
		if si != sj {
			return si < sj
		}
		return sorted[i].Letter < sorted[j].Letter
	})
	if playerAtEnd {
		for i, j := 0, len(sorted)-1; i < j; i, j = i+1, j-1 {
			sorted[i], sorted[j] = sorted[j], sorted[i]
		}
	}
	return sorted
}

// ── topology adjacency (for hold city picking) ───────────────────────

func buildTopologyAdjacency(settings *config.GeneratorConfig, playerLetters []string, neutralZones []neutralZonePlan) map[string]map[string]bool {
	adj := map[string]map[string]bool{}
	link := func(a, b string) {
		if adj[a] == nil {
			adj[a] = map[string]bool{}
		}
		if adj[b] == nil {
			adj[b] = map[string]bool{}
		}
		adj[a][b] = true
		adj[b][a] = true
	}
	switch settings.Topology {
	case config.TopologyChain:
		ordered := buildOrderedLetters(settings, playerLetters, neutralZones, false)
		isolate := settings.NoDirectPlayerConnections && len(playerLetters) > 1
		playerSet := toSet(playerLetters)
		for i := 0; i < len(ordered)-1; i++ {
			if isolate && playerSet[ordered[i]] && playerSet[ordered[i+1]] {
				continue
			}
			link(ordered[i], ordered[i+1])
		}
	case config.TopologyDefault, config.TopologyBalanced:
		ordered := buildOrderedLetters(settings, playerLetters, neutralZones, true)
		isolate := settings.NoDirectPlayerConnections && len(playerLetters) > 1
		playerSet := toSet(playerLetters)
		for i := 0; i < len(ordered); i++ {
			next := (i + 1) % len(ordered)
			if isolate && playerSet[ordered[i]] && playerSet[ordered[next]] {
				continue
			}
			link(ordered[i], ordered[next])
		}
	default:
		ordered := buildOrderedLetters(settings, playerLetters, neutralZones, true)
		for i := 0; i < len(ordered); i++ {
			link(ordered[i], ordered[(i+1)%len(ordered)])
		}
	}
	return adj
}

func pickHoldCityNeutralLetter(neutralZones []neutralZonePlan, playerLetters []string, adjacency map[string]map[string]bool) string {
	if len(neutralZones) == 0 {
		return ""
	}

	bfs := func(start string) map[string]int {
		dist := map[string]int{start: 0}
		queue := []string{start}
		for len(queue) > 0 {
			cur := queue[0]
			queue = queue[1:]
			for nb := range adjacency[cur] {
				if _, ok := dist[nb]; !ok {
					dist[nb] = dist[cur] + 1
					queue = append(queue, nb)
				}
			}
		}
		return dist
	}

	var distsByPlayer []map[string]int
	for _, p := range playerLetters {
		distsByPlayer = append(distsByPlayer, bfs(p))
	}

	type candidate struct {
		letter    string
		minDist   int
		variance  float64
		quality   int
		hasCastle int
	}
	var candidates []candidate
	for _, plan := range neutralZones {
		var dists []int
		for _, d := range distsByPlayer {
			v, ok := d[plan.Letter]
			if !ok {
				v = 999999
			}
			dists = append(dists, v)
		}
		minD := dists[0]
		sum := 0
		for _, d := range dists {
			if d < minD {
				minD = d
			}
			sum += d
		}
		mean := float64(sum) / float64(len(dists))
		variance := 0.0
		for _, d := range dists {
			diff := float64(d) - mean
			variance += diff * diff
		}
		variance /= float64(len(dists))
		hc := 0
		if plan.CastleCount > 0 {
			hc = 1
		}
		candidates = append(candidates, candidate{plan.Letter, minD, variance, int(plan.Quality), hc})
	}
	sort.SliceStable(candidates, func(i, j int) bool {
		a, b := candidates[i], candidates[j]
		if a.minDist != b.minDist {
			return a.minDist > b.minDist
		}
		if a.variance != b.variance {
			return a.variance < b.variance
		}
		if a.quality != b.quality {
			return a.quality > b.quality
		}
		return a.hasCastle > b.hasCastle
	})
	return candidates[0].letter
}

// ── helpers ──────────────────────────────────────────────────────────

func contains(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}

func indexOf(ss []string, s string) int {
	for i, v := range ss {
		if v == s {
			return i
		}
	}
	return -1
}

func zoneName(letter string, playerLetters []string) string {
	if contains(playerLetters, letter) {
		return "Spawn-" + letter
	}
	return "Neutral-" + letter
}

func mapNeutralByLetter(neutralZones []neutralZonePlan) map[string]neutralZonePlan {
	m := make(map[string]neutralZonePlan, len(neutralZones))
	for _, nz := range neutralZones {
		m[nz.Letter] = nz
	}
	return m
}

func toSet(ss []string) map[string]bool {
	m := make(map[string]bool, len(ss))
	for _, s := range ss {
		m[s] = true
	}
	return m
}

func uniqueStrings(ss []string) []string {
	seen := map[string]bool{}
	var result []string
	for _, s := range ss {
		if s != "" && !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}
	return result
}

func findZone(zones []template.Zone, name string) *template.Zone {
	for i := range zones {
		if zones[i].Name == name {
			return &zones[i]
		}
	}
	return nil
}
