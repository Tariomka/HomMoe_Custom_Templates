package topology

import (
	"fmt"
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/builders/variant_content"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
)

type RingTopologyService struct {
	base.TopologyBase
}

func NewRingTopologyService() *RingTopologyService {
	return &RingTopologyService{
		TopologyBase: base.NewTopologyBase(),
	}
}

func (this *RingTopologyService) CreateTopologyVariant(
	configuration config.GeneratorConfig,
	playerLabels []string,
	neutralZones models.NeutralZonePlans,
	tuning models.GenerationTuning,
	holdCityNeutralLabel string) template.Variant {
	orderedLabels := this.ZoneLabelProvider.CreateOrderedZoneLabels(configuration, playerLabels, neutralZones, true)
	isIsolated := configuration.NoDirectPlayerConnections && len(playerLabels) > 1

	zones := this.createZones(configuration, playerLabels, orderedLabels, tuning, isIsolated, neutralZones, holdCityNeutralLabel)
	conns := this.createConnections(playerLabels, orderedLabels, tuning, isIsolated, neutralZones)
	if configuration.RandomPortals {
		conns = append(conns, this.CreateRandomPortalConnections(playerLabels, orderedLabels, tuning, configuration.MaxPortalConnections)...)
	}
	if isIsolated {
		conns = append(conns, this.CreateMissingPlayerConnections(playerLabels, zones, conns, tuning)...)
	}
	return this.CreateVariant(playerLabels, orderedLabels[0], len(orderedLabels), zones, conns)
}

func (this *RingTopologyService) createZones(
	configuration config.GeneratorConfig,
	playerLabels, orderedLabels []string,
	tuning models.GenerationTuning,
	isIsolated bool,
	neutralZones models.NeutralZonePlans,
	holdCityNeutralLabel string) []template.Zone {

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
		var connNames []string
		if ringConnLeft[i] != "" {
			connNames = append(connNames, ringConnLeft[i])
		}
		if ringConnRight[i] != "" && ringConnRight[i] != ringConnLeft[i] {
			connNames = append(connNames, ringConnRight[i])
		}
		if pi := slices.Index(playerLabels, label); pi >= 0 {
			zones = append(zones,
				this.CreateSpawnZone(
					label, fmt.Sprintf("Player%d", pi+1), connNames, configuration.ZoneConfiguration.PlayerZoneCastles,
					configuration.MatchPlayerCastleFactions, configuration.ZoneConfiguration.Advanced.PlayerZoneSize,
					configuration.SpawnRemoteFootholds, configuration.GenerateRoads, tuning))
		} else {
			zones = append(zones,
				this.CreateNeutralZone(
					linq.FromSlice(neutralZones).FirstOrDefault(func(x models.NeutralZonePlan) bool { return x.Label == label }),
					connNames, configuration.ZoneConfiguration.Advanced.NeutralZoneSize,
					configuration.SpawnRemoteFootholds, configuration.GenerateRoads, tuning, label == holdCityNeutralLabel))
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

	var connections []template.Connection
	for i := 0; i < count; i++ {
		labelFrom := orderedLabels[i]
		labelTo := orderedLabels[(i+1)%count]
		if isIsolated && slices.Contains(playerLabels, labelFrom) && slices.Contains(playerLabels, labelTo) {
			continue
		}

		zoneFrom := this.ZoneLabelProvider.CreateZoneName(labelFrom, playerLabels)
		zoneTo := this.ZoneLabelProvider.CreateZoneName(labelTo, playerLabels)
		connections = append(connections, variant_content.NewConnectionBuilder().
			WithName(fmt.Sprintf("Ring-%s-%s", labelFrom, labelTo)).
			WithFrom(zoneFrom).
			WithTo(zoneTo).
			WithConnectionTypeDirect().
			WithGuardZone(zoneFrom).
			WithSimTurnSquad().
			WithGuardValue(this.GetBorderGuardValue(labelFrom, labelTo, playerLabels, neutralZones, tuning)).
			WithGuardWeeklyIncrement(0.15).
			WithGuardMatchGroup(fmt.Sprintf("ring_guard_%s_%s", labelFrom, labelTo)).
			Build())
	}
	return connections
}
