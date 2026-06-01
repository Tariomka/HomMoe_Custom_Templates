package topology

import (
	"math"

	"github.com/Tariomka/hommoe_custom_templates/internal/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
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
	letter, player string,
	ringConns []string, castleCount int, matchFactions bool, zoneSize float64,
	spawnFootholds,
	generateRoads bool, tuning models.GenerationTuning) template.Zone {
	mainObjects := []template.MainObject{
		{
			Type: "Spawn", Spawn: player, RemoveGuardIfHasOwner: true,
			GuardChance:              1,
			GuardValue:               tuning.ScaleByNeutralGuardStrength(5000),
			GuardWeeklyIncrement:     0.10,
			BuildingsConstructionSid: "default_buildings_construction",
			Placement:                "Uniform", PlacementArgs: []string{"true", "0.7", "0"},
		},
	}
	for i := 1; i < castleCount; i++ {
		// TODO: add player own castles
		mo := template.MainObject{
			Type: "City", GuardChance: 1, GuardValue: tuning.ScaleByNeutralGuardStrength(2500),
			GuardWeeklyIncrement:     0.10,
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
	crossroadsPosition := 0
	biome := template.TypedRef{Type: "MatchMainObject", Args: []string{"0"}}
	var roads []template.Road
	if castleCount > 0 {
		roads = buildOuterZoneRoads(ringConns, castleCount, spawnFootholds, generateRoads)
	} else {
		roads = buildConnectorZoneRoads(ringConns, generateRoads)
	}
	return template.Zone{
		Name:                         "Spawn-" + letter,
		Size:                         utils.NormalizeZoneSize(zoneSize),
		Layout:                       "zone_layout_spawns",
		GuardCutoffValue:             2000,
		GuardRandomization:           tuning.GuardRandomization,
		GuardMultiplier:              scaleGuardMultiplier(1.0, tuning),
		GuardWeeklyIncrement:         0.20,
		GuardReactionDistribution:    []int{60, 20, 10, 10, 2, 0},
		DiplomacyModifier:            -0.5,
		GuardedContentPool:           linq.FromSlice(t2Guarded).ToSlice(),
		UnguardedContentPool:         linq.FromSlice(t2Unguarded).ToSlice(),
		ResourcesContentPool:         linq.FromSlice(generalResourcesPoor).ToSlice(),
		MandatoryContent:             template.StringList{"mandatory_content_side_" + letter},
		ContentCountLimits:           buildSideContentLimits(),
		GuardedContentValue:          scaleStructureValue(200000*tuning.ContentScale, tuning),
		GuardedContentValuePerArea:   scaleStructureValue(2000*math.Sqrt(tuning.ContentScale), tuning),
		UnguardedContentValue:        scaleStructureValue(50000*tuning.ContentScale, tuning),
		UnguardedContentValuePerArea: scaleStructureValue(400*math.Sqrt(tuning.ContentScale), tuning),
		ResourcesValue:               scaleResourceValue(80000*tuning.ContentScale, tuning),
		ResourcesValuePerArea:        scaleResourceValue(600*math.Sqrt(tuning.ContentScale), tuning),
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
	if plan.Quality == models.QualityHigh {
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
