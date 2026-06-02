package topology

import (
	"fmt"
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
)

type RingTopologyService struct {
	topologyBase
}

func NewRingTopologyService() *RingTopologyService {
	return &RingTopologyService{
		topologyBase: newTopologyBase(),
	}
}

func (this *RingTopologyService) CreateTopologyVariant(
	configuration config.GeneratorConfig,
	playerLetters []string,
	neutralZones []models.NeutralZonePlan,
	tuning models.GenerationTuning,
	holdCityNeutralLetter string) template.Variant {
	neutralByLetter := mapNeutralByLetter(neutralZones) // TODO: Remove

	orderedLabels := this.zoneLabelProvider.CreateOrderedZoneLabels(configuration, playerLetters, neutralZones, true)
	labelCount := len(orderedLabels)
	isolate := configuration.NoDirectPlayerConnections && len(playerLetters) > 1

	ringConnRight := make([]string, labelCount)
	ringConnLeft := make([]string, labelCount)
	for i := range labelCount {
		next := (i + 1) % labelCount
		if isolate && slices.Contains(playerLetters, orderedLabels[i]) && slices.Contains(playerLetters, orderedLabels[next]) {
			continue
		}
		name := fmt.Sprintf("Ring-%s-%s", orderedLabels[i], orderedLabels[next])
		ringConnRight[i] = name
		ringConnLeft[next] = name
	}

	var zones []template.Zone
	for i := range labelCount {
		label := orderedLabels[i]
		var myConns []string
		if ringConnLeft[i] != "" {
			myConns = append(myConns, ringConnLeft[i])
		}
		if ringConnRight[i] != "" && ringConnRight[i] != ringConnLeft[i] {
			myConns = append(myConns, ringConnRight[i])
		}
		if pi := slices.Index(playerLetters, label); pi >= 0 {
			zones = append(zones,
				this.GetSpawnZone(
					label, fmt.Sprintf("Player%d", pi+1), myConns, configuration.ZoneConfiguration.PlayerZoneCastles,
					configuration.MatchPlayerCastleFactions, configuration.ZoneConfiguration.Advanced.PlayerZoneSize,
					configuration.SpawnRemoteFootholds, configuration.GenerateRoads, tuning))
		} else {
			zones = append(zones,
				this.CreateNeutralZone(
					linq.FromSlice(neutralZones).FirstOrDefault(func(nz models.NeutralZonePlan) bool { return nz.Label == label }),
					myConns, configuration.ZoneConfiguration.Advanced.NeutralZoneSize,
					configuration.SpawnRemoteFootholds, configuration.GenerateRoads, tuning, label == holdCityNeutralLetter))
		}
	}

	conns := buildRingConnections(playerLetters, orderedLabels, tuning, isolate, neutralByLetter)
	if configuration.RandomPortals {
		conns = append(conns, buildRandomPortalConnections(playerLetters, orderedLabels, tuning, configuration.MaxPortalConnections)...)
	}
	if isolate {
		ensurePlayerZonesConnected(playerLetters, zones, &conns, tuning)
	}
	return makeVariant(playerLetters, orderedLabels[0], labelCount, zones, conns)
}
