package topology

import (
	"fmt"
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
)

type SharedWebTopologyService struct {
	topologyBase
}

func NewSharedWebTopologyService() *SharedWebTopologyService {
	return &SharedWebTopologyService{
		topologyBase: newTopologyBase(),
	}
}

func (this *SharedWebTopologyService) GetTopologyVariant(
	configuration config.GeneratorConfig,
	playerLetters []string,
	neutralZones []models.NeutralZonePlan,
	tuning models.GenerationTuning,
	holdCityNeutralLetter string) template.Variant {
	neutralByLetter := mapNeutralByLetter(neutralZones)
	var neutrals []string
	if configuration.Topology == config.TopologyBalanced {
		neutrals = buildBalancedNeutralRing(neutralZones, len(playerLetters))
	} else {
		for _, nz := range neutralZones {
			neutrals = append(neutrals, nz.Label)
		}
	}

	p := len(playerLetters)
	nn := len(neutrals)

	spokeByPlayer := map[string][]string{}
	spokeByNeutral := map[string][]string{}
	for _, l := range playerLetters {
		spokeByPlayer[l] = nil
	}
	for _, l := range neutrals {
		spokeByNeutral[l] = nil
	}
	addSpoke := func(pl, nl string) {
		cn := fmt.Sprintf("Web-%s-%s", pl, nl)
		spokeByPlayer[pl] = append(spokeByPlayer[pl], cn)
		spokeByNeutral[nl] = append(spokeByNeutral[nl], cn)
	}
	for i := 0; i < p; i++ {
		n1 := (i * nn / p) % nn
		n2 := ((i * nn / p) + 1) % nn
		addSpoke(playerLetters[i], neutrals[n1])
		if n1 != n2 {
			addSpoke(playerLetters[i], neutrals[n2])
		}
	}

	var zones []template.Zone
	var conns []template.Connection

	neutralRingConns := make([]string, nn)
	for i := 0; i < nn; i++ {
		next := (i + 1) % nn
		neutralRingConns[i] = fmt.Sprintf("NRing-%s-%s", neutrals[i], neutrals[next])
	}

	for i := 0; i < nn; i++ {
		prev := (i - 1 + nn) % nn
		var nConns []string
		if nn > 1 {
			nConns = append(nConns, neutralRingConns[prev], neutralRingConns[i])
		}
		nConns = append(nConns, spokeByNeutral[neutrals[i]]...)
		nConns = uniqueStrings(nConns)
		z := buildNeutralZone(neutralByLetter[neutrals[i]], nConns, configuration.ZoneConfiguration.Advanced.NeutralZoneSize, configuration.SpawnRemoteFootholds, configuration.GenerateRoads, tuning, neutrals[i] == holdCityNeutralLetter)
		if neutralByLetter[neutrals[i]].CastleCount == 0 {
			z.Roads = buildConnectorZoneRoads(nConns, configuration.GenerateRoads)
		}
		zones = append(zones, z)
	}

	for i := 0; i < p; i++ {
		pl := playerLetters[i]
		sc := spokeByPlayer[pl]
		zones = append(zones, buildSpawnZone(pl, fmt.Sprintf("Player%d", i+1), sc, configuration.ZoneConfiguration.PlayerZoneCastles, configuration.MatchPlayerCastleFactions, configuration.ZoneConfiguration.Advanced.PlayerZoneSize, configuration.SpawnRemoteFootholds, configuration.GenerateRoads, tuning))

		for _, cn := range sc {
			parts := strings.Split(cn, "-")
			nl := parts[2]
			conns = append(conns, template.Connection{
				Name: cn, From: "Spawn-" + pl, To: "Neutral-" + nl,
				ConnectionType: "Direct", GuardZone: "Neutral-" + nl, SimTurnSquad: true,
				GuardValue: borderGuardValue(pl, nl, playerLetters, neutralByLetter, tuning), GuardWeeklyIncrement: 0.15,
				GuardMatchGroup: fmt.Sprintf("web_guard_%s_%s", pl, nl),
			})
		}
	}

	if nn > 1 {
		for i := 0; i < nn; i++ {
			next := (i + 1) % nn
			conns = append(conns, template.Connection{
				Name: neutralRingConns[i], From: "Neutral-" + neutrals[i], To: "Neutral-" + neutrals[next],
				ConnectionType: "Direct", GuardZone: "Neutral-" + neutrals[i], SimTurnSquad: true,
				GuardValue: borderGuardValue(neutrals[i], neutrals[next], playerLetters, neutralByLetter, tuning), GuardWeeklyIncrement: 0.15,
				GuardMatchGroup: fmt.Sprintf("nring_guard_%s_%s", neutrals[i], neutrals[next]),
			})
		}
	}

	if configuration.RandomPortals {
		all := append(append([]string{}, playerLetters...), neutrals...)
		conns = append(conns, buildRandomPortalConnections(playerLetters, all, tuning, configuration.MaxPortalConnections)...)
	}
	if configuration.NoDirectPlayerConnections && len(playerLetters) > 1 {
		ensurePlayerZonesConnected(playerLetters, zones, &conns, tuning)
	}
	return makeVariant(playerLetters, playerLetters[0], len(zones), zones, conns)
}
