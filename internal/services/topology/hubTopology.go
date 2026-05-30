package topology

import (
	"fmt"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
)

func buildVariantHubAndSpoke(settings *config.GeneratorConfig, playerLetters []string, neutralZones []neutralZonePlan, tuning generationTuning, hubIsHoldCity bool) template.Variant {
	neutralByLetter := mapNeutralByLetter(neutralZones)
	neutralLetters := make([]string, len(neutralZones))
	for i, nz := range neutralZones {
		neutralLetters[i] = nz.Letter
	}

	var outerLetters []string
	if settings.Topology == config.TopologyBalanced {
		sep := 0
		if settings.MinNeutralZonesBetweenPlayers > 0 && canHonorNeutralSeparation(settings, len(neutralZones)) {
			sep = settings.MinNeutralZonesBetweenPlayers
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
	zones = append(zones, buildHubZone(hubConns, tuning, hubIsHoldCity, settings.ZoneConfiguration.HubZoneSize, settings.ZoneConfiguration.HubZoneCastles, settings.GenerateRoads))

	for i, letter := range outerLetters {
		spokeConns := []string{"Hub-" + letter}
		if pi := indexOf(playerLetters, letter); pi >= 0 {
			zones = append(zones, buildSpawnZone(letter, fmt.Sprintf("Player%d", pi+1), spokeConns, settings.ZoneConfiguration.PlayerZoneCastles, settings.MatchPlayerCastleFactions, settings.ZoneConfiguration.Advanced.PlayerZoneSize, settings.SpawnRemoteFootholds, settings.GenerateRoads, tuning))
		} else {
			zones = append(zones, buildNeutralZone(neutralByLetter[letter], spokeConns, settings.ZoneConfiguration.Advanced.NeutralZoneSize, settings.SpawnRemoteFootholds, settings.GenerateRoads, tuning, false))
		}
		_ = i
	}

	for _, letter := range outerLetters {
		outerZone := zoneName(letter, playerLetters)
		hubAnchor := letter
		if len(playerLetters) > 0 {
			hubAnchor = playerLetters[0]
		}
		hubGuard := borderGuardValue(hubAnchor, letter, playerLetters, neutralByLetter, tuning)
		conns = append(conns,
			template.Connection{
				Name: "Hub-" + letter, From: "Hub", To: outerZone,
				ConnectionType: "Direct", GuardZone: "Hub", SimTurnSquad: true,
				GuardValue: hubGuard, GuardWeeklyIncrement: 0.15,
				GuardMatchGroup: "hub_guard_" + letter,
			},
			template.Connection{
				From: "Hub", To: outerZone, ConnectionType: "Direct",
				GuardZone: "Hub", SimTurnSquad: true,
				GuardValue: hubGuard, GuardWeeklyIncrement: 0.15,
				GuardMatchGroup: fmt.Sprintf("hub_guard_%s_%d", letter, 1),
			})
	}

	// Proximity ring
	for i := 0; i < len(outerLetters); i++ {
		next := (i + 1) % len(outerLetters)
		from := outerLetters[i]
		to := outerLetters[next]
		if settings.NoDirectPlayerConnections && contains(playerLetters, from) && contains(playerLetters, to) {
			continue
		}
		conns = append(conns, template.Connection{
			Name: fmt.Sprintf("Pseudo-%s-%s", from, to),
			From: zoneName(from, playerLetters), To: zoneName(to, playerLetters),
			ConnectionType: "Proximity",
		})
	}

	if settings.RandomPortals {
		conns = append(conns, buildRandomPortalConnections(playerLetters, outerLetters, tuning, settings.MaxPortalConnections)...)
	}
	return makeVariant(playerLetters, outerLetters[0], len(outerLetters)+1, zones, conns)
}
