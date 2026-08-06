package topology

import (
	"fmt"
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_connections"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
	zone_services "github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
)

type RingTopologyService struct {
	base.TopologyBase
}

func NewRingTopologyService(
	zoneFactory *zone_services.ZoneFactory,
	roadFactory *zone_services.RoadFactory,
	zoneLabelProvider zone_services.IZoneLabelProvider,
	connectionService *base.TopologyConnectionService) *RingTopologyService {
	return &RingTopologyService{
		TopologyBase: base.NewTopologyBase(zoneFactory, roadFactory, zoneLabelProvider, connectionService),
	}
}

func (this *RingTopologyService) CreateTopologyVariant(
	configuration config.GeneratorConfig,
	playerLabels []string,
	neutralZones neutral_zone.Plans,
	tuning models.GenerationTuning,
	holdCityNeutralLabel string) entities.Variant {
	orderedLabels := this.ZoneLabelProvider.CreateOrderedZoneLabels(configuration, playerLabels, neutralZones, true)
	isIsolated := configuration.NoDirectPlayerConnections && len(playerLabels) > 1

	zones := this.createZones(
		configuration, playerLabels, orderedLabels, tuning, isIsolated, neutralZones, holdCityNeutralLabel)
	conns := this.createConnections(playerLabels, orderedLabels, tuning, isIsolated, neutralZones)
	if configuration.RandomPortals {
		conns = append(conns,
			this.CreateRandomPortalConnections(
				playerLabels, orderedLabels, tuning, configuration.MaxPortalConnections)...)
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
	neutralZones neutral_zone.Plans,
	holdCityNeutralLabel string) []entities.Zone {
	labelCount := len(orderedLabels)

	ringConnRight := make([]string, labelCount)
	ringConnLeft := make([]string, labelCount)
	for i := range labelCount {
		next := (i + 1) % labelCount
		if isIsolated &&
			slices.Contains(playerLabels, orderedLabels[i]) &&
			slices.Contains(playerLabels, orderedLabels[next]) {
			continue
		}
		name := fmt.Sprintf("Ring-%s-%s", orderedLabels[i], orderedLabels[next])
		ringConnRight[i] = name
		ringConnLeft[next] = name
	}

	var zones []entities.Zone
	for i, label := range orderedLabels {
		var connNames []string
		if ringConnLeft[i] != "" {
			connNames = append(connNames, ringConnLeft[i])
		}
		if ringConnRight[i] != "" && ringConnRight[i] != ringConnLeft[i] {
			connNames = append(connNames, ringConnRight[i])
		}
		playerIndex := slices.Index(playerLabels, label)
		zones = append(zones, this.CreateClusterZone(
			configuration, label, connNames, playerIndex, playerIndex >= 0,
			label == holdCityNeutralLabel, tuning, neutralZones))
	}
	return zones
}

func (this *RingTopologyService) createConnections(
	playerLabels, orderedLabels []string,
	tuning models.GenerationTuning,
	isIsolated bool,
	neutralZones neutral_zone.Plans) []entities.Connection {
	count := len(orderedLabels)
	if count < 2 {
		return nil
	}

	var connections []entities.Connection
	for i := range count {
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
			WithGuardWeeklyIncrement(common_connections.GetGuardWeeklyIncrements().Standard).
			WithGuardMatchGroup(fmt.Sprintf("ring_guard_%s_%s", labelFrom, labelTo)).
			Build())
	}
	return connections
}
