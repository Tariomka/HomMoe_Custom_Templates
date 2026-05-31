package topology

import (
	"fmt"
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
)

type RingTopologyService struct {
	topologyBase
}

func NewRingTopologyService() *RingTopologyService {
	return &RingTopologyService{}
}

func (this *RingTopologyService) GetTopologyVariant(
	settings config.GeneratorConfig,
	playerLetters []string,
	neutralZones []models.NeutralZonePlan,
	tuning models.GenerationTuning,
	holdCityNeutralLetter string) template.Variant {
	neutralByLetter := mapNeutralByLetter(neutralZones)
	ordered := buildOrderedLetters(settings, playerLetters, neutralZones, true)
	n := len(ordered)
	isolate := settings.NoDirectPlayerConnections && len(playerLetters) > 1

	ringConnRight := make([]string, n)
	ringConnLeft := make([]string, n)
	for i := range n {
		next := (i + 1) % n
		if isolate && slices.Contains(playerLetters, ordered[i]) && slices.Contains(playerLetters, ordered[next]) {
			continue
		}
		name := fmt.Sprintf("Ring-%s-%s", ordered[i], ordered[next])
		ringConnRight[i] = name
		ringConnLeft[next] = name
	}

	var zones []template.Zone
	for i := range n {
		letter := ordered[i]
		var myConns []string
		if ringConnLeft[i] != "" {
			myConns = append(myConns, ringConnLeft[i])
		}
		if ringConnRight[i] != "" && ringConnRight[i] != ringConnLeft[i] {
			myConns = append(myConns, ringConnRight[i])
		}
		if pi := slices.Index(playerLetters, letter); pi >= 0 {
			zones = append(zones,
				this.getSpawnZone(
					letter, fmt.Sprintf("Player%d", pi+1), myConns, settings.ZoneConfiguration.PlayerZoneCastles,
					settings.MatchPlayerCastleFactions, settings.ZoneConfiguration.Advanced.PlayerZoneSize,
					settings.SpawnRemoteFootholds, settings.GenerateRoads, tuning))
		} else {
			zones = append(zones,
				buildNeutralZone(
					neutralByLetter[letter],
					myConns, settings.ZoneConfiguration.Advanced.NeutralZoneSize,
					settings.SpawnRemoteFootholds, settings.GenerateRoads, tuning, letter == holdCityNeutralLetter))
		}
	}

	conns := buildRingConnections(playerLetters, ordered, tuning, isolate, neutralByLetter)
	if settings.RandomPortals {
		conns = append(conns, buildRandomPortalConnections(playerLetters, ordered, tuning, settings.MaxPortalConnections)...)
	}
	if isolate {
		ensurePlayerZonesConnected(playerLetters, zones, &conns, tuning)
	}
	return makeVariant(playerLetters, ordered[0], n, zones, conns)
}
