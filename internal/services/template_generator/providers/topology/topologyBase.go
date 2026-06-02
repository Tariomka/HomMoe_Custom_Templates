package topology

import (
	"fmt"
	"math"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/builders/variant_content"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/utils"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
)

var resourceContentPool = registry.GetResourcesContentPoolValues()

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
	mainObjects := []template.MainObject{
		this.createPlayerSpawnCity(playerName, tuning.ScaleByNeutralGuardStrength(5000)),
	}
	// TODO: add player owned castles
	mainObjects = append(mainObjects,
		this.createPlayerUnclaimedCities(matchFactions, tuning.ScaleByNeutralGuardStrength(2500), castleCount-1)...)

	crossroadsPosition := 0
	biome := variant_content.NewRefBuilder().BuildBiomeMatchMainObjectType("0")
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
		GuardedContentPool:           registry.GetGuardedContentPoolT2List(),
		UnguardedContentPool:         registry.GetUnguardedContentPoolT2List(),
		ResourcesContentPool:         []string{resourceContentPool.StartZonePoor},
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
		Roads:                        this.createOuterZoneRoads(ringConns, castleCount, spawnFootholds, generateRoads),
	}
}

func (this *topologyBase) CreateNeutralZone(
	plan models.NeutralZonePlan,
	ringConns []string,
	zoneSize float64,
	spawnFootholds, generateRoads bool,
	tuning models.GenerationTuning,
	isHoldCity bool) template.Zone {
	if isHoldCity && plan.CastleCount < 1 {
		plan.CastleCount = 1
	}
	profile := models.NewNeutralZoneProfile(plan.Quality)

	var mainObjects []template.MainObject
	if plan.CastleCount > 0 {
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
	for i := 1; i < plan.CastleCount; i++ {
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
	if plan.CastleCount > 0 {
		biome = template.TypedRef{Type: "MatchMainObject", Args: []string{"0"}}
	}
	crossroadsPosition := 0
	return template.Zone{
		Name:                         "Neutral-" + plan.Label,
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
		MandatoryContent:             template.StringList{"mandatory_content_neutral_" + plan.Label},
		ContentCountLimits:           buildSideContentLimits(),
		GuardedContentValue:          tuning.ScaleByStructureDensity(float64(profile.GuardedContentValue) * tuning.ContentScale),
		GuardedContentValuePerArea:   tuning.ScaleByStructureDensity(float64(profile.GuardedContentValuePerArea) * math.Sqrt(tuning.ContentScale)),
		UnguardedContentValue:        tuning.ScaleByStructureDensity(float64(profile.UnguardedContentValue) * tuning.ContentScale),
		UnguardedContentValuePerArea: tuning.ScaleByStructureDensity(float64(profile.UnguardedContentValuePerArea) * math.Sqrt(tuning.ContentScale)),
		ResourcesValue:               tuning.ScaleByResourceDensity(float64(profile.ResourcesValue) * tuning.ContentScale),
		ResourcesValuePerArea:        tuning.ScaleByResourceDensity(float64(profile.ResourcesValuePerArea) * math.Sqrt(tuning.ContentScale)),
		MainObjects:                  mainObjects, ZoneBiome: biome, ContentBiome: biome, MetaObjectsBiome: biome,
		CrossroadsPosition: &crossroadsPosition,
		Roads:              this.createOuterZoneRoads(ringConns, plan.CastleCount, spawnFootholds, generateRoads),
	}
}

func (this *topologyBase) CreateHubZone(
	connectionNames []string,
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
	crossroadsPosition := 0
	return template.Zone{
		Name: "Hub", Size: size, Layout: "zone_layout_center",
		GuardCutoffValue: 2000, GuardRandomization: 0.05,
		GuardMultiplier: tuning.ScaleByNeutralGuardStrengthPrecise(1.5), GuardWeeklyIncrement: 0.20,
		GuardReactionDistribution: []int{0, 10, 10, 20, 10, 0}, DiplomacyModifier: -0.5,
		GuardedContentPool:   registry.GetGuardedContentPoolT3List(),
		UnguardedContentPool: registry.GetUnguardedContentPoolT3List(),
		ResourcesContentPool: []string{resourceContentPool.StartZoneMedium},
		MandatoryContent:     template.StringList{}, ContentCountLimits: buildSideContentLimits(),
		GuardedContentValue:          tuning.ScaleByStructureDensity(300000 * tuning.ContentScale),
		GuardedContentValuePerArea:   tuning.ScaleByStructureDensity(2400 * math.Sqrt(tuning.ContentScale)),
		UnguardedContentValue:        tuning.ScaleByStructureDensity(50000 * tuning.ContentScale),
		UnguardedContentValuePerArea: tuning.ScaleByStructureDensity(600 * math.Sqrt(tuning.ContentScale)),
		ResourcesValue:               tuning.ScaleByResourceDensity(80000 * tuning.ContentScale),
		ResourcesValuePerArea:        tuning.ScaleByResourceDensity(600 * math.Sqrt(tuning.ContentScale)),
		MainObjects:                  mainObjects, ZoneBiome: biome, ContentBiome: biome, MetaObjectsBiome: biome,
		CrossroadsPosition: &crossroadsPosition,
		Roads:              this.createOuterZoneRoads(connectionNames, effectiveCastleCount, false, generateRoads),
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

func (this *topologyBase) createPlayerUnclaimedCities(
	matchPlayerFaction bool,
	guardValue, castleCount int) []template.MainObject {
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

func (this *topologyBase) createOuterZoneRoads(
	connectionNames []string,
	castleCount int,
	includeFoothold, generateRoads bool) []template.Road {
	if !generateRoads {
		return nil
	}

	if castleCount == 0 {
		return this.createConnectorZoneRoads(connectionNames, generateRoads)
	}

	var roads []template.Road
	for i := 1; i < castleCount; i++ {
		roads = append(roads,
			variant_content.NewRoadBuilder().
				WithFrom(variant_content.NewRefBuilder().BuildMainObjectType("0")).
				WithTo(variant_content.NewRefBuilder().BuildMainObjectType(fmt.Sprintf("%d", i))).
				Build())
	}
	if includeFoothold {
		roads = append(roads,
			variant_content.NewRoadBuilder().
				WithFrom(variant_content.NewRefBuilder().BuildMainObjectType("0")).
				WithTo(variant_content.NewRefBuilder().BuildMandatoryContentType("name_remote_foothold_1")).
				Build())
	}
	for _, name := range connectionNames {
		roads = append(roads,
			variant_content.NewRoadBuilder().
				WithFrom(variant_content.NewRefBuilder().BuildMainObjectType("0")).
				WithTo(variant_content.NewRefBuilder().BuildConnectionType(name)).
				Build())
	}
	return roads
}

func (this *topologyBase) createConnectorZoneRoads(connectionNames []string, generateRoads bool) []template.Road {
	if !generateRoads {
		return nil
	}

	distinctNames := helpers.GetUniqueElements(connectionNames)
	if len(distinctNames) == 0 {
		return nil
	}

	if len(distinctNames) == 1 {
		return []template.Road{
			variant_content.NewRoadBuilder().
				WithFrom(variant_content.NewRefBuilder().BuildConnectionType(distinctNames[0])).
				WithTo(variant_content.NewRefBuilder().BuildConnectionType(distinctNames[0])).
				Build()}
	}
	var roads []template.Road
	for _, name := range distinctNames[1:] {
		roads = append(roads,
			variant_content.NewRoadBuilder().
				WithFrom(variant_content.NewRefBuilder().BuildConnectionType(distinctNames[0])).
				WithTo(variant_content.NewRefBuilder().BuildConnectionType(name)).
				Build())
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
