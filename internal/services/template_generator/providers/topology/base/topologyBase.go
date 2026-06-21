package base

import (
	"fmt"
	"math"
	"math/rand/v2"
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/builders/placement_rule"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/builders/variant_content"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base/utils"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
)

var resourceContentPool = registry.GetResourcesContentPoolValues()

type TopologyBase struct {
	ZoneLabelProvider *zones.ZoneLabelProvider
}

func NewTopologyBase() TopologyBase {
	return TopologyBase{
		ZoneLabelProvider: zones.NewZoneLabelProvider(),
	}
}

func (this *TopologyBase) CreateVariant(
	playerLabels []string,
	firstLabel string,
	zoneCount int,
	zones []entities.Zone,
	connections []entities.Connection) entities.Variant {
	orientationBuilder := variant_content.NewOrientationBuilder().
		WithBaseAngleMin(45).
		WithBaseAngleMax(45).
		WithRandomAngleAmplitude(360)

	if zoneCount > 0 {
		orientationBuilder.WithRandomAngleStep(360 / zoneCount)
	}

	if slices.Contains(playerLabels, firstLabel) {
		orientationBuilder.WithZeroAngleZone("Spawn-" + firstLabel)
	} else {
		orientationBuilder.WithZeroAngleZone("Neutral-" + firstLabel)
	}

	return variant_content.NewVariantBuilder().
		WithOrientation(orientationBuilder.Build()).
		WithBorder(variant_content.NewBorderBuilder().
			WithCornerRadius(0).
			WithObstaclesWidth(3).
			WithObstaclesNoise(1, 12).
			WithWaterWidth(0).
			WithWaterNoise(1, 12).
			WithWaterTypeWaterGrass().
			Build()).
		WithZones(zones...).
		WithConnections(connections...).
		Build()
}

func (this *TopologyBase) CreateSpawnZone(
	label, playerName string,
	connectionNames []string,
	castleCount int,
	matchFactions bool,
	zoneSize float64,
	spawnFootholds,
	generateRoads bool,
	tuning models.GenerationTuning) entities.Zone {
	mainObjects := []entities.MainObject{
		this.createPlayerSpawnCastle(playerName, tuning.ScaleByNeutralGuardStrength(5000)),
	}
	mainObjects = append(mainObjects,
		this.createPlayerOwnedCastles(matchFactions, playerName, tuning.PlayerOwnedCastles)...)
	mainObjects = append(mainObjects,
		this.createPlayerUnclaimedCastles(matchFactions, tuning.ScaleByNeutralGuardStrength(5000), castleCount)...)

	// Roads connect the spawn castle (main object 0) to every other castle in
	// the zone; player-owned extras are road-linked just like unclaimed ones.
	// A zone with no extra castles at all stays a pass-through connector.
	roadCastleCount := castleCount + tuning.PlayerOwnedCastles

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
		WithRoads(this.createOuterZoneRoads(connectionNames, roadCastleCount, spawnFootholds, generateRoads)).
		Build()
}

func (this *TopologyBase) CreateNeutralZone(
	plan models.NeutralZonePlan,
	connectionNames []string,
	zoneSize float64,
	spawnFootholds, generateRoads bool,
	tuning models.GenerationTuning,
	isHoldCity bool) entities.Zone {
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
		WithMainObjects(CreateNeutralZoneCastles(profile, tuning, plan.CastleCount, isHoldCity)).
		WithCrossroadsPosition(0).
		WithRoads(this.createOuterZoneRoads(connectionNames, plan.CastleCount, spawnFootholds, generateRoads))

	if plan.CastleCount > 0 {
		zoneBuilder = zoneBuilder.WithBiomeMatchMainObject("0")
	} else {
		zoneBuilder = zoneBuilder.WithBiomeMatchZone()
	}

	return zoneBuilder.Build()
}

func (this *TopologyBase) CreateHubZone(
	connectionNames []string,
	tuning models.GenerationTuning,
	isHoldCity bool,
	size float64,
	castleCount int,
	generateRoads bool) entities.Zone {
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

func (this *TopologyBase) CreateRandomPortalConnections(
	playerLabels, orderedLabels []string,
	tuning models.GenerationTuning,
	maxCount int) []entities.Connection {
	count := len(orderedLabels)
	if count < 2 {
		return nil
	}
	dest := buildNonAdjacentDerangement(count)
	indices := make([]int, count)
	for i := range count {
		indices[i] = i
	}
	rand.Shuffle(len(indices), func(i, j int) { indices[i], indices[j] = indices[j], indices[i] })

	rule := placement_rule.NewPlacementRuleBuilder().BuildCrossroadsRule(placement_rule.DistanceNear, 2)
	var conns []entities.Connection
	for i := range min(count, maxCount) {
		idx := indices[i]
		fromLabel := orderedLabels[idx]
		toLabel := orderedLabels[dest[idx]]
		fromName := this.ZoneLabelProvider.CreateZoneName(fromLabel, playerLabels)
		toName := this.ZoneLabelProvider.CreateZoneName(toLabel, playerLabels)
		conns = append(conns, variant_content.NewConnectionBuilder().
			WithName(fmt.Sprintf("Portal-%s-%s", fromLabel, toLabel)).
			WithFrom(fromName).
			WithTo(toName).
			WithConnectionTypePortal().
			WithPortalPlacementRulesFrom(rule).
			WithPortalPlacementRulesTo(rule).
			WithRoad(true).
			WithGuardValue(tuning.ScaleByBorderGuardStrength(25000)).
			WithGuardWeeklyIncrement(0.15).
			Build())
	}
	return conns
}

func (this *TopologyBase) CreateMissingPlayerConnections(
	playerLabels []string,
	zones []entities.Zone,
	connections []entities.Connection,
	tuning models.GenerationTuning) []entities.Connection {
	if len(playerLabels) < 2 {
		return nil
	}

	createFallbackConnName := func(label string, partner string) string {
		if label > partner {
			label, partner = partner, label
		}
		return fmt.Sprintf("Fallback-%s-%s", label, partner)
	}

	connNames := map[string]bool{}
	for _, c := range connections {
		if c.Name != "" {
			connNames[c.Name] = true
		}
	}
	var additionalConns []entities.Connection
	for _, label := range playerLabels {
		zoneName := "Spawn-" + label
		zone, ok := linq.FromSlice(zones).First(func(z entities.Zone) bool { return z.Name == zoneName })
		if !ok {
			continue
		}

		hasConn := false
		for _, r := range zone.Roads {
			if r.To.Type == "Connection" && len(r.To.Args) > 0 && connNames[r.To.Args[0]] {
				hasConn = true
				break
			}
		}
		if hasConn {
			continue
		}

		partner := linq.FromSlice(playerLabels).FirstOrDefault(func(x string) bool { return x != label })
		if partner == "" {
			continue
		}

		fallbackName := createFallbackConnName(label, partner)
		if connNames[fallbackName] {
			continue
		}

		additionalConns = append(additionalConns, entities.Connection{
			Name: fallbackName, From: zoneName, To: "Spawn-" + partner,
			ConnectionType: "Direct", GuardZone: zoneName, SimTurnSquad: true,
			GuardValue: this.GetBorderGuardValue(label, partner, playerLabels, nil, tuning), GuardWeeklyIncrement: 0.15,
			GuardMatchGroup: "fallback_guard_" + fallbackName,
		})
		connNames[fallbackName] = true
		for _, pl := range []string{label, partner} {
			if pz, ok := linq.FromSlice(zones).First(func(z entities.Zone) bool { return z.Name == "Spawn-"+pl }); ok {
				pz.Roads = append(pz.Roads, variant_content.NewRoadBuilder().
					WithFrom(variant_content.NewRefBuilder().BuildMainObjectType("0")).
					WithTo(variant_content.NewRefBuilder().BuildConnectionType(fallbackName)).
					Build())
			}
		}
	}
	return additionalConns
}

func (this *TopologyBase) CreateMissingConnections(
	playerLabels, allLabels []string,
	positions models.Positions,
	zones []entities.Zone,
	connections []entities.Connection,
	tuning models.GenerationTuning,
	neutralZones models.NeutralZonePlans) []entities.Connection {
	if len(allLabels) < 2 {
		return nil
	}

	getFallbackConnName := func(zone entities.Zone) string {
		existingConn := ""
		for _, road := range zone.Roads {
			if road.From.Type == "Connection" && len(road.From.Args) > 0 {
				existingConn = road.From.Args[0]
				break
			}
			if road.To.Type == "Connection" && len(road.To.Args) > 0 {
				existingConn = road.To.Args[0]
				break
			}
		}
		return existingConn
	}

	adjacency := models.NewZoneIndexAdjacency(len(allLabels))
	// TODO: move out to a separate function
	zoneNameToIdx := map[string]int{}
	for index, label := range allLabels {
		zoneNameToIdx[this.ZoneLabelProvider.CreateZoneName(label, playerLabels)] = index
	}
	for connection := range linq.FromSlice(connections).
		Where(func(x entities.Connection) bool { return x.ConnectionType == "Direct" || x.ConnectionType == "Portal" }).
		Where(func(x entities.Connection) bool {
			_, okA := zoneNameToIdx[x.From]
			_, okB := zoneNameToIdx[x.To]
			return okA && okB
		}).Iterate {
		adjacency.Link(zoneNameToIdx[connection.From], zoneNameToIdx[connection.To])
	}

	connNameSet := map[string]bool{}
	for connection := range linq.FromSlice(connections).
		Where(func(x entities.Connection) bool { return x.Name != "" }).Iterate {
		connNameSet[connection.Name] = true
	}

	var additionalConns []entities.Connection
	for {
		components := adjacency.FindIndexes(len(allLabels))
		bestIndexA, bestIndexB, ok := positions.GetShortestDistanceIndex(components)
		if !ok {
			break
		}

		labelA, labelB := allLabels[bestIndexA], allLabels[bestIndexB]
		if labelA > labelB {
			labelA, labelB = labelB, labelA
		}
		bridgeName := fmt.Sprintf("Bridge-%s-%s", labelA, labelB)
		if connNameSet[bridgeName] {
			continue
		}

		zoneFrom := this.ZoneLabelProvider.CreateZoneName(allLabels[bestIndexA], playerLabels)
		zoneTo := this.ZoneLabelProvider.CreateZoneName(allLabels[bestIndexB], playerLabels)
		additionalConns = append(additionalConns, variant_content.NewConnectionBuilder().
			WithName(bridgeName).
			WithFrom(zoneFrom).
			WithTo(zoneTo).
			WithConnectionTypeDirect().
			WithGuardZone(zoneFrom).
			WithSimTurnSquad().
			WithGuardValue(this.GetBorderGuardValue(labelA, labelB, playerLabels, neutralZones, tuning)).
			WithGuardWeeklyIncrement(0.15).
			WithGuardMatchGroup(fmt.Sprintf("bridge_guard_%s-%s", labelA, labelB)).
			Build())
		connNameSet[bridgeName] = true

		for _, zoneName := range []string{zoneFrom, zoneTo} {
			if zone, ok := linq.FromSlice(zones).First(func(x entities.Zone) bool { return x.Name == zoneName }); ok {
				roadBuilder := variant_content.NewRoadBuilder().WithTo(
					variant_content.NewRefBuilder().BuildConnectionType(bridgeName))
				if len(zone.MainObjects) > 0 {
					zone.Roads = append(zone.Roads,
						roadBuilder.WithFrom(variant_content.NewRefBuilder().BuildMainObjectType("0")).Build())
				} else if len(zone.Roads) > 0 {
					connectionName := getFallbackConnName(zone)
					if connectionName == "" {
						connectionName = bridgeName
					}
					zone.Roads = append(zone.Roads,
						roadBuilder.WithFrom(variant_content.NewRefBuilder().BuildConnectionType(connectionName)).Build())
				} else {
					zone.Roads = append(zone.Roads,
						roadBuilder.WithFrom(variant_content.NewRefBuilder().BuildConnectionType(bridgeName)).Build())
				}
			}
		}

		adjacency.Link(bestIndexA, bestIndexB)
	}

	return additionalConns
}

func (this *TopologyBase) CreateConnectorZoneRoads(connectionNames []string, generateRoads bool) []entities.Road {
	if !generateRoads {
		return nil
	}

	distinctNames := helpers.GetUniqueElements(connectionNames)
	if len(distinctNames) == 0 {
		return nil
	}

	if len(distinctNames) == 1 {
		return []entities.Road{
			variant_content.NewRoadBuilder().
				WithFrom(variant_content.NewRefBuilder().BuildConnectionType(distinctNames[0])).
				WithTo(variant_content.NewRefBuilder().BuildConnectionType(distinctNames[0])).
				Build()}
	}
	var roads []entities.Road
	for _, name := range distinctNames[1:] {
		roads = append(roads,
			variant_content.NewRoadBuilder().
				WithFrom(variant_content.NewRefBuilder().BuildConnectionType(distinctNames[0])).
				WithTo(variant_content.NewRefBuilder().BuildConnectionType(name)).
				Build())
	}
	return roads
}

func (this *TopologyBase) GetBorderGuardValue(
	labelA, labelB string,
	playerLabels []string,
	neutralZones models.NeutralZonePlans,
	tuning models.GenerationTuning) int {
	aIsPlayer := slices.Contains(playerLabels, labelA)
	bIsPlayer := slices.Contains(playerLabels, labelB)
	if aIsPlayer && bIsPlayer {
		return tuning.ScaleByBorderGuardStrength(30_000)
	}

	if !aIsPlayer && !bIsPlayer {
		qualityA := neutralZones.GetQuality(labelA)
		qualityB := neutralZones.GetQuality(labelB)
		higher := qualityA
		if int(qualityB) > int(qualityA) {
			higher = qualityB
		}
		return tuning.ScaleByBorderGuardStrength(higher.GetGuardValue())
	}

	neutralLabel := labelB
	if !aIsPlayer {
		neutralLabel = labelA
	}
	return tuning.ScaleByBorderGuardStrength(neutralZones.GetQuality(neutralLabel).GetGuardValue())
}

func (this *TopologyBase) createPlayerSpawnCastle(playerName string, guardValue int) entities.MainObject {
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

// createPlayerOwnedCastles builds the extra City castles that the player owns
// from the very start of the game. Because they already have an owner, their
// guards are dropped immediately so the player can use them right away.
func (this *TopologyBase) createPlayerOwnedCastles(
	matchPlayerFaction bool,
	owner string,
	castleCount int) []entities.MainObject {
	var castles []entities.MainObject
	for range castleCount {
		objectBuilder := variant_content.NewObjectBuilder().
			WithTypeCity().
			WithOwner(owner).
			WithCastleQualityPoor().
			WithPlacementUniform()
		if matchPlayerFaction {
			objectBuilder = objectBuilder.WithFaction("Match", "0")
		} else {
			objectBuilder = objectBuilder.WithFaction("Random")
		}
		castles = append(castles, objectBuilder.Build())
	}
	return castles
}

// createPlayerUnclaimedCastles builds the extra neutral City castles that sit
// inside a player's zone but stay unowned until someone captures them.
func (this *TopologyBase) createPlayerUnclaimedCastles(
	matchPlayerFaction bool,
	guardValue, castleCount int) []entities.MainObject {
	var castles []entities.MainObject
	for range castleCount {
		objectBuilder := variant_content.NewObjectBuilder().
			WithTypeCity().
			WithGuardChance(1).
			WithGuardValue(guardValue).
			WithGuardWeeklyIncrement(0.15).
			WithCastleQualityMedium().
			WithPlacementUniform().
			WithPlacementArgs("false", "-0.8", "3")
		if matchPlayerFaction {
			objectBuilder = objectBuilder.WithFaction("Match", "0")
		} else {
			objectBuilder = objectBuilder.WithFaction("Random")
		}
		castles = append(castles, objectBuilder.Build())
	}
	return castles
}

// CreateNeutralZoneCastles builds the City main objects of a neutral zone.
// Exported so the manual zone editor can rebuild castles when the user edits
// a zone's quality or castle count.
func CreateNeutralZoneCastles(
	profile models.NeutralZoneProfile,
	tuning models.GenerationTuning,
	castleCount int,
	isHoldCityZone bool) []entities.MainObject {
	var castles []entities.MainObject

	if castleCount > 0 {
		objectBuilder := variant_content.NewObjectBuilder().
			WithGuardChance(1).
			WithGuardWeeklyIncrement(0.10).
			WithFaction("FromList")

		if tuning.SpawnAbandonedOutposts {
			objectBuilder = objectBuilder.WithTypeAbandonedOutpost()
		} else {
			objectBuilder = objectBuilder.WithTypeCity()
		}

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
		objectBuilder := variant_content.NewObjectBuilder().
			WithGuardChance(1).
			WithGuardValue(tuning.ScaleByBorderGuardStrength(profile.ExtraCityGuardValue)).
			WithGuardWeeklyIncrement(0.10).
			WithCastleQuality(profile.ExtraBuildingsCSid).
			WithFaction("FromList").
			WithPlacementUniform().
			WithPlacementArgs("false", "-0.8", "3")

		if tuning.SpawnAbandonedOutposts {
			objectBuilder = objectBuilder.WithTypeAbandonedOutpost()
		} else {
			objectBuilder = objectBuilder.WithTypeCity()
		}

		castles = append(castles, objectBuilder.Build())
	}

	return castles
}

func (this *TopologyBase) createHubZoneCastles(
	tuning models.GenerationTuning,
	castleCount int,
	isHoldCityZone bool) []entities.MainObject {
	var castles []entities.MainObject
	newCastleBuilder := func() *variant_content.MainObjectBuilder {
		return variant_content.NewObjectBuilder().
			WithTypeCity().
			WithGuardWeeklyIncrement(0.10).
			WithFaction("FromList")
	}
	buildHoldCityCastle := func(builder *variant_content.MainObjectBuilder) entities.MainObject {
		return builder.
			WithGuardChance(1).
			WithGuardValue(tuning.ScaleByNeutralGuardStrength(25_000)).
			WithCastleQualityUltraRich().
			WithPlacementCenter().
			WithHoldCityWinCon().
			Build()
	}
	buildCastle := func(builder *variant_content.MainObjectBuilder) entities.MainObject {
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

func (this *TopologyBase) createOuterZoneRoads(
	connectionNames []string,
	castleCount int,
	includeFoothold, generateRoads bool) []entities.Road {
	if !generateRoads {
		return nil
	}

	if castleCount == 0 {
		return this.CreateConnectorZoneRoads(connectionNames, generateRoads)
	}

	var roads []entities.Road
	for i := range castleCount {
		roads = append(roads,
			variant_content.NewRoadBuilder().
				WithStoneType().
				WithFrom(variant_content.NewRefBuilder().BuildMainObjectType("0")).
				WithTo(variant_content.NewRefBuilder().BuildMainObjectType(fmt.Sprintf("%d", i+1))).
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

func buildSideContentLimits() entities.StringList {
	var limits []string
	for a := 1; a <= 5; a++ {
		for b := a + 1; b <= 6; b++ {
			limits = append(limits, fmt.Sprintf("content_limits_side_%d_%d", a, b))
		}
	}
	return limits
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
				candidate := candidates[j]
				if candidate != i && candidate != (i+1)%count && candidate != (i-1+count)%count {
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
