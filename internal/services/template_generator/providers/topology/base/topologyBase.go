package base

import (
	"fmt"
	"math"
	"math/rand/v2"
	"slices"
	"strconv"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_zones"
	"github.com/Tariomka/hommoe_custom_templates/internal/common/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/geometry"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/graph"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/placement_rule"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base/utils"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
)

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
		orientationBuilder.WithZeroAngleZone(constants.PlayerZonePrefix + firstLabel)
	} else {
		orientationBuilder.WithZeroAngleZone(constants.NeutralZonePrefix + firstLabel)
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
	footholdCount int,
	generateRoads bool,
	tuning models.GenerationTuning) entities.Zone {
	mainObjects := []entities.MainObject{
		this.createPlayerSpawnCastle(playerName, tuning.ScaleByNeutralGuardStrength(5000)),
	}
	mainObjects = append(mainObjects,
		this.CreatePlayerOwnedCastles(matchFactions, playerName, tuning.PlayerOwnedCastles)...)
	mainObjects = append(mainObjects,
		this.CreatePlayerUnclaimedCastles(matchFactions, tuning.ScaleByNeutralGuardStrength(5000), castleCount)...)

	// Roads connect the spawn castle (main object 0) to every other castle in
	// the zone; player-owned extras are road-linked just like unclaimed ones.
	// A zone with no extra castles at all stays a pass-through connector, so
	// pass 0 main objects in that case; otherwise pass the full main-object
	// count so every extra castle gets a road.
	roadCastleCount := castleCount + tuning.PlayerOwnedCastles
	roadMainObjectCount := 0
	if roadCastleCount > 0 {
		roadMainObjectCount = len(mainObjects)
	}

	return variant_content.NewZoneBuilder().
		WithName(constants.PlayerZonePrefix + label).
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
		WithResourcesContentPool([]string{registry.GetResourcesContentPoolValues().StartZonePoor}).
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
		WithRoads(this.CreateOuterZoneRoads(connectionNames, roadMainObjectCount, footholdCount, generateRoads)).
		Build()
}

func (this *TopologyBase) CreateNeutralZone(
	plan neutral_zone.Plan,
	connectionNames []string,
	zoneSize float64,
	footholdCount int,
	generateRoads bool,
	tuning models.GenerationTuning,
	isHoldCity bool) entities.Zone {
	if isHoldCity && plan.CastleCount < 1 {
		plan.CastleCount = 1
	}
	profile := common_zones.GetNeutralZoneProfile(plan.Quality)

	// Abandoned outposts are spawned in addition to the zone's castles, with
	// their own count slider instead of being tied to the castle count.
	mainObjects := CreateNeutralZoneCastles(profile, tuning, plan.CastleCount, isHoldCity)
	mainObjects = append(mainObjects, createAbandonedOutposts(profile, tuning, tuning.AbandonedOutpostCount)...)
	totalMainObjects := plan.CastleCount + tuning.AbandonedOutpostCount

	zoneBuilder := variant_content.NewZoneBuilder().
		WithName(constants.NeutralZonePrefix + plan.Label).
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
		WithMainObjects(mainObjects).
		WithCrossroadsPosition(0).
		WithRoads(this.CreateOuterZoneRoads(connectionNames, totalMainObjects, footholdCount, generateRoads))

	if totalMainObjects > 0 {
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
	generateRoads bool,
	mandatoryContentName string) entities.Zone {
	if isHoldCity && castleCount < 1 {
		castleCount = 1
	}
	profile := common_zones.GetNeutralZoneProfile(neutral_zone.QualityHighest)

	zoneBuilder := variant_content.NewZoneBuilder().
		WithName(constants.HubZoneName).
		WithSize(utils.NormalizeZoneSize(size)).
		WithLayout(profile.Layout).
		WithGuardCutoffValue(2000).
		WithGuardRandomization(0.05).
		WithGuardMultiplier(tuning.ScaleByNeutralGuardStrengthPrecise(profile.GuardMultiplier)).
		WithGuardWeeklyIncrement(0.20).
		WithGuardReactionDistribution(profile.GuardReactionDistribution).
		WithDiplomacyModifier(-0.5).
		WithGuardedContentPool(profile.GuardedContentPool).
		WithUnguardedContentPool(profile.UnguardedContentPool).
		WithResourcesContentPool(profile.ResourcesContentPool).
		WithContentCountLimits(buildSideContentLimits()).
		WithGuardedContentValue(tuning.ScaleByStructureDensity(float64(profile.GuardedContentValue) * tuning.ContentScale)).
		WithGuardedContentValuePerArea(tuning.ScaleByStructureDensity(float64(profile.GuardedContentValuePerArea) * math.Sqrt(tuning.ContentScale))).
		WithUnguardedContentValue(tuning.ScaleByStructureDensity(float64(profile.UnguardedContentValue) * tuning.ContentScale)).
		WithUnguardedContentValuePerArea(tuning.ScaleByStructureDensity(float64(profile.UnguardedContentValuePerArea) * math.Sqrt(tuning.ContentScale))).
		WithResourcesValue(tuning.ScaleByResourceDensity(float64(profile.ResourcesValue) * tuning.ContentScale)).
		WithResourcesValuePerArea(tuning.ScaleByResourceDensity(float64(profile.ResourcesValuePerArea) * math.Sqrt(tuning.ContentScale))).
		WithMainObjects(this.CreateHubZoneCastles(tuning, castleCount, isHoldCity)).
		WithCrossroadsPosition(0).
		WithRoads(this.CreateOuterZoneRoads(connectionNames, castleCount, 0, generateRoads))

	if mandatoryContentName != "" {
		zoneBuilder = zoneBuilder.WithMandatoryContent(mandatoryContentName)
	}

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
		zoneName := constants.PlayerZonePrefix + label
		zone, ok := linq.FromSlice(zones).First(func(z entities.Zone) bool { return z.Name == zoneName })
		if !ok {
			continue
		}

		if spawnZoneHasConnection(zone, connNames) {
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
			Name: fallbackName, From: zoneName, To: constants.PlayerZonePrefix + partner,
			ConnectionType: "Direct", GuardZone: zoneName, SimTurnSquad: true,
			GuardValue: this.GetBorderGuardValue(label, partner, playerLabels, nil, tuning), GuardWeeklyIncrement: 0.15,
			GuardMatchGroup: "fallback_guard_" + fallbackName,
		})
		connNames[fallbackName] = true
		appendSpawnFallbackRoads(zones, label, partner, fallbackName)
	}
	return additionalConns
}

// spawnZoneHasConnection reports whether the spawn zone already has a road
// leading to any of the known connections.
func spawnZoneHasConnection(zone entities.Zone, connNames map[string]bool) bool {
	connectionType := registry.GetRoadConnectionTypeValues().Connection
	for _, road := range zone.Roads {
		if road.To.Type == connectionType && len(road.To.Args) > 0 && connNames[road.To.Args[0]] {
			return true
		}
	}
	return false
}

// appendSpawnFallbackRoads adds a road from each of the two spawn zones' first
// main object to the freshly created fallback connection.
func appendSpawnFallbackRoads(zones []entities.Zone, label, partner, fallbackName string) {
	for _, playerLabel := range []string{label, partner} {
		spawnZoneName := constants.PlayerZonePrefix + playerLabel
		zoneIndex := slices.IndexFunc(
			zones,
			func(candidate entities.Zone) bool { return candidate.Name == spawnZoneName })
		if zoneIndex >= 0 {
			zones[zoneIndex].Roads = append(zones[zoneIndex].Roads, variant_content.NewRoadBuilder().
				WithFrom(variant_content.NewRefBuilder().BuildMainObjectType("0")).
				WithTo(variant_content.NewRefBuilder().BuildConnectionType(fallbackName)).
				Build())
		}
	}
}

func (this *TopologyBase) CreateMissingConnections(
	playerLabels, allLabels []string,
	positions models.Positions,
	zones []entities.Zone,
	connections []entities.Connection,
	tuning models.GenerationTuning,
	neutralZones neutral_zone.Plans) []entities.Connection {
	if len(allLabels) < 2 {
		return nil
	}

	adjacency := this.buildZoneAdjacency(playerLabels, allLabels, connections)

	connNameSet := map[string]bool{}
	for connection := range linq.FromSlice(connections).
		Where(func(x entities.Connection) bool { return x.Name != "" }).Iterate {
		connNameSet[connection.Name] = true
	}

	var additionalConns []entities.Connection
	for {
		nodes := make([]int, len(allLabels))
		for index := range nodes {
			nodes[index] = index
		}
		components := graph.ConnectedComponents(adjacency, nodes)
		bestIndexes, ok := geometry.FindClosestAcrossComponents(positions, components)
		if !ok {
			break
		}

		labelA, labelB := allLabels[bestIndexes.X], allLabels[bestIndexes.Y]
		if labelA > labelB {
			labelA, labelB = labelB, labelA
		}
		bridgeName := fmt.Sprintf("Bridge-%s-%s", labelA, labelB)
		if connNameSet[bridgeName] {
			// The existing bridge already connects the two components; link them so the
			// loop progresses instead of reselecting the same pair forever.
			graph.Link(adjacency, bestIndexes.X, bestIndexes.Y)
			continue
		}

		zoneFrom := this.ZoneLabelProvider.CreateZoneName(allLabels[bestIndexes.X], playerLabels)
		zoneTo := this.ZoneLabelProvider.CreateZoneName(allLabels[bestIndexes.Y], playerLabels)
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

		appendBridgeRoads(zones, zoneFrom, zoneTo, bridgeName)

		graph.Link(adjacency, bestIndexes.X, bestIndexes.Y)
	}

	return additionalConns
}

// findExistingConnName returns the name of the first connection referenced by
// the zone's roads, or "" when the zone has no connection-bound road.
func findExistingConnName(zone entities.Zone) string {
	connectionType := registry.GetRoadConnectionTypeValues().Connection
	for _, road := range zone.Roads {
		if road.From.Type == connectionType && len(road.From.Args) > 0 {
			return road.From.Args[0]
		}
		if road.To.Type == connectionType && len(road.To.Args) > 0 {
			return road.To.Args[0]
		}
	}
	return ""
}

// appendBridgeRoads wires the new bridge connection into both endpoint zones:
// from the zone's first main object when it has one, otherwise chained off an
// existing connection road (or the bridge itself as a last resort).
func appendBridgeRoads(zones []entities.Zone, zoneFrom, zoneTo, bridgeName string) {
	for _, zoneName := range []string{zoneFrom, zoneTo} {
		zoneIndex := slices.IndexFunc(
			zones,
			func(candidate entities.Zone) bool { return candidate.Name == zoneName })
		if zoneIndex < 0 {
			continue
		}
		roadBuilder := variant_content.NewRoadBuilder().WithTo(
			variant_content.NewRefBuilder().BuildConnectionType(bridgeName))
		switch {
		case len(zones[zoneIndex].MainObjects) > 0:
			zones[zoneIndex].Roads = append(zones[zoneIndex].Roads,
				roadBuilder.WithFrom(variant_content.NewRefBuilder().BuildMainObjectType("0")).Build())
		case len(zones[zoneIndex].Roads) > 0:
			connectionName := findExistingConnName(zones[zoneIndex])
			if connectionName == "" {
				connectionName = bridgeName
			}
			zones[zoneIndex].Roads = append(zones[zoneIndex].Roads,
				roadBuilder.WithFrom(variant_content.NewRefBuilder().BuildConnectionType(connectionName)).
					Build())
		default:
			zones[zoneIndex].Roads = append(zones[zoneIndex].Roads,
				roadBuilder.WithFrom(variant_content.NewRefBuilder().BuildConnectionType(bridgeName)).Build())
		}
	}
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
	neutralZones neutral_zone.Plans,
	tuning models.GenerationTuning) int {
	aIsPlayer := slices.Contains(playerLabels, labelA)
	bIsPlayer := slices.Contains(playerLabels, labelB)
	if aIsPlayer && bIsPlayer {
		return tuning.ScaleByBorderGuardStrength(neutral_zone.QualityUnknown.GetGuardValue())
	}

	if aIsPlayer || bIsPlayer {
		neutralLabel := labelB
		if !aIsPlayer {
			neutralLabel = labelA
		}
		return tuning.ScaleByBorderGuardStrength(neutralZones.GetQuality(neutralLabel).GetGuardValue())
	}

	higher := max(neutralZones.GetQuality(labelA), neutralZones.GetQuality(labelB))
	return tuning.ScaleByBorderGuardStrength(higher.GetGuardValue())
}

// CreatePlayerOwnedCastles builds the extra City castles that the player owns
// from the very start of the game. Because they already have an owner, their
// guards are dropped immediately so the player can use them right away.
// Exported so the manual zone editor can rebuild a spawn zone's castles when
// the player-castle options change after manual editing.
func (this *TopologyBase) CreatePlayerOwnedCastles(
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
			objectBuilder = objectBuilder.WithFactionMatch()
		} else {
			objectBuilder = objectBuilder.WithFaction("Random")
		}
		castles = append(castles, objectBuilder.Build())
	}
	return castles
}

// CreatePlayerUnclaimedCastles builds the extra neutral City castles that sit
// inside a player's zone but stay unowned until someone captures them.
// Exported so the manual zone editor can rebuild a spawn zone's castles when
// the player-castle options change after manual editing.
func (this *TopologyBase) CreatePlayerUnclaimedCastles(
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
			objectBuilder = objectBuilder.WithFactionMatch()
		} else {
			objectBuilder = objectBuilder.WithFaction("Random")
		}
		castles = append(castles, objectBuilder.Build())
	}
	return castles
}

// CreateHubZoneCastles builds the City main objects of a hub zone. Exported
// so the manual zone editor can rebuild hub castles when the hub-castle
// option changes after manual editing.
func (this *TopologyBase) CreateHubZoneCastles(
	tuning models.GenerationTuning,
	castleCount int,
	isHoldCityZone bool) []entities.MainObject {
	var castles []entities.MainObject
	newCastleBuilder := func() *variant_content.MainObjectBuilder {
		return variant_content.NewObjectBuilder().
			WithTypeCity().
			WithGuardWeeklyIncrement(0.10).
			WithFactionFromList()
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

func (this *TopologyBase) CreateOuterZoneRoads(
	connectionNames []string,
	mainObjectCount int,
	footholdCount int, generateRoads bool) []entities.Road {
	if !generateRoads {
		return nil
	}

	if mainObjectCount == 0 {
		return this.CreateConnectorZoneRoads(connectionNames, generateRoads)
	}

	var roads []entities.Road
	for i := range mainObjectCount - 1 {
		roads = append(roads,
			variant_content.NewRoadBuilder().
				WithStoneType().
				WithFrom(variant_content.NewRefBuilder().BuildMainObjectType("0")).
				WithTo(variant_content.NewRefBuilder().BuildMainObjectType(strconv.Itoa(i+1))).
				Build())
	}
	for i := range footholdCount {
		roads = append(roads,
			variant_content.NewRoadBuilder().
				WithFrom(variant_content.NewRefBuilder().BuildMainObjectType("0")).
				WithTo(variant_content.NewRefBuilder().
					BuildMandatoryContentType(fmt.Sprintf("name_remote_foothold_%d", i+1))).
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

// buildZoneAdjacency indexes every zone label and links the indexes of zones
// already joined by a Direct or Portal connection.
func (this *TopologyBase) buildZoneAdjacency(
	playerLabels, allLabels []string,
	connections []entities.Connection) graph.Adjacency[int] {
	nodes := make([]int, len(allLabels))
	for index := range nodes {
		nodes[index] = index
	}
	adjacency := graph.NewAdjacency(nodes)
	zoneNameToIdx := map[string]int{}
	for index, label := range allLabels {
		zoneNameToIdx[this.ZoneLabelProvider.CreateZoneName(label, playerLabels)] = index
	}
	for connection := range linq.FromSlice(connections).
		Where(func(x entities.Connection) bool {
			connectionTypes := registry.GetConnectionTypeValues()
			return x.ConnectionType == connectionTypes.Direct || x.ConnectionType == connectionTypes.Portal
		}).
		Where(func(x entities.Connection) bool {
			_, okA := zoneNameToIdx[x.From]
			_, okB := zoneNameToIdx[x.To]
			return okA && okB
		}).Iterate {
		graph.Link(adjacency, zoneNameToIdx[connection.From], zoneNameToIdx[connection.To])
	}
	return adjacency
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

// CreateNeutralZoneCastles builds the City main objects of a neutral zone.
// Exported so the manual zone editor can rebuild castles when the user edits
// a zone's quality or castle count.
func CreateNeutralZoneCastles(
	profile neutral_zone.Profile,
	tuning models.GenerationTuning,
	castleCount int,
	isHoldCityZone bool) []entities.MainObject {
	var castles []entities.MainObject

	if castleCount > 0 {
		objectBuilder := variant_content.NewObjectBuilder().
			WithTypeCity().
			WithGuardChance(1).
			WithGuardWeeklyIncrement(0.10).
			WithFactionFromList()

		if isHoldCityZone {
			objectBuilder = objectBuilder.
				WithGuardValue(tuning.ScaleByBorderGuardStrength(max(profile.PrimaryCityGuardValue, 20_000))).
				WithCastleQualityUltraRich().
				WithPlacementCenter().
				WithHoldCityWinCon()
		} else {
			objectBuilder = objectBuilder.
				WithGuardValue(tuning.ScaleByBorderGuardStrength(profile.PrimaryCityGuardValue)).
				WithCastleQuality(profile.PrimaryBuildingsSid).
				WithPlacementUniform().
				WithPlacementArgs("true", "0.8", "2")
		}

		castles = append(castles, objectBuilder.Build())
	}

	for range castleCount - 1 {
		objectBuilder := variant_content.NewObjectBuilder().
			WithTypeCity().
			WithGuardChance(1).
			WithGuardValue(tuning.ScaleByBorderGuardStrength(profile.ExtraCityGuardValue)).
			WithGuardWeeklyIncrement(0.10).
			WithCastleQuality(profile.ExtraBuildingsSid).
			WithFactionFromList().
			WithPlacementUniform().
			WithPlacementArgs("false", "-0.8", "3")

		castles = append(castles, objectBuilder.Build())
	}

	return castles
}

// createAbandonedOutposts builds extra AbandonedOutpost main objects that sit
// in a neutral zone alongside its City castles. The number of outposts is
// driven by the dedicated count rather than the zone's castle count.
func createAbandonedOutposts(
	profile neutral_zone.Profile,
	tuning models.GenerationTuning,
	count int) []entities.MainObject {
	var outposts []entities.MainObject
	for range count {
		outposts = append(outposts,
			variant_content.NewObjectBuilder().
				WithTypeAbandonedOutpost().
				WithGuardChance(1).
				WithGuardValue(tuning.ScaleByBorderGuardStrength(profile.ExtraCityGuardValue)).
				WithGuardWeeklyIncrement(0.10).
				WithCastleQuality(profile.ExtraBuildingsSid).
				WithFactionFromList().
				WithPlacementUniform().
				WithPlacementArgs("false", "-0.8", "3").
				Build())
	}
	return outposts
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
	for range 100 {
		if dest, ok := tryRandomDerangement(count); ok {
			return dest
		}
	}
	return buildShiftDerangement(count)
}

// tryRandomDerangement attempts to assign every index a random destination that
// is neither itself nor one of its ring neighbours; it reports failure when the
// shuffled candidate order leaves some index without a legal destination.
func tryRandomDerangement(count int) ([]int, bool) {
	candidates := make([]int, count)
	for i := range candidates {
		candidates[i] = i
	}
	rand.Shuffle(len(candidates), func(i, j int) { candidates[i], candidates[j] = candidates[j], candidates[i] })

	dest := make([]int, count)
	used := make([]bool, count)
	for i := range count {
		found := pickDerangementTarget(candidates, used, i, count)
		if found < 0 {
			return nil, false
		}
		dest[i] = candidates[found]
		used[candidates[found]] = true
	}
	return dest, true
}

// pickDerangementTarget returns the position (within candidates) of the first
// unused candidate that is neither i nor a ring neighbor of i; when none
// exists it falls back to any unused candidate other than i, or -1.
func pickDerangementTarget(candidates []int, used []bool, i, count int) int {
	for j := range candidates {
		if used[candidates[j]] {
			continue
		}
		candidate := candidates[j]
		if candidate != i && candidate != (i+1)%count && candidate != (i-1+count)%count {
			return j
		}
	}
	for j := range candidates {
		if !used[candidates[j]] && candidates[j] != i {
			return j
		}
	}
	return -1
}

// buildShiftDerangement is the deterministic fallback: every index maps to the
// one half a ring away, which is never itself or an adjacent index for the
// zone counts this is called with.
func buildShiftDerangement(count int) []int {
	dest := make([]int, count)
	shift := max(1, count/2)
	for i := range count {
		dest[i] = (i + shift) % count
	}
	return dest
}
