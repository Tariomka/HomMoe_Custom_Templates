package topology

import (
	"fmt"
	"sort"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
)

type TournamentTopologyService struct {
	topologyBase
}

func NewTournamentTopologyService() *TournamentTopologyService {
	return &TournamentTopologyService{
		topologyBase: newTopologyBase(),
	}
}

func (this *TournamentTopologyService) GetTopologyVariant(
	configuration config.GeneratorConfig,
	playerLetters []string,
	neutralZones models.NeutralZonePlans,
	tuning models.GenerationTuning) template.Variant {
	neutralByLetter := mapNeutralByLetter(neutralZones)

	// Distribute neutrals in a balanced round-robin: sort by descending quality
	// (then castle count, then letter) so that quality tiers are split evenly
	// across the two players (e91e79f / v0.7 ordering).
	sorted := models.NewNeutralZonePlansSorted(neutralZones)

	neutralsForPlayer := [2]models.NeutralZonePlans{}
	for i, nz := range *sorted {
		neutralsForPlayer[i%2] = append(neutralsForPlayer[i%2], nz)
	}

	for p := range 2 {
		sort.SliceStable(neutralsForPlayer[p], func(i, j int) bool {
			ai, aj := neutralsForPlayer[p][i], neutralsForPlayer[p][j]
			si, sj := neutralZoneBalanceScore(ai), neutralZoneBalanceScore(aj)
			if si != sj {
				return si < sj
			}
			return ai.Label < aj.Label
		})
	}

	var zones []template.Zone
	var conns []template.Connection

	switch configuration.Topology {
	case config.TopologyHubAndSpoke:
		for p := range 2 {
			buildTournamentHubCluster(p, playerLetters[p], neutralsForPlayer[p], neutralByLetter, configuration, tuning, &zones, &conns)
		}
	case config.TopologyBalanced:
		for p := range 2 {
			buildTournamentBalancedCluster(p, playerLetters[p], neutralsForPlayer[p], neutralByLetter, configuration, tuning, &zones, &conns)
		}
	case config.TopologyDefault:
		for p := range 2 {
			buildTournamentRingCluster(p, playerLetters[p], neutralsForPlayer[p], neutralByLetter, configuration, tuning, &zones, &conns)
		}
	default:
		// Chain, SharedWeb, Random → chain-per-cluster fallback.
		for p := range 2 {
			buildTournamentChainCluster(p, playerLetters[p], neutralsForPlayer[p], neutralByLetter, configuration, tuning, &zones, &conns)
		}
	}

	if configuration.RandomPortals {
		for p := range 2 {
			clusterLetters := []string{playerLetters[p]}
			for _, n := range neutralsForPlayer[p] {
				clusterLetters = append(clusterLetters, n.Label)
			}
			conns = append(conns, buildRandomPortalConnections(playerLetters, clusterLetters, tuning, configuration.MaxPortalConnections)...)
		}
	}

	return makeVariant(playerLetters, playerLetters[0], len(zones), zones, conns)
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
