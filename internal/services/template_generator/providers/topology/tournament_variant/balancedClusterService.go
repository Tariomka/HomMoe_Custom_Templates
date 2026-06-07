package tournament_variant

import (
	"fmt"
	"math"
	"slices"
	"sort"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/tournament_variant/misc"
)

type BalancedClusterService struct {
	base.TopologyBase
}

func NewBalancedClusterService() *BalancedClusterService {
	return &BalancedClusterService{
		TopologyBase: base.NewTopologyBase(),
	}
}

func (this *BalancedClusterService) CreateClusterVariant(
	configuration config.GeneratorConfig,
	tuning models.GenerationTuning,
	allNeutralZonePlans, playerNeutralZonePlans models.NeutralZonePlans,
	playerIndex int,
	playerLabel string) ([]template.Zone, []template.Connection) {
	singlePlayerList := []string{playerLabel}
	orderedLabels := this.ZoneLabelProvider.CreateBalancedRingZoneLabels(singlePlayerList, playerNeutralZonePlans, 0)
	rawPositions := models.CreatePositionsFromPlans(orderedLabels, singlePlayerList, allNeutralZonePlans)
	positions := this.createPositions(rawPositions, playerIndex)

	sortedPairs := this.createSortedPairs(orderedLabels, rawPositions, allNeutralZonePlans)

	connectionNames := this.createConnectionNames(orderedLabels, sortedPairs)
	zones := this.createZones(configuration, playerLabel, playerIndex, orderedLabels, tuning, allNeutralZonePlans, connectionNames)

	clusterStart := len(zones)
	for index, label := range orderedLabels {
		position := positions[index]
		zones[clusterStart+index].GeneratorPosition = &[2]float64{position.X, position.Y}
		tier := allNeutralZonePlans.GetTier(label)
		zones[clusterStart+index].GeneratorRing = &tier
	}

	connections := this.createConnections(playerLabel, orderedLabels, tuning, allNeutralZonePlans, connectionNames, sortedPairs)
	connections = this.CreateMissingConnections(singlePlayerList, orderedLabels, positions, zones[clusterStart:], connections, tuning, allNeutralZonePlans)
	return zones, connections
}

func (this *BalancedClusterService) createPositions(rawPositions models.Positions, playerIndex int) models.Positions {
	if len(rawPositions) < 1 {
		return models.Positions{}
	}

	min, max := rawPositions.GetMinAndMaxPositions()
	spanX := math.Max(max.X-min.X, 0.001)
	spanY := math.Max(max.Y-min.Y, 0.001)

	xMin, xMax := 0.03, 0.43
	if playerIndex != 0 {
		xMin, xMax = 0.57, 0.97
	}

	scale := math.Min((xMax-xMin)/spanX, 0.9/spanY)
	xCentre := (xMin + xMax) / 2.0
	yCentre := 0.5

	positions := models.Positions{}
	for _, position := range rawPositions {
		positions.Add(models.NewPosition(
			xCentre+(position.X-(min.X+max.X)/2.0)*scale,
			yCentre+(position.Y-(min.Y+max.Y)/2.0)*scale))
	}

	return positions
}

func (this *BalancedClusterService) createSortedPairs(
	orderedLabels []string,
	rawPositions models.Positions,
	allNeutralZonePlans models.NeutralZonePlans) [][2]int {
	// Build connections from pure ring structure (a040c98).
	angDist := func(a, b float64) float64 {
		d := math.Mod(math.Abs(a-b), 2*math.Pi)
		if d > math.Pi {
			d = 2*math.Pi - d
		}
		return d
	}
	tierIndices := map[int][]int{}
	for index, label := range orderedLabels {
		tier := allNeutralZonePlans.GetTier(label)
		tierIndices[tier] = append(tierIndices[tier], index)
	}
	tierKeys := make([]int, 0, len(tierIndices))
	for tier := range tierIndices {
		tierKeys = append(tierKeys, tier)
	}
	slices.Sort(tierKeys)

	tierSorted := map[int][]int{}
	tierAngles := map[int][]float64{}
	for tier, idx := range tierIndices {
		s := make([]int, len(idx))
		copy(s, idx)
		sort.SliceStable(s, func(i, j int) bool {
			return math.Atan2(rawPositions[s[i]].Y-0.5, rawPositions[s[i]].X-0.5) <
				math.Atan2(rawPositions[s[j]].Y-0.5, rawPositions[s[j]].X-0.5)
		})
		tierSorted[tier] = s
		ang := make([]float64, len(s))
		for j, ii := range s {
			ang[j] = math.Atan2(rawPositions[ii].Y-0.5, rawPositions[ii].X-0.5)
		}
		tierAngles[tier] = ang
	}

	pairSet := misc.NewPairSet()

	// Same-ring: circle-neighbors only; skip degenerate < 3 rings.
	for _, sorted := range tierSorted {
		nn := len(sorted)
		if nn < 3 {
			continue
		}
		for j := 0; j < nn; j++ {
			pairSet.Add(sorted[j], sorted[(j+1)%nn])
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
				pairSet.Add(outerSorted[oi], innerSorted[best])
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
					pairSet.Add(innerSorted[ii], outerSorted[oi])
				}
			}
		}
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
	allNeutralZonePlans models.NeutralZonePlans,
	connectionNames [][]string) []template.Zone {
	var zones []template.Zone
	for index, label := range orderedLabels {
		myConns := connectionNames[index]
		if label == playerLabel {
			zones = append(zones, this.CreateSpawnZone(
				label, fmt.Sprintf("Player%d", playerIndex+1), myConns, configuration.ZoneConfiguration.PlayerZoneCastles,
				configuration.MatchPlayerCastleFactions, configuration.ZoneConfiguration.Advanced.PlayerZoneSize,
				configuration.SpawnRemoteFootholds, configuration.GenerateRoads, tuning))
		} else {
			zones = append(zones, this.CreateNeutralZone(
				linq.FromSlice(allNeutralZonePlans).FirstOrDefault(func(x models.NeutralZonePlan) bool { return x.Label == label }),
				myConns, configuration.ZoneConfiguration.Advanced.NeutralZoneSize, configuration.SpawnRemoteFootholds,
				configuration.GenerateRoads, tuning, false))
		}
	}
	return zones
}

func (this *BalancedClusterService) createConnections(
	playerLabel string,
	orderedLabels []string,
	tuning models.GenerationTuning,
	allNeutralZonePlans models.NeutralZonePlans,
	connectionNames [][]string,
	sortedPairs [][2]int) []template.Connection {
	nameLookup := make(map[int]int, len(orderedLabels))

	var connections []template.Connection
	for _, pair := range sortedPairs {
		indexA, indexB := pair[0], pair[1]
		labelFrom := orderedLabels[indexA]
		labelTo := orderedLabels[indexB]

		connName := connectionNames[indexA][nameLookup[indexA]]
		nameLookup[indexA]++
		nameLookup[indexB]++

		fromZone := "Spawn-" + labelFrom
		if labelFrom != playerLabel {
			fromZone = "Neutral-" + labelFrom
		}
		toZone := "Spawn-" + labelTo
		if labelTo != playerLabel {
			toZone = "Neutral-" + labelTo
		}
		connections = append(connections, template.Connection{
			Name: connName, From: fromZone, To: toZone,
			ConnectionType: "Direct", GuardZone: fromZone, SimTurnSquad: true,
			GuardValue: this.GetBorderGuardValue(labelFrom, labelTo, []string{playerLabel}, allNeutralZonePlans, tuning), GuardWeeklyIncrement: 0.15,
			GuardMatchGroup: fmt.Sprintf("tourney_bal_guard_%s_%s", labelFrom, labelTo),
		})
	}
	return connections
}
