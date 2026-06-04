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
	playerLabels []string,
	neutralZones models.NeutralZonePlans,
	tuning models.GenerationTuning,
	holdCityNeutralLabel string) template.Variant {

	orderedLabels := this.zoneLabelProvider.CreateOrderedZoneLabels(configuration, playerLabels, neutralZones, true)
	labelCount := len(orderedLabels)
	isolate := configuration.NoDirectPlayerConnections && len(playerLabels) > 1

	zones := this.createZones(configuration, playerLabels, orderedLabels, tuning, isolate, neutralZones, holdCityNeutralLabel)
	conns := this.createConnections(playerLabels, orderedLabels, tuning, isolate, neutralZones)
	if configuration.RandomPortals {
		conns = append(conns, this.CreateRandomPortalConnections(playerLabels, orderedLabels, tuning, configuration.MaxPortalConnections)...)
	}
	if isolate {
		this.EnsurePlayerZonesConnected(playerLabels, zones, &conns, tuning)
	}
	return this.CreateVariant(playerLabels, orderedLabels[0], labelCount, zones, conns)
}

func (this *RingTopologyService) createZones(
	configuration config.GeneratorConfig,
	playerLabels, orderedLabels []string,
	tuning models.GenerationTuning,
	isIsolated bool,
	neutralZones models.NeutralZonePlans,
	holdCityNeutralLetter string) []template.Zone {

	labelCount := len(orderedLabels)

	ringConnRight := make([]string, labelCount)
	ringConnLeft := make([]string, labelCount)
	for i := range labelCount {
		next := (i + 1) % labelCount
		if isIsolated && slices.Contains(playerLabels, orderedLabels[i]) && slices.Contains(playerLabels, orderedLabels[next]) {
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
		if pi := slices.Index(playerLabels, label); pi >= 0 {
			zones = append(zones,
				this.CreateSpawnZone(
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
	return zones
}

func (this *RingTopologyService) createConnections(
	playerLabels, orderedLabels []string,
	tuning models.GenerationTuning,
	isIsolated bool,
	neutralZones models.NeutralZonePlans) []template.Connection {
	count := len(orderedLabels)
	if count < 2 {
		return nil
	}
	var conns []template.Connection
	for i := 0; i < count; i++ {
		next := (i + 1) % count
		from := orderedLabels[i]
		to := orderedLabels[next]
		if isIsolated && slices.Contains(playerLabels, from) && slices.Contains(playerLabels, to) {
			continue
		}
		fromZone := createZoneName(from, playerLabels)
		toZone := createZoneName(to, playerLabels)
		conns = append(conns, template.Connection{
			Name: fmt.Sprintf("Ring-%s-%s", from, to), From: fromZone, To: toZone,
			ConnectionType: "Direct", GuardZone: fromZone, SimTurnSquad: true,
			GuardValue:           this.GetBorderGuardValue(from, to, playerLabels, neutralZones, tuning),
			GuardWeeklyIncrement: 0.15,
			GuardMatchGroup:      fmt.Sprintf("ring_guard_%s_%s", from, to),
		})
	}
	return conns
}
