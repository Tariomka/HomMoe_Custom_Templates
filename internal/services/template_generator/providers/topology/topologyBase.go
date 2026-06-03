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
		this.createPlayerSpawnCastle(playerName, tuning.ScaleByNeutralGuardStrength(5000)),
	}
	// TODO: add player owned castles
	mainObjects = append(mainObjects,
		this.createPlayerUnclaimedCastles(matchFactions, tuning.ScaleByNeutralGuardStrength(2500), castleCount-1)...)

	return variant_content.NewZoneBuilder().
		WithName("Spawn-" + label).
		WithSize(utils.NormalizeZoneSize(zoneSize)).
		WithLayoutSpawns().
		WithGuardCutoffValue(2000).
		WithGuardRandomization(tuning.GuardRandomization).
		WithGuardMultiplier(tuning.ScaleByNeutralGuardStrengthPrecise(1.0)).
		WithGuardWeeklyIncrement(0.20).
		WithGuardReactionDistribution([]int{60, 20, 10, 10, 2, 0}).
		WithDiplomacyModifier(-0.5).
		WithGuardedContentPool(registry.GetGuardedContentPoolT2List()).
		WithUnguardedContentPool(registry.GetUnguardedContentPoolT2List()).
		WithResourcesContentPool([]string{resourceContentPool.StartZonePoor}).
		WithMandatoryContent("mandatory_content_side_" + label).
		WithContentCountLimits(buildSideContentLimits()).
		WithGuardedContentValue(tuning.ScaleByStructureDensity(200000 * tuning.ContentScale)).
		WithGuardedContentValuePerArea(tuning.ScaleByStructureDensity(2000 * math.Sqrt(tuning.ContentScale))).
		WithUnguardedContentValue(tuning.ScaleByStructureDensity(50000 * tuning.ContentScale)).
		WithUnguardedContentValuePerArea(tuning.ScaleByStructureDensity(400 * math.Sqrt(tuning.ContentScale))).
		WithResourcesValue(tuning.ScaleByResourceDensity(80000 * tuning.ContentScale)).
		WithResourcesValuePerArea(tuning.ScaleByResourceDensity(600 * math.Sqrt(tuning.ContentScale))).
		WithMainObjects(mainObjects).
		WithBiomeMatchMainObject("0").
		WithCrossroadsPosition(0).
		WithRoads(this.createOuterZoneRoads(ringConns, castleCount, spawnFootholds, generateRoads)).
		Build()
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

	zoneBuilder := variant_content.NewZoneBuilder().
		WithName("Neutral-" + plan.Label).
		WithSize(utils.NormalizeZoneSize(zoneSize)).
		WithLayout(profile.Layout).
		WithGuardCutoffValue(2000).
		WithGuardRandomization(tuning.GuardRandomization).
		WithGuardMultiplier(tuning.ScaleByNeutralGuardStrengthPrecise(profile.GuardMultiplier)).
		WithGuardWeeklyIncrement(0.20).
		WithGuardReactionDistribution(profile.GuardReactionDistribution).
		WithDiplomacyModifier(-0.5).
		WithGuardedContentPool(profile.GuardedContentPool).
		WithUnguardedContentPool(profile.UnguardedContentPool).
		WithResourcesContentPool(profile.ResourcesContentPool).
		WithMandatoryContent("mandatory_content_neutral_" + plan.Label).
		WithContentCountLimits(buildSideContentLimits()).
		WithGuardedContentValue(tuning.ScaleByStructureDensity(float64(profile.GuardedContentValue) * tuning.ContentScale)).
		WithGuardedContentValuePerArea(tuning.ScaleByStructureDensity(float64(profile.GuardedContentValuePerArea) * math.Sqrt(tuning.ContentScale))).
		WithUnguardedContentValue(tuning.ScaleByStructureDensity(float64(profile.UnguardedContentValue) * tuning.ContentScale)).
		WithUnguardedContentValuePerArea(tuning.ScaleByStructureDensity(float64(profile.UnguardedContentValuePerArea) * math.Sqrt(tuning.ContentScale))).
		WithResourcesValue(tuning.ScaleByResourceDensity(float64(profile.ResourcesValue) * tuning.ContentScale)).
		WithResourcesValuePerArea(tuning.ScaleByResourceDensity(float64(profile.ResourcesValuePerArea) * math.Sqrt(tuning.ContentScale))).
		WithMainObjects(this.createNeutralZoneCastles(profile, tuning, plan.CastleCount, isHoldCity)).
		WithCrossroadsPosition(0).
		WithRoads(this.createOuterZoneRoads(ringConns, plan.CastleCount, spawnFootholds, generateRoads))

	if plan.CastleCount > 0 {
		zoneBuilder = zoneBuilder.WithBiomeMatchMainObject("0")
	} else {
		zoneBuilder = zoneBuilder.WithBiomeMatchZone()
	}

	return zoneBuilder.Build()
}

func (this *topologyBase) CreateHubZone(
	connectionNames []string,
	tuning models.GenerationTuning,
	isHoldCity bool,
	size float64,
	castleCount int,
	generateRoads bool) template.Zone {
	if isHoldCity && castleCount < 1 {
		castleCount = 1
	}

	zoneBuilder := variant_content.NewZoneBuilder().
		WithName("Hub").
		WithSize(utils.NormalizeZoneSize(size)).
		WithLayoutCenter().
		WithGuardCutoffValue(2000).
		WithGuardRandomization(0.05).
		WithGuardMultiplier(tuning.ScaleByNeutralGuardStrengthPrecise(1.5)).
		WithGuardWeeklyIncrement(0.20).
		WithGuardReactionDistribution([]int{0, 10, 10, 20, 10, 0}).
		WithDiplomacyModifier(-0.5).
		WithGuardedContentPool(registry.GetGuardedContentPoolT3List()).
		WithUnguardedContentPool(registry.GetUnguardedContentPoolT3List()).
		WithResourcesContentPool([]string{resourceContentPool.StartZoneMedium}).
		WithContentCountLimits(buildSideContentLimits()).
		WithGuardedContentValue(tuning.ScaleByStructureDensity(300000 * tuning.ContentScale)).
		WithGuardedContentValuePerArea(tuning.ScaleByStructureDensity(2400 * math.Sqrt(tuning.ContentScale))).
		WithUnguardedContentValue(tuning.ScaleByStructureDensity(50000 * tuning.ContentScale)).
		WithUnguardedContentValuePerArea(tuning.ScaleByStructureDensity(600 * math.Sqrt(tuning.ContentScale))).
		WithResourcesValue(tuning.ScaleByResourceDensity(80000 * tuning.ContentScale)).
		WithResourcesValuePerArea(tuning.ScaleByResourceDensity(600 * math.Sqrt(tuning.ContentScale))).
		WithMainObjects(this.createHubZoneCastles(tuning, castleCount, isHoldCity)).
		WithCrossroadsPosition(0).
		WithRoads(this.createOuterZoneRoads(connectionNames, castleCount, false, generateRoads))

	if castleCount > 0 {
		zoneBuilder = zoneBuilder.WithBiomeMatchMainObject("0")
	} else {
		zoneBuilder = zoneBuilder.WithBiomeMatchZone()
	}

	return zoneBuilder.Build()
}

func (this *topologyBase) createPlayerSpawnCastle(playerName string, guardValue int) template.MainObject {
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

// func (this *topologyBase) createPlayerOwnedCastles(playerIndex int, guardValue int, castleCount int) []template.MainObject {  }

func (this *topologyBase) createPlayerUnclaimedCastles(
	matchPlayerFaction bool,
	guardValue, castleCount int) []template.MainObject {
	var castles []template.MainObject
	for range castleCount {
		objectBuilder := variant_content.NewObjectBuilder().
			WithTypeCity().
			WithGuardChance(1).
			WithGuardValue(guardValue).
			WithGuardWeeklyIncrement(0.10).
			WithCastleQualityPoor().
			WithPlacementUniform().
			WithPlacementArgs("false", "-0.8", "3")
		if matchPlayerFaction {
			objectBuilder = objectBuilder.WithFaction("Match", "0")
		} else {
			objectBuilder = objectBuilder.WithFaction("Random") // TODO: is this valid?
		}
		castles = append(castles, objectBuilder.Build())
	}
	return castles
}

func (this *topologyBase) createNeutralZoneCastles(
	profile models.NeutralZoneProfile,
	tuning models.GenerationTuning,
	castleCount int,
	isHoldCityZone bool) []template.MainObject {
	var castles []template.MainObject

	if castleCount > 0 {
		objectBuilder := variant_content.NewObjectBuilder().
			WithTypeCity().
			WithGuardChance(1).
			WithGuardWeeklyIncrement(0.10).
			WithFaction("FromList")
		if isHoldCityZone {
			objectBuilder = objectBuilder.
				WithGuardValue(tuning.ScaleByBorderGuardStrength(max(profile.PrimaryCityGuardValue, 20_000))).
				WithCastleQualityUltraRich().
				WithPlacementCenter().
				WithHoldCityWinCon()
		} else {
			objectBuilder = objectBuilder.
				WithGuardValue(tuning.ScaleByBorderGuardStrength(profile.PrimaryCityGuardValue)).
				WithCastleQuality(profile.PrimaryBuildingsCSid).
				WithPlacementUniform().
				WithPlacementArgs("true", "0.8", "2")
		}
		castles = append(castles, objectBuilder.Build())
	}

	for range castleCount - 1 {
		castles = append(castles,
			variant_content.NewObjectBuilder().
				WithTypeCity().
				WithGuardChance(1).
				WithGuardValue(tuning.ScaleByBorderGuardStrength(profile.ExtraCityGuardValue)).
				WithGuardWeeklyIncrement(0.10).
				WithCastleQuality(profile.ExtraBuildingsCSid).
				WithFaction("FromList").
				WithPlacementUniform().
				WithPlacementArgs("false", "-0.8", "3").
				Build())
	}

	return castles
}

func (this *topologyBase) createHubZoneCastles(
	tuning models.GenerationTuning,
	castleCount int,
	isHoldCityZone bool) []template.MainObject {
	var castles []template.MainObject
	newCastleBuilder := func() *variant_content.MainObjectBuilder {
		return variant_content.NewObjectBuilder().
			WithTypeCity().
			WithGuardWeeklyIncrement(0.10).
			WithFaction("FromList")
	}
	buildHoldCityCastle := func(builder *variant_content.MainObjectBuilder) template.MainObject {
		return builder.
			WithGuardChance(1).
			WithGuardValue(tuning.ScaleByNeutralGuardStrength(25_000)).
			WithCastleQualityUltraRich().
			WithPlacementCenter().
			WithHoldCityWinCon().
			Build()
	}
	buildCastle := func(builder *variant_content.MainObjectBuilder) template.MainObject {
		return builder.
			WithGuardChance(0.5).
			WithGuardValue(tuning.ScaleByNeutralGuardStrength(16_000)).
			WithCastleQualityRich().
			WithPlacementUniform().
			WithPlacementArgs("true", "0.8", "2").
			Build()
	}

	if castleCount > 0 && isHoldCityZone {
		castles = append(castles, buildHoldCityCastle(newCastleBuilder()))
	} else if castleCount > 0 {
		castles = append(castles, buildCastle(newCastleBuilder()))
	}

	for range castleCount - 1 {
		castles = append(castles, buildCastle(newCastleBuilder()))
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
