package topology

import (
	"fmt"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
)

type ChainTopologyService struct {
	topologyBase
}

func NewChainTopologyService() *ChainTopologyService {
	return &ChainTopologyService{
		topologyBase: newTopologyBase(),
	}
}

func (this *ChainTopologyService) GetTopologyVariant(
	configuration config.GeneratorConfig,
	playerLetters []string,
	neutralZones []models.NeutralZonePlan,
	tuning models.GenerationTuning,
	holdCityNeutralLetter string) template.Variant {
	neutralByLetter := mapNeutralByLetter(neutralZones)
	ordered := buildOrderedLetters(configuration, playerLetters, neutralZones, false)
	n := len(ordered)
	isolate := configuration.NoDirectPlayerConnections && len(playerLetters) > 1

	connNames := make([]string, n-1)
	for i := 0; i < n-1; i++ {
		if isolate && contains(playerLetters, ordered[i]) && contains(playerLetters, ordered[i+1]) {
			continue
		}
		connNames[i] = fmt.Sprintf("Chain-%s-%s", ordered[i], ordered[i+1])
	}

	var zones []template.Zone
	for i := range n {
		letter := ordered[i]
		var myConns []string
		if i > 0 && connNames[i-1] != "" {
			myConns = append(myConns, connNames[i-1])
		}
		if i < n-1 && connNames[i] != "" {
			myConns = append(myConns, connNames[i])
		}
		if pi := indexOf(playerLetters, letter); pi >= 0 {
			zones = append(zones, buildSpawnZone(letter, fmt.Sprintf("Player%d", pi+1), myConns, configuration.ZoneConfiguration.PlayerZoneCastles, configuration.MatchPlayerCastleFactions, configuration.ZoneConfiguration.Advanced.PlayerZoneSize, configuration.SpawnRemoteFootholds, configuration.GenerateRoads, tuning))
		} else {
			zones = append(zones, buildNeutralZone(neutralByLetter[letter], myConns, configuration.ZoneConfiguration.Advanced.NeutralZoneSize, configuration.SpawnRemoteFootholds, configuration.GenerateRoads, tuning, letter == holdCityNeutralLetter))
		}
	}

	var conns []template.Connection
	for i := 0; i < n-1; i++ {
		if connNames[i] == "" {
			continue
		}
		from := ordered[i]
		to := ordered[i+1]
		fromZone := zoneName(from, playerLetters)
		toZone := zoneName(to, playerLetters)
		conns = append(conns, template.Connection{
			Name: connNames[i], From: fromZone, To: toZone,
			ConnectionType: "Direct", GuardZone: fromZone, SimTurnSquad: true,
			GuardValue: borderGuardValue(from, to, playerLetters, neutralByLetter, tuning), GuardWeeklyIncrement: 0.15,
			GuardMatchGroup: fmt.Sprintf("chain_guard_%s_%s", from, to),
		})
	}
	if configuration.RandomPortals {
		conns = append(conns, buildRandomPortalConnections(playerLetters, ordered, tuning, configuration.MaxPortalConnections)...)
	}
	if isolate {
		ensurePlayerZonesConnected(playerLetters, zones, &conns, tuning)
	}
	return makeVariant(playerLetters, ordered[0], n, zones, conns)
}
