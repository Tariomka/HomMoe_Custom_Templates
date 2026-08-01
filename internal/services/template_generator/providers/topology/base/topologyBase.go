package base

import (
	"fmt"
	"math/rand/v2"
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/geometry"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/graph"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/placement_rule"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
)

type TopologyBase struct {
	ZoneLabelProvider *zones.ZoneLabelProvider
	castleFactory     *zones.CastleFactory
	roadFactory       *zones.RoadFactory
	zoneFactory       *zones.ZoneFactory
}

func NewTopologyBase() TopologyBase {
	return NewTopologyBaseWithCreationServices(zones.NewCreationServices(nil, nil))
}

func NewTopologyBaseWithCreationServices(creationServices *zones.CreationServices) TopologyBase {
	if creationServices == nil {
		creationServices = zones.NewCreationServices(nil, nil)
	}
	return TopologyBase{
		ZoneLabelProvider: zones.NewZoneLabelProvider(),
		castleFactory:     creationServices.CastleFactory,
		roadFactory:       creationServices.RoadFactory,
		zoneFactory:       creationServices.ZoneFactory,
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
	return this.zoneFactory.CreateSpawnZone(
		label,
		playerName,
		connectionNames,
		castleCount,
		matchFactions,
		zoneSize,
		footholdCount,
		generateRoads,
		tuning,
	)
}

func (this *TopologyBase) CreateNeutralZone(
	plan neutral_zone.Plan,
	connectionNames []string,
	zoneSize float64,
	footholdCount int,
	generateRoads bool,
	tuning models.GenerationTuning,
	isHoldCity bool) entities.Zone {
	return this.zoneFactory.CreateNeutralZone(models.NeutralZoneCreation{
		Name:                 constants.NeutralZonePrefix + plan.Label,
		Quality:              plan.Quality,
		Size:                 zoneSize,
		ConnectionNames:      connectionNames,
		MandatoryContentName: "mandatory_content_neutral_" + plan.Label,
		CastleCount:          plan.CastleCount,
		HoldCity:             isHoldCity,
		OutpostCount:         tuning.AbandonedOutpostCount,
		FootholdCount:        footholdCount,
		GuardRandomization:   tuning.GuardRandomization,
		GenerateRoads:        generateRoads,
		Tuning:               tuning,
	})
}

func (this *TopologyBase) CreateHubZone(
	name string,
	connectionNames []string,
	tuning models.GenerationTuning,
	isHoldCity bool,
	size float64,
	castleCount int,
	generateRoads bool,
	mandatoryContentName string) entities.Zone {
	return this.zoneFactory.CreateHubZone(models.HubZoneCreation{
		Name:                 name,
		Size:                 size,
		ConnectionNames:      connectionNames,
		MandatoryContentName: mandatoryContentName,
		CastleCount:          castleCount,
		HoldCity:             isHoldCity,
		GuardRandomization:   0.05,
		GenerateRoads:        generateRoads,
		Tuning:               tuning,
	})
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

	rule := placement_rule.NewPlacementRuleBuilder().BuildNearCrossroadsRule(2)
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

		additionalConns = append(additionalConns, variant_content.NewConnectionBuilder().
			WithName(fallbackName).
			WithFrom(zoneName).
			WithTo(constants.PlayerZonePrefix+partner).
			WithConnectionTypeDirect().
			WithGuardZone(zoneName).
			WithSimTurnSquad().
			WithGuardValue(this.GetBorderGuardValue(label, partner, playerLabels, nil, tuning)).
			WithGuardWeeklyIncrement(0.15).
			WithGuardMatchGroup("fallback_guard_"+fallbackName).
			Build())
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
	return this.roadFactory.CreateConnectorZoneRoads(connectionNames, generateRoads)
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
	return this.castleFactory.CreatePlayerOwnedCastles(matchPlayerFaction, owner, castleCount)
}

// CreatePlayerUnclaimedCastles builds the extra neutral City castles that sit
// inside a player's zone but stay unowned until someone captures them.
// Exported so the manual zone editor can rebuild a spawn zone's castles when
// the player-castle options change after manual editing.
func (this *TopologyBase) CreatePlayerUnclaimedCastles(
	matchPlayerFaction bool,
	guardValue, castleCount int) []entities.MainObject {
	return this.castleFactory.CreatePlayerUnclaimedCastles(matchPlayerFaction, guardValue, castleCount)
}

// CreateHubZoneCastles builds the City main objects of a hub zone. Exported
// so the manual zone editor can rebuild hub castles when the hub-castle
// option changes after manual editing.
func (this *TopologyBase) CreateHubZoneCastles(
	tuning models.GenerationTuning,
	castleCount int,
	isHoldCityZone bool) []entities.MainObject {
	return this.castleFactory.CreateHubZoneCastles(tuning, castleCount, isHoldCityZone)
}

func (this *TopologyBase) CreateOuterZoneRoads(
	connectionNames []string,
	mainObjectCount int,
	footholdCount int, generateRoads bool) []entities.Road {
	return this.roadFactory.CreateOuterZoneRoads(
		connectionNames,
		mainObjectCount,
		footholdCount,
		generateRoads,
	)
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

// CreateNeutralZoneCastles builds the City main objects of a neutral zone.
// Exported so the manual zone editor can rebuild castles when the user edits
// a zone's quality or castle count.
func CreateNeutralZoneCastles(
	profile neutral_zone.Profile,
	tuning models.GenerationTuning,
	castleCount int,
	isHoldCityZone bool) []entities.MainObject {
	return zones.NewCastleFactory().CreateNeutralZoneCastles(
		profile,
		tuning,
		castleCount,
		isHoldCityZone,
	)
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
