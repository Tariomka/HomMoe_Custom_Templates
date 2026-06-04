package topology

import (
	"fmt"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
)

type HubTopologyService struct {
	topologyBase
}

func NewHubTopologyService() *HubTopologyService {
	return &HubTopologyService{
		topologyBase: newTopologyBase(),
	}
}

func (this *HubTopologyService) GetTopologyVariant(
	configuration config.GeneratorConfig,
	playerLetters []string,
	neutralZones models.NeutralZonePlans,
	tuning models.GenerationTuning,
	hubIsHoldCity bool) template.Variant {
	neutralByLetter := mapNeutralByLetter(neutralZones)
	neutralLetters := make([]string, len(neutralZones))
	for i, nz := range neutralZones {
		neutralLetters[i] = nz.Label
	}

	var outerLetters []string
	if configuration.Topology == config.TopologyBalanced {
		sep := 0
		if configuration.MinNeutralZonesBetweenPlayers > 0 && canHonorNeutralSeparation(configuration, len(neutralZones)) {
			sep = configuration.MinNeutralZonesBetweenPlayers
		}
		outerLetters = buildBalancedRingLetters(playerLetters, neutralZones, sep)
	} else {
		outerLetters = append(append([]string{}, playerLetters...), neutralLetters...)
	}

	var zones []template.Zone
	var conns []template.Connection

	hubConns := make([]string, len(outerLetters))
	for i, l := range outerLetters {
		hubConns[i] = "Hub-" + l
	}
	zones = append(zones, buildHubZone(hubConns, tuning, hubIsHoldCity, configuration.ZoneConfiguration.HubZoneSize, configuration.ZoneConfiguration.HubZoneCastles, configuration.GenerateRoads))

	for i, letter := range outerLetters {
		spokeConns := []string{"Hub-" + letter}
		if pi := indexOf(playerLetters, letter); pi >= 0 {
			zones = append(zones, buildSpawnZone(letter, fmt.Sprintf("Player%d", pi+1), spokeConns, configuration.ZoneConfiguration.PlayerZoneCastles, configuration.MatchPlayerCastleFactions, configuration.ZoneConfiguration.Advanced.PlayerZoneSize, configuration.SpawnRemoteFootholds, configuration.GenerateRoads, tuning))
		} else {
			zones = append(zones, buildNeutralZone(neutralByLetter[letter], spokeConns, configuration.ZoneConfiguration.Advanced.NeutralZoneSize, configuration.SpawnRemoteFootholds, configuration.GenerateRoads, tuning, false))
		}
		_ = i
	}

	for _, letter := range outerLetters {
		outerZone := createZoneName(letter, playerLetters)
		hubAnchor := letter
		if len(playerLetters) > 0 {
			hubAnchor = playerLetters[0]
		}
		hubGuard := borderGuardValue(hubAnchor, letter, playerLetters, neutralByLetter, tuning)
		conns = append(conns,
			template.Connection{
				Name:                 "Hub-" + letter,
				From:                 "Hub",
				To:                   outerZone,
				ConnectionType:       "Direct",
				GuardZone:            "Hub",
				SimTurnSquad:         true,
				GuardValue:           hubGuard,
				GuardWeeklyIncrement: 0.15,
				GuardMatchGroup:      "hub_guard_" + letter,
			},
			template.Connection{
				From:                 "Hub",
				To:                   outerZone,
				ConnectionType:       "Direct",
				GuardZone:            "Hub",
				SimTurnSquad:         true,
				GuardValue:           hubGuard,
				GuardWeeklyIncrement: 0.15,
				GuardMatchGroup:      fmt.Sprintf("hub_guard_%s_%d", letter, 1),
			})
	}

	// Proximity ring
	for i := 0; i < len(outerLetters); i++ {
		next := (i + 1) % len(outerLetters)
		from := outerLetters[i]
		to := outerLetters[next]
		if configuration.NoDirectPlayerConnections && contains(playerLetters, from) && contains(playerLetters, to) {
			continue
		}
		conns = append(conns, template.Connection{
			Name: fmt.Sprintf("Pseudo-%s-%s", from, to),
			From: createZoneName(from, playerLetters), To: createZoneName(to, playerLetters),
			ConnectionType: "Proximity",
		})
	}

	if configuration.RandomPortals {
		conns = append(conns, buildRandomPortalConnections(playerLetters, outerLetters, tuning, configuration.MaxPortalConnections)...)
	}
	return makeVariant(playerLetters, outerLetters[0], len(outerLetters)+1, zones, conns)
}
