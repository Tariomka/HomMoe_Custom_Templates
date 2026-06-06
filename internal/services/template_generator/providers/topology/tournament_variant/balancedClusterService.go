package tournament_variant

import (
	"fmt"
	"math"
	"sort"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
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
	var zones []template.Zone
	var connections []template.Connection

	singlePlayerList := []string{playerLabel}
	orderedLetters := buildBalancedRingLetters(singlePlayerList, playerNeutralZonePlans, 0)
	rawPos := buildBalancedRandomPositions(orderedLetters, singlePlayerList, allNeutralZonePlans)

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
		t := zoneTierRank(l, singlePlayerList, allNeutralZonePlans)
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
		if from != playerLabel {
			fromZone = "Neutral-" + from
		}
		toZone := "Spawn-" + to
		if to != playerLabel {
			toZone = "Neutral-" + to
		}
		*connections = append(*connections, template.Connection{
			Name: connName, From: fromZone, To: toZone,
			ConnectionType: "Direct", GuardZone: fromZone, SimTurnSquad: true,
			GuardValue: GetBorderGuardValue(from, to, []string{playerLabel}, allNeutralZonePlans, tuning), GuardWeeklyIncrement: 0.15,
			GuardMatchGroup: fmt.Sprintf("tourney_bal_guard_%s_%s", from, to),
		})
	}

	clusterStart := len(*zones)
	for i, letter := range orderedLetters {
		myConns := connsByZone[i]
		if letter == playerLabel {
			*zones = append(*zones, buildSpawnZone(letter, fmt.Sprintf("Player%d", playerIndex+1), myConns, configuration.ZoneConfiguration.PlayerZoneCastles, configuration.MatchPlayerCastleFactions, configuration.ZoneConfiguration.Advanced.PlayerZoneSize, configuration.SpawnRemoteFootholds, configuration.GenerateRoads, tuning))
		} else {
			*zones = append(*zones, buildNeutralZone(allNeutralZonePlans[letter], myConns, configuration.ZoneConfiguration.Advanced.NeutralZoneSize, configuration.SpawnRemoteFootholds, configuration.GenerateRoads, tuning, false))
		}
	}

	// Stamp generator positions and ring indices onto the freshly built cluster
	// zones so the preview renderer can reproduce the tournament-balanced
	// geometry without re-deriving it from connections.
	for i := 0; i < len(orderedLetters); i++ {
		p := pos[i]
		(*zones)[clusterStart+i].GeneratorPosition = &[2]float64{p[0], p[1]}
		r := zoneTierRank(orderedLetters[i], singlePlayerList, allNeutralZonePlans)
		(*zones)[clusterStart+i].GeneratorRing = &r
	}

	// Ensure the cluster is fully connected (same guarantee as the standard
	// balanced variant). Operate on the slice header of the cluster's zones.
	clusterZones := (*zones)[clusterStart:]
	ensureFullConnectivity(singlePlayerList, orderedLetters, pos, clusterZones, connections, tuning, allNeutralZonePlans)
	return zones, connections
}
