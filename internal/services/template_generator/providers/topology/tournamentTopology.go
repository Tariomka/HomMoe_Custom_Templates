package topology

import (
	"fmt"
	"math"
	"sort"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/tournament_variant"
)

type TournamentTopologyService struct {
	base.TopologyBase
	clusterService tournament_variant.IClusterService
}

func NewTournamentTopologyService() *TournamentTopologyService {
	return &TournamentTopologyService{
		TopologyBase: base.NewTopologyBase(),
	}
}

func (this *TournamentTopologyService) CreateTopologyVariant(
	configuration config.GeneratorConfig,
	playerLabels []string,
	neutralZones models.NeutralZonePlans,
	tuning models.GenerationTuning) template.Variant {
	neutralByLetter := mapNeutralByLetter(neutralZones)

	// Distribute neutrals in a balanced round-robin: sort by descending quality
	// (then castle count, then letter) so that quality tiers are split evenly
	// across the two players.
	sorted := models.NewNeutralZonePlansSorted(neutralZones)

	playerNeutralZones := [2]models.NeutralZonePlans{}
	for i, nz := range *sorted {
		playerNeutralZones[i%2] = append(playerNeutralZones[i%2], nz)
	}

	for p := range 2 {
		sort.SliceStable(playerNeutralZones[p], func(i, j int) bool {
			ai, aj := playerNeutralZones[p][i], playerNeutralZones[p][j]
			si, sj := neutralZoneBalanceScore(ai), neutralZoneBalanceScore(aj)
			if si != sj {
				return si < sj
			}
			return ai.Label < aj.Label
		})
	}

	switch configuration.Topology {
	case config.TopologyHubAndSpoke:
		this.clusterService = tournament_variant.NewHubClusterService()
	case config.TopologyBalanced:
		this.clusterService = tournament_variant.NewBalancedClusterService()
	case config.TopologyDefault:
		this.clusterService = tournament_variant.NewRingClusterService()
	default:
		// Chain, SharedWeb, Random → chain-per-cluster fallback.
		this.clusterService = tournament_variant.NewChainClusterService()
	}

	var zones []template.Zone
	var conns []template.Connection
	for playerIndex := range 2 {
		zones, conns = this.clusterService.CreateClusterVariant(
			configuration,
			tuning,
			neutralZones,
			playerNeutralZones[playerIndex],
			playerIndex,
			playerLabels[playerIndex])
	}

	switch configuration.Topology {
	case config.TopologyHubAndSpoke:
		for p := range 2 {
			buildTournamentHubCluster(p, playerLabels[p], playerNeutralZones[p], neutralByLetter, configuration, tuning, &zones, &conns)
		}
	case config.TopologyBalanced:
		for p := range 2 {
			buildTournamentBalancedCluster(p, playerLabels[p], playerNeutralZones[p], neutralByLetter, configuration, tuning, &zones, &conns)
		}
	case config.TopologyDefault:
		for p := range 2 {
			buildTournamentRingCluster(p, playerLabels[p], playerNeutralZones[p], neutralByLetter, configuration, tuning, &zones, &conns)
		}
	default:
		// Chain, SharedWeb, Random → chain-per-cluster fallback.
		for p := range 2 {
			buildTournamentChainCluster(p, playerLabels[p], playerNeutralZones[p], neutralByLetter, configuration, tuning, &zones, &conns)
		}
	}

	if configuration.RandomPortals {
		for playerIndex := range 2 {
			clusterLabels := linq.FromSlice(playerNeutralZones[playerIndex]).
				SelectString(func(x models.NeutralZonePlan) string { return x.Label }).
				ToSlice()
			conns = append(conns, this.CreateRandomPortalConnections(playerLabels, clusterLabels, tuning, configuration.MaxPortalConnections)...)
		}
	}
	return this.CreateVariant(playerLabels, playerLabels[0], len(zones), zones, conns)
}

// buildTournamentHubCluster — one player's isolated cluster as a private
// hub-and-spoke layout. A dedicated mini-hub zone "Hub-{playerLetter}" sits
// at the centre and connects directly to the player spawn and all of their
// exclusive neutrals
func buildTournamentHubCluster(playerIndex int, playerLetter string, myNeutrals models.NeutralZonePlans, neutralByLetter map[string]models.NeutralZonePlan, settings *config.GeneratorConfig, tuning models.GenerationTuning, zones *[]template.Zone, connections *[]template.Connection) {
	hubName := "Hub-" + playerLetter

	spokeLetters := []string{playerLetter}
	for _, nz := range myNeutrals {
		spokeLetters = append(spokeLetters, nz.Label)
	}

	spokeConnNames := make([]string, len(spokeLetters))
	for i, l := range spokeLetters {
		spokeConnNames[i] = fmt.Sprintf("THubSpoke-%s-%s", playerLetter, l)
	}

	hubZone := buildHubZone(spokeConnNames, tuning, false, settings.ZoneConfiguration.HubZoneSize, settings.ZoneConfiguration.HubZoneCastles, settings.GenerateRoads)
	hubZone.Name = hubName
	*zones = append(*zones, hubZone)

	for i, letter := range spokeLetters {
		conn := []string{spokeConnNames[i]}
		if i == 0 {
			*zones = append(*zones, buildSpawnZone(letter, fmt.Sprintf("Player%d", playerIndex+1), conn, settings.ZoneConfiguration.PlayerZoneCastles, settings.MatchPlayerCastleFactions, settings.ZoneConfiguration.Advanced.PlayerZoneSize, settings.SpawnRemoteFootholds, settings.GenerateRoads, tuning))
		} else {
			*zones = append(*zones, buildNeutralZone(neutralByLetter[letter], conn, settings.ZoneConfiguration.Advanced.NeutralZoneSize, settings.SpawnRemoteFootholds, settings.GenerateRoads, tuning, false))
		}
	}

	for i, spokeLetter := range spokeLetters {
		spokeZone := "Spawn-" + spokeLetter
		if i != 0 {
			spokeZone = "Neutral-" + spokeLetter
		}
		*connections = append(*connections, template.Connection{
			Name: spokeConnNames[i], From: hubName, To: spokeZone,
			ConnectionType: "Direct", GuardZone: hubName, SimTurnSquad: true,
			GuardValue: borderGuardValue(playerLetter, spokeLetter, []string{playerLetter}, neutralByLetter, tuning), GuardWeeklyIncrement: 0.15,
			GuardMatchGroup: fmt.Sprintf("tourney_hub_guard_%s_%s", playerLetter, spokeLetter),
		})
	}

	// Proximity ring around spokes so the engine arranges them sensibly.
	for i := 0; i < len(spokeLetters); i++ {
		next := (i + 1) % len(spokeLetters)
		from := spokeLetters[i]
		to := spokeLetters[next]
		fromZone := "Spawn-" + from
		if i != 0 {
			fromZone = "Neutral-" + from
		}
		toZone := "Spawn-" + to
		if next != 0 {
			toZone = "Neutral-" + to
		}
		*connections = append(*connections, template.Connection{
			Name:           fmt.Sprintf("TPseudo-%s-%s-%s", playerLetter, from, to),
			From:           fromZone,
			To:             toZone,
			ConnectionType: "Proximity",
		})
	}
}

func buildTournamentBalancedCluster(playerIndex int, playerLetter string, myNeutrals []neutralZonePlan, neutralByLetter map[string]neutralZonePlan, settings *config.GeneratorConfig, tuning generationTuning, zones *[]template.Zone, connections *[]template.Connection) {
	singlePlayerList := []string{playerLetter}
	orderedLetters := buildBalancedRingLetters(singlePlayerList, myNeutrals, 0)
	rawPos := buildBalancedRandomPositions(orderedLetters, singlePlayerList, neutralByLetter)

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
		t := zoneTierRank(l, singlePlayerList, neutralByLetter)
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
		if from != playerLetter {
			fromZone = "Neutral-" + from
		}
		toZone := "Spawn-" + to
		if to != playerLetter {
			toZone = "Neutral-" + to
		}
		*connections = append(*connections, template.Connection{
			Name: connName, From: fromZone, To: toZone,
			ConnectionType: "Direct", GuardZone: fromZone, SimTurnSquad: true,
			GuardValue: GetBorderGuardValue(from, to, []string{playerLetter}, neutralByLetter, tuning), GuardWeeklyIncrement: 0.15,
			GuardMatchGroup: fmt.Sprintf("tourney_bal_guard_%s_%s", from, to),
		})
	}

	clusterStart := len(*zones)
	for i, letter := range orderedLetters {
		myConns := connsByZone[i]
		if letter == playerLetter {
			*zones = append(*zones, buildSpawnZone(letter, fmt.Sprintf("Player%d", playerIndex+1), myConns, settings.ZoneConfiguration.PlayerZoneCastles, settings.MatchPlayerCastleFactions, settings.ZoneConfiguration.Advanced.PlayerZoneSize, settings.SpawnRemoteFootholds, settings.GenerateRoads, tuning))
		} else {
			*zones = append(*zones, buildNeutralZone(neutralByLetter[letter], myConns, settings.ZoneConfiguration.Advanced.NeutralZoneSize, settings.SpawnRemoteFootholds, settings.GenerateRoads, tuning, false))
		}
	}

	// Stamp generator positions and ring indices onto the freshly built cluster
	// zones so the preview renderer can reproduce the tournament-balanced
	// geometry without re-deriving it from connections.
	for i := 0; i < len(orderedLetters); i++ {
		p := pos[i]
		(*zones)[clusterStart+i].GeneratorPosition = &[2]float64{p[0], p[1]}
		r := zoneTierRank(orderedLetters[i], singlePlayerList, neutralByLetter)
		(*zones)[clusterStart+i].GeneratorRing = &r
	}

	// Ensure the cluster is fully connected (same guarantee as the standard
	// balanced variant). Operate on the slice header of the cluster's zones.
	clusterZones := (*zones)[clusterStart:]
	ensureFullConnectivity(singlePlayerList, orderedLetters, pos, clusterZones, connections, tuning, neutralByLetter)
}

func buildTournamentChainCluster(playerIndex int, playerLetter string, myNeutrals []neutralZonePlan, neutralByLetter map[string]neutralZonePlan, settings *config.GeneratorConfig, tuning generationTuning, zones *[]template.Zone, connections *[]template.Connection) {
	chain := []string{playerLetter}
	for _, n := range myNeutrals {
		chain = append(chain, n.Letter)
	}
	connNames := make([]string, len(chain)-1)
	for i := range connNames {
		connNames[i] = fmt.Sprintf("Tourney-%s-%s", chain[i], chain[i+1])
	}
	for i, letter := range chain {
		var myConns []string
		if i > 0 {
			myConns = append(myConns, connNames[i-1])
		}
		if i < len(connNames) {
			myConns = append(myConns, connNames[i])
		}
		if i == 0 {
			*zones = append(*zones, buildSpawnZone(letter, fmt.Sprintf("Player%d", playerIndex+1), myConns, settings.ZoneConfiguration.PlayerZoneCastles, settings.MatchPlayerCastleFactions, settings.ZoneConfiguration.Advanced.PlayerZoneSize, settings.SpawnRemoteFootholds, settings.GenerateRoads, tuning))
		} else {
			*zones = append(*zones, buildNeutralZone(neutralByLetter[letter], myConns, settings.ZoneConfiguration.Advanced.NeutralZoneSize, settings.SpawnRemoteFootholds, settings.GenerateRoads, tuning, false))
		}
	}
	for i := range connNames {
		from := chain[i]
		to := chain[i+1]
		fromZone := "Spawn-" + from
		if i > 0 {
			fromZone = "Neutral-" + from
		}
		*connections = append(*connections, template.Connection{
			Name: connNames[i], From: fromZone, To: "Neutral-" + to,
			ConnectionType: "Direct", GuardZone: fromZone, SimTurnSquad: true,
			GuardValue: GetBorderGuardValue(from, to, []string{playerLetter}, neutralByLetter, tuning), GuardWeeklyIncrement: 0.15,
			GuardMatchGroup: fmt.Sprintf("tourney_guard_%s_%s", from, to),
		})
	}
}

// buildTournamentRingCluster — one player's isolated cluster as a ring:
// player → low → … → high → … → low → player. Sorts neutrals by balance
// score, then fills outward-from-player so the highest-quality zones sit
// at the ring midpoint
func buildTournamentRingCluster(playerIndex int, playerLetter string, myNeutrals []neutralZonePlan, neutralByLetter map[string]neutralZonePlan, settings *config.GeneratorConfig, tuning generationTuning, zones *[]template.Zone, connections *[]template.Connection) {
	sortedNeutrals := make([]neutralZonePlan, len(myNeutrals))
	copy(sortedNeutrals, myNeutrals)
	sort.SliceStable(sortedNeutrals, func(i, j int) bool {
		si, sj := neutralZoneBalanceScore(sortedNeutrals[i]), neutralZoneBalanceScore(sortedNeutrals[j])
		if si != sj {
			return si < sj
		}
		return sortedNeutrals[i].Letter < sortedNeutrals[j].Letter
	})

	n := len(sortedNeutrals)
	orderedNeutrals := make([]neutralZonePlan, n)
	lo, hi := 0, n-1
	for i := range n {
		if i%2 == 0 {
			orderedNeutrals[lo] = sortedNeutrals[i]
			lo++
		} else {
			orderedNeutrals[hi] = sortedNeutrals[i]
			hi--
		}
	}

	ringLetters := []string{playerLetter}
	for _, nz := range orderedNeutrals {
		ringLetters = append(ringLetters, nz.Letter)
	}
	count := len(ringLetters)
	if count < 2 {
		// Lone player zone — no ring edges possible.
		*zones = append(*zones, buildSpawnZone(playerLetter, fmt.Sprintf("Player%d", playerIndex+1), nil, settings.ZoneConfiguration.PlayerZoneCastles, settings.MatchPlayerCastleFactions, settings.ZoneConfiguration.Advanced.PlayerZoneSize, settings.SpawnRemoteFootholds, settings.GenerateRoads, tuning))
		return
	}

	connNames := make([]string, count)
	for i := range count {
		next := (i + 1) % count
		connNames[i] = fmt.Sprintf("TRing-%s-%s", ringLetters[i], ringLetters[next])
	}

	for i, letter := range ringLetters {
		prev := (i - 1 + count) % count
		seen := map[string]bool{}
		var myConns []string
		for _, name := range []string{connNames[prev], connNames[i]} {
			if !seen[name] {
				seen[name] = true
				myConns = append(myConns, name)
			}
		}
		if i == 0 {
			*zones = append(*zones, buildSpawnZone(letter, fmt.Sprintf("Player%d", playerIndex+1), myConns, settings.ZoneConfiguration.PlayerZoneCastles, settings.MatchPlayerCastleFactions, settings.ZoneConfiguration.Advanced.PlayerZoneSize, settings.SpawnRemoteFootholds, settings.GenerateRoads, tuning))
		} else {
			*zones = append(*zones, buildNeutralZone(neutralByLetter[letter], myConns, settings.ZoneConfiguration.Advanced.NeutralZoneSize, settings.SpawnRemoteFootholds, settings.GenerateRoads, tuning, false))
		}
	}

	for i := range count {
		next := (i + 1) % count
		from := ringLetters[i]
		to := ringLetters[next]
		fromZone := "Spawn-" + from
		if i != 0 {
			fromZone = "Neutral-" + from
		}
		toZone := "Spawn-" + to
		if next != 0 {
			toZone = "Neutral-" + to
		}
		*connections = append(*connections, template.Connection{
			Name: connNames[i], From: fromZone, To: toZone,
			ConnectionType: "Direct", GuardZone: fromZone, SimTurnSquad: true,
			GuardValue: GetBorderGuardValue(from, to, []string{playerLetter}, neutralByLetter, tuning), GuardWeeklyIncrement: 0.15,
			GuardMatchGroup: fmt.Sprintf("tourney_ring_guard_%s_%s", from, to),
		})
	}
}
