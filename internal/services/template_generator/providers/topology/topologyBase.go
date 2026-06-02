package topology

import (
	"fmt"
	"math"

	"github.com/Tariomka/hommoe_custom_templates/internal/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/builders/variant_content"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/utils"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
)

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

type topologyBase struct {
	zoneLabelProvider *zones.ZoneLabelProvider
}

func newTopologyBase() topologyBase {
	return topologyBase{
		zoneLabelProvider: zones.NewZoneLabelProvider(),
	}
}

func (this *topologyBase) GetSpawnZone(
	label, playerName string,
	ringConns []string,
	castleCount int,
	matchFactions bool,
	zoneSize float64,
	spawnFootholds,
	generateRoads bool,
	tuning models.GenerationTuning) template.Zone {
	mainObjects := []template.MainObject{this.createPlayerSpawnCity(playerName, tuning.ScaleByNeutralGuardStrength(5000))}
	// TODO: add player owned castles
	mainObjects = append(mainObjects, this.createPlayerUnclaimedCities(matchFactions, tuning.ScaleByNeutralGuardStrength(2500), castleCount-1)...)

	crossroadsPosition := 0
	biome := template.TypedRef{Type: "MatchMainObject", Args: []string{"0"}}
	var roads []template.Road
	if castleCount > 0 {
		roads = buildOuterZoneRoads(ringConns, castleCount, spawnFootholds, generateRoads)
	} else {
		roads = buildConnectorZoneRoads(ringConns, generateRoads)
	}
	return template.Zone{
		Name:                         "Spawn-" + label,
		Size:                         utils.NormalizeZoneSize(zoneSize),
		Layout:                       "zone_layout_spawns",
		GuardCutoffValue:             2000,
		GuardRandomization:           tuning.GuardRandomization,
		GuardMultiplier:              tuning.ScaleByNeutralGuardStrengthPrecise(1.0),
		GuardWeeklyIncrement:         0.20,
		GuardReactionDistribution:    []int{60, 20, 10, 10, 2, 0},
		DiplomacyModifier:            -0.5,
		GuardedContentPool:           linq.FromSlice(t2Guarded).ToSlice(),
		UnguardedContentPool:         linq.FromSlice(t2Unguarded).ToSlice(),
		ResourcesContentPool:         linq.FromSlice(generalResourcesPoor).ToSlice(),
		MandatoryContent:             template.StringList{"mandatory_content_side_" + label},
		ContentCountLimits:           buildSideContentLimits(),
		GuardedContentValue:          tuning.ScaleByStructureDensity(200000 * tuning.ContentScale),
		GuardedContentValuePerArea:   tuning.ScaleByStructureDensity(2000 * math.Sqrt(tuning.ContentScale)),
		UnguardedContentValue:        tuning.ScaleByStructureDensity(50000 * tuning.ContentScale),
		UnguardedContentValuePerArea: tuning.ScaleByStructureDensity(400 * math.Sqrt(tuning.ContentScale)),
		ResourcesValue:               tuning.ScaleByResourceDensity(80000 * tuning.ContentScale),
		ResourcesValuePerArea:        tuning.ScaleByResourceDensity(600 * math.Sqrt(tuning.ContentScale)),
		MainObjects:                  mainObjects,
		ZoneBiome:                    biome,
		ContentBiome:                 biome,
		MetaObjectsBiome:             biome,
		CrossroadsPosition:           &crossroadsPosition,
		Roads:                        roads,
	}
}

func (this *topologyBase) CreateNeutralZone(
	plan models.NeutralZonePlan,
	ringConns []string,
	zoneSize float64,
	spawnFootholds, generateRoads bool,
	tuning models.GenerationTuning,
	isHoldCity bool) template.Zone {
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
			Type: "City", GuardChance: 1,
			GuardValue:               tuning.ScaleByBorderGuardStrength(guardVal),
			GuardWeeklyIncrement:     0.10,
			BuildingsConstructionSid: bcsid,
			Faction:                  &template.TypedRef{Type: "FromList", Args: []string{}},
			Placement:                placement, PlacementArgs: placementArgs, HoldCityWinCon: isHoldCity,
		}
		mainObjects = append(mainObjects, mo)
	}
	for i := 1; i < castleCount; i++ {
		mainObjects = append(mainObjects, template.MainObject{
			Type: "City", GuardChance: 1,
			GuardValue:               tuning.ScaleByBorderGuardStrength(profile.ExtraCityGuardValue),
			GuardWeeklyIncrement:     0.10,
			BuildingsConstructionSid: profile.ExtraBuildingsCSid,
			Faction:                  &template.TypedRef{Type: "FromList", Args: []string{}},
			Placement:                "Uniform", PlacementArgs: []string{"false", "-0.8", "3"},
		})
	}

	reaction := []int{0, 10, 10, 10, 10, 0}
	if plan.Quality == models.QualityHigh {
		reaction = []int{0, 10, 10, 20, 10, 0}
	}

	biome := template.TypedRef{Type: "MatchZone", Args: []string{}}
	if castleCount > 0 {
		biome = template.TypedRef{Type: "MatchMainObject", Args: []string{"0"}}
	}
	crossroadsPosition := 0
	var roads []template.Road
	if castleCount > 0 {
		roads = buildOuterZoneRoads(ringConns, castleCount, spawnFootholds, generateRoads)
	} else {
		roads = buildConnectorZoneRoads(ringConns, generateRoads)
	}
	return template.Zone{
		Name:                         "Neutral-" + plan.Letter,
		Size:                         utils.NormalizeZoneSize(zoneSize),
		Layout:                       profile.Layout,
		GuardCutoffValue:             2000,
		GuardRandomization:           tuning.GuardRandomization,
		GuardMultiplier:              tuning.ScaleByNeutralGuardStrengthPrecise(profile.GuardMultiplier),
		GuardWeeklyIncrement:         0.20,
		GuardReactionDistribution:    reaction,
		DiplomacyModifier:            -0.5,
		GuardedContentPool:           profile.GuardedContentPool,
		UnguardedContentPool:         profile.UnguardedContentPool,
		ResourcesContentPool:         profile.ResourcesContentPool,
		MandatoryContent:             template.StringList{"mandatory_content_neutral_" + plan.Letter},
		ContentCountLimits:           buildSideContentLimits(),
		GuardedContentValue:          tuning.ScaleByStructureDensity(float64(profile.GuardedContentValue) * tuning.ContentScale),
		GuardedContentValuePerArea:   tuning.ScaleByStructureDensity(float64(profile.GuardedContentValuePerArea) * math.Sqrt(tuning.ContentScale)),
		UnguardedContentValue:        tuning.ScaleByStructureDensity(float64(profile.UnguardedContentValue) * tuning.ContentScale),
		UnguardedContentValuePerArea: tuning.ScaleByStructureDensity(float64(profile.UnguardedContentValuePerArea) * math.Sqrt(tuning.ContentScale)),
		ResourcesValue:               tuning.ScaleByResourceDensity(float64(profile.ResourcesValue) * tuning.ContentScale),
		ResourcesValuePerArea:        tuning.ScaleByResourceDensity(float64(profile.ResourcesValuePerArea) * math.Sqrt(tuning.ContentScale)),
		MainObjects:                  mainObjects, ZoneBiome: biome, ContentBiome: biome, MetaObjectsBiome: biome,
		CrossroadsPosition: &crossroadsPosition,
		Roads:              roads,
	}
}

func (this *topologyBase) CreateHubZone(
	spokeConns []string,
	tuning models.GenerationTuning,
	isHoldCity bool,
	size float64,
	castleCount int,
	generateRoads bool) template.Zone {
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
			Type: "City", GuardChance: gc,
			GuardValue:               tuning.ScaleByNeutralGuardStrength(gv),
			GuardWeeklyIncrement:     0.10,
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
		Name: "Hub", Size: size, Layout: "zone_layout_center",
		GuardCutoffValue: 2000, GuardRandomization: 0.05,
		GuardMultiplier: tuning.ScaleByNeutralGuardStrengthPrecise(1.5), GuardWeeklyIncrement: 0.20,
		GuardReactionDistribution: []int{0, 10, 10, 20, 10, 0}, DiplomacyModifier: -0.5,
		GuardedContentPool:   linq.FromSlice(t3Guarded).ToSlice(),
		UnguardedContentPool: linq.FromSlice(t3Unguarded).ToSlice(),
		ResourcesContentPool: linq.FromSlice(generalResourcesMedium).ToSlice(),
		MandatoryContent:     template.StringList{}, ContentCountLimits: buildSideContentLimits(),
		GuardedContentValue:          tuning.ScaleByStructureDensity(300000 * tuning.ContentScale),
		GuardedContentValuePerArea:   tuning.ScaleByStructureDensity(2400 * math.Sqrt(tuning.ContentScale)),
		UnguardedContentValue:        tuning.ScaleByStructureDensity(50000 * tuning.ContentScale),
		UnguardedContentValuePerArea: tuning.ScaleByStructureDensity(600 * math.Sqrt(tuning.ContentScale)),
		ResourcesValue:               tuning.ScaleByResourceDensity(80000 * tuning.ContentScale),
		ResourcesValuePerArea:        tuning.ScaleByResourceDensity(600 * math.Sqrt(tuning.ContentScale)),
		MainObjects:                  mainObjects, ZoneBiome: biome, ContentBiome: biome, MetaObjectsBiome: biome,
		CrossroadsPosition: &cp0, Roads: roads,
	}
}

func (this *topologyBase) createPlayerSpawnCity(playerName string, guardValue int) template.MainObject {
	return variant_content.NewObjectBuilder().
		WithTypeSpawn().
		WithSpawn(playerName).
		WithNoGuardWhenOwned().
		WithGuardChance(1).
		WithGuardValue(guardValue).
		WithGuardWeeklyIncrement(0.10).
		WithCastleQualityDefault().
		WithPlacementUniform().
		WithPlacementArgs("true", "0.7", "0").
		Build()
}

func (this *topologyBase) createPlayerUnclaimedCities(matchPlayerFaction bool, guardValue int, castleCount int) []template.MainObject {
	var castles []template.MainObject
	for range castleCount {
		object := variant_content.NewObjectBuilder().
			WithTypeCity().
			WithGuardChance(1).
			WithGuardValue(guardValue).
			WithGuardWeeklyIncrement(0.10).
			WithCastleQualityPoor().
			WithPlacementUniform().
			WithPlacementArgs("false", "-0.8", "3")
		if matchPlayerFaction {
			object = object.WithFaction("Match", "0")
		} else {
			object = object.WithFaction("Random") // TODO: is this valid?
		}
		castles = append(castles, object.Build())
	}
	return castles
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
