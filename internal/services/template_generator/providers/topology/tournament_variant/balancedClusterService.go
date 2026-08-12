package tournament_variant

import (
	"fmt"
	"math"
	"slices"
	"sort"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_connections"
	"github.com/Tariomka/hommoe_custom_templates/internal/common/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/geometry_helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/position_layout"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/tournament_variant/misc"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones/zone_interfaces"
)

type BalancedClusterService struct {
	base.TopologyBase

	positionLayoutService *position_layout.PositionLayoutService
}

func NewBalancedClusterService(
	zoneFactory zone_interfaces.IZoneFactory,
	roadFactory zone_interfaces.IRoadFactory,
	zoneLabelProvider zone_interfaces.IZoneLabelProvider,
	connectionService base.ITopologyConnectionService,
) *BalancedClusterService {
	return &BalancedClusterService{
		TopologyBase:          base.NewTopologyBase(zoneFactory, roadFactory, zoneLabelProvider, connectionService),
		positionLayoutService: position_layout.NewPositionLayoutService(),
	}
}

func (this *BalancedClusterService) CreateClusterVariant(
	configuration config.GeneratorConfig,
	tuning models.GenerationTuning,
	allNeutralZonePlans, playerNeutralZonePlans neutral_zone.Plans,
	playerIndex int,
	playerLabel string) ([]entities.Zone, []entities.Connection) {
	singlePlayerList := []string{playerLabel}
	orderedLabels := this.ZoneLabelProvider.CreateBalancedRingZoneLabels(singlePlayerList, playerNeutralZonePlans)
	rawPositions := this.positionLayoutService.CreatePositionsFromPlans(
		orderedLabels,
		singlePlayerList,
		allNeutralZonePlans,
	)
	positions := this.createPositions(rawPositions, playerIndex)

	sortedPairs := this.createSortedPairs(orderedLabels, rawPositions, allNeutralZonePlans)

	connectionNames := this.createConnectionNames(orderedLabels, sortedPairs)
	zones := this.createZones(
		configuration,
		playerLabel,
		playerIndex,
		orderedLabels,
		tuning,
		allNeutralZonePlans,
		connectionNames,
	)

	for index, label := range orderedLabels {
		position := positions[index]
		zones[index].GeneratorPosition = &[2]float64{position.X, position.Y}
		// Seed the editor/preview layout with the mirrored position too, so the
		// generated tournament opens as a full-size, two-sided mirror (honored
		// by layoutManualPositions for every topology) while staying fully
		// editable. ManualPosition is editor-only (json:"-") and is never saved.
		zones[index].ManualPosition = &[2]float64{position.X, position.Y}
		tier := allNeutralZonePlans.GetTier(label)
		zones[index].GeneratorRing = &tier
	}

	connections := this.createConnections(
		playerLabel, orderedLabels,
		tuning, allNeutralZonePlans,
		connectionNames, sortedPairs)
	connections = append(
		connections,
		this.CreateMissingConnections(
			singlePlayerList, orderedLabels,
			positions, zones, connections,
			tuning, allNeutralZonePlans)...)
	return zones, connections
}

func (this *BalancedClusterService) createPositions(rawPositions models.Positions, playerIndex int) models.Positions {
	if len(rawPositions) < 1 {
		return models.Positions{}
	}

	minimumPosition, maximumPosition := geometry_helpers.GetPositionBounds(rawPositions)
	spanX := math.Max(maximumPosition.X-minimumPosition.X, 0.001)
	spanY := math.Max(maximumPosition.Y-minimumPosition.Y, 0.001)

	// Player 0 fills the left half, player 1 the right half. Player 1 is laid
	// out as a true mirror of player 0 - reflected on both the horizontal and
	// vertical axes - so the generated tournament starts as a symmetric,
	// mirrored layout for the two players.
	xMin, xMax := 0.03, 0.43
	xSign, ySign := 1.0, 1.0
	if playerIndex != 0 {
		xMin, xMax = 0.57, 0.97
		xSign, ySign = -1.0, -1.0
	}

	scale := math.Min((xMax-xMin)/spanX, 0.9/spanY)
	xCenter := (xMin + xMax) / 2.0
	yCenter := 0.5

	positions := models.Positions{}
	for _, position := range rawPositions {
		positions.Add(data.NewVec2(
			xCenter+xSign*(position.X-(minimumPosition.X+maximumPosition.X)/2.0)*scale,
			yCenter+ySign*(position.Y-(minimumPosition.Y+maximumPosition.Y)/2.0)*scale))
	}

	return positions
}

func (this *BalancedClusterService) createSortedPairs(
	orderedLabels []string,
	rawPositions models.Positions,
	allNeutralZonePlans neutral_zone.Plans) [][2]int {
	tierIndices := bucketIndicesByTier(orderedLabels, allNeutralZonePlans)
	tierKeys := make([]int, 0, len(tierIndices))
	for tier := range tierIndices {
		tierKeys = append(tierKeys, tier)
	}
	slices.Sort(tierKeys)

	tierSorted, tierAngles := sortTiersByAngle(tierIndices, rawPositions)

	pairSet := misc.NewPairSet()

	// Same-ring: circle-neighbors only; skip degenerate < 3 rings.
	addSameRingPairs(pairSet.Add, tierSorted)

	// Cross-ring: bidirectional nearest-neighbor between adjacent tiers.
	for tierIndex := range tierKeys[:len(tierKeys)-1] {
		addCrossRingPairs(
			pairSet.Add,
			tierSorted[tierKeys[tierIndex]], tierAngles[tierKeys[tierIndex]],
			tierSorted[tierKeys[tierIndex+1]], tierAngles[tierKeys[tierIndex+1]])
	}

	sortedPairs := linq.FromMap(*pairSet).SelectKeys().ToSlice()
	sort.Slice(sortedPairs, func(i, j int) bool {
		if sortedPairs[i][0] != sortedPairs[j][0] {
			return sortedPairs[i][0] < sortedPairs[j][0]
		}
		return sortedPairs[i][1] < sortedPairs[j][1]
	})
	return sortedPairs
}

// bucketIndicesByTier groups the ordered label indexes by their neutral-plan tier.
func bucketIndicesByTier(orderedLabels []string, allNeutralZonePlans neutral_zone.Plans) map[int][]int {
	tierIndices := map[int][]int{}
	for index, label := range orderedLabels {
		tier := allNeutralZonePlans.GetTier(label)
		tierIndices[tier] = append(tierIndices[tier], index)
	}
	return tierIndices
}

// sortTiersByAngle orders every tier's indexes by their angle around the map
// center and returns the sorted indexes alongside the matching angles.
func sortTiersByAngle(
	tierIndices map[int][]int,
	rawPositions models.Positions) (map[int][]int, map[int][]float64) {
	tierSorted := map[int][]int{}
	tierAngles := map[int][]float64{}
	for tier, indexes := range tierIndices {
		sorted := make([]int, len(indexes))
		copy(sorted, indexes)
		sort.SliceStable(sorted, func(i, j int) bool {
			return math.Atan2(rawPositions[sorted[i]].Y-0.5, rawPositions[sorted[i]].X-0.5) <
				math.Atan2(rawPositions[sorted[j]].Y-0.5, rawPositions[sorted[j]].X-0.5)
		})
		tierSorted[tier] = sorted
		angles := make([]float64, len(sorted))
		for j, zoneIndex := range sorted {
			angles[j] = math.Atan2(rawPositions[zoneIndex].Y-0.5, rawPositions[zoneIndex].X-0.5)
		}
		tierAngles[tier] = angles
	}
	return tierSorted, tierAngles
}

// addSameRingPairs links each ring's circle-neighbours; rings with fewer than
// three zones are skipped as degenerate.
func addSameRingPairs(addPair func(a, b int), tierSorted map[int][]int) {
	for _, sorted := range tierSorted {
		ringSize := len(sorted)
		if ringSize < 3 {
			continue
		}
		for i := range sorted {
			addPair(sorted[i], sorted[(i+1)%ringSize])
		}
	}
}

// addCrossRingPairs links two adjacent rings bidirectionally: every outer zone
// to its angularly nearest inner zone, and every inner zone to all of its
// angularly nearest outer zones (ties within epsilon all count).
func addCrossRingPairs(
	addPair func(a, b int),
	outerSorted []int, outerAngles []float64,
	innerSorted []int, innerAngles []float64) {
	for outerIndex := range outerSorted {
		best, bestD := 0, math.MaxFloat64
		for innerIndex := range innerSorted {
			if distance := misc.GetShortestAngleDistance(
				outerAngles[outerIndex],
				innerAngles[innerIndex]); distance < bestD {
				bestD = distance
				best = innerIndex
			}
		}
		if len(innerSorted) > 0 {
			addPair(outerSorted[outerIndex], innerSorted[best])
		}
	}

	const epsilon = 1e-9
	for innerIndex := range innerSorted {
		bestDistance := math.MaxFloat64
		for outerIndex := range outerSorted {
			if distance := misc.GetShortestAngleDistance(
				innerAngles[innerIndex],
				outerAngles[outerIndex],
			); distance < bestDistance {
				bestDistance = distance
			}
		}
		for outerIndex := range outerSorted {
			if misc.GetShortestAngleDistance(
				innerAngles[innerIndex],
				outerAngles[outerIndex]) <= bestDistance+epsilon {
				addPair(innerSorted[innerIndex], outerSorted[outerIndex])
			}
		}
	}
}

func (this *BalancedClusterService) createConnectionNames(orderedLabels []string, sortedPairs [][2]int) [][]string {
	connectionNamesByZone := make([][]string, len(orderedLabels))
	for _, pair := range sortedPairs {
		labelFrom := orderedLabels[pair[0]]
		labelTo := orderedLabels[pair[1]]
		connectionName := fmt.Sprintf("TBal-%s-%s", labelFrom, labelTo)
		connectionNamesByZone[pair[0]] = append(connectionNamesByZone[pair[0]], connectionName)
		connectionNamesByZone[pair[1]] = append(connectionNamesByZone[pair[1]], connectionName)
	}
	return connectionNamesByZone
}

func (this *BalancedClusterService) createZones(
	configuration config.GeneratorConfig,
	playerLabel string,
	playerIndex int,
	orderedLabels []string,
	tuning models.GenerationTuning,
	allNeutralZonePlans neutral_zone.Plans,
	connectionNames [][]string) []entities.Zone {
	var zones []entities.Zone
	for index, label := range orderedLabels {
		myConns := connectionNames[index]
		zones = append(zones, this.CreateClusterZone(
			configuration, label, myConns, playerIndex, label == playerLabel, false, tuning, allNeutralZonePlans))
	}
	return zones
}

func (this *BalancedClusterService) createConnections(
	playerLabel string,
	orderedLabels []string,
	tuning models.GenerationTuning,
	allNeutralZonePlans neutral_zone.Plans,
	connectionNames [][]string,
	sortedPairs [][2]int) []entities.Connection {
	nameLookup := make(map[int]int, len(orderedLabels))

	var connections []entities.Connection
	for _, pair := range sortedPairs {
		indexA, indexB := pair[0], pair[1]
		labelFrom := orderedLabels[indexA]
		labelTo := orderedLabels[indexB]

		connName := connectionNames[indexA][nameLookup[indexA]]
		nameLookup[indexA]++
		nameLookup[indexB]++

		fromZone := constants.PlayerZonePrefix + labelFrom
		if labelFrom != playerLabel {
			fromZone = constants.NeutralZonePrefix + labelFrom
		}
		toZone := constants.PlayerZonePrefix + labelTo
		if labelTo != playerLabel {
			toZone = constants.NeutralZonePrefix + labelTo
		}
		connections = append(connections, variant_content.NewConnectionBuilder().
			WithName(connName).
			WithFrom(fromZone).
			WithTo(toZone).
			WithConnectionTypeDirect().
			WithGuardZone(fromZone).
			WithSimTurnSquad().
			WithGuardValue(this.GetBorderGuardValue(
				labelFrom, labelTo, []string{playerLabel}, allNeutralZonePlans, tuning)).
			WithGuardWeeklyIncrement(common_connections.GetGuardWeeklyIncrements().Standard).
			WithGuardMatchGroup(fmt.Sprintf("tourney_bal_guard_%s_%s", labelFrom, labelTo)).
			Build())
	}
	return connections
}
