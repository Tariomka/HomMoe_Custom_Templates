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

type ChainTopologyService struct {
	base.TopologyBase
}

func NewChainTopologyService(
	zoneFactory *zone_services.ZoneFactory,
	roadFactory *zone_services.RoadFactory,
	zoneLabelProvider zone_services.IZoneLabelProvider,
	connectionService *base.TopologyConnectionService,
) *ChainTopologyService {
	return &ChainTopologyService{
		TopologyBase: base.NewTopologyBase(zoneFactory, roadFactory, zoneLabelProvider, connectionService),
	}
}

func (this *ChainTopologyService) CreateTopologyVariant(
	configuration config.GeneratorConfig,
	playerLetters []string,
	neutralZones neutral_zone.Plans,
	tuning models.GenerationTuning,
	holdCityNeutralLetter string) entities.Variant {
	orderedLabels := this.ZoneLabelProvider.CreateOrderedZoneLabels(configuration, playerLetters, neutralZones, false)
	isIsolated := configuration.NoDirectPlayerConnections && len(playerLetters) > 1
	connNames := this.createConnectionNames(playerLetters, orderedLabels, isIsolated)

	zones := this.createZones(
		configuration, playerLetters, orderedLabels, tuning, neutralZones, holdCityNeutralLetter, connNames)
	conns := this.createConnections(playerLetters, orderedLabels, tuning, neutralZones, connNames)
	if configuration.RandomPortals {
		conns = append(conns,
			this.CreateRandomPortalConnections(
				playerLetters, orderedLabels, tuning, configuration.MaxPortalConnections)...)
	}
	if isIsolated {
		conns = append(conns, this.CreateMissingPlayerConnections(playerLetters, zones, conns, tuning)...)
	}
	return this.CreateVariant(playerLetters, orderedLabels[0], len(orderedLabels), zones, conns)
}

func (this *ChainTopologyService) createConnectionNames(
	playerLetters, orderedLabels []string,
	isIsolated bool) []string {
	labelCount := len(orderedLabels)

	connNames := make([]string, labelCount-1)
	for i := range labelCount - 1 {
		if isIsolated &&
			slices.Contains(playerLetters, orderedLabels[i]) &&
			slices.Contains(playerLetters, orderedLabels[i+1]) {
			continue
		}

		connNames[i] = fmt.Sprintf("Chain-%s-%s", orderedLabels[i], orderedLabels[i+1])
	}
	return connNames
}

func (this *ChainTopologyService) createZones(
	configuration config.GeneratorConfig,
	playerLabels, orderedLabels []string,
	tuning models.GenerationTuning,
	neutralZones neutral_zone.Plans,
	holdCityNeutralLabel string,
	connectionNames []string) []entities.Zone {
	labelCount := len(orderedLabels)

	var zones []entities.Zone
	for index, label := range orderedLabels {
		var tempConnectionNames []string
		if index > 0 && connectionNames[index-1] != "" {
			tempConnectionNames = append(tempConnectionNames, connectionNames[index-1])
		}
		if index < labelCount-1 && connectionNames[index] != "" {
			tempConnectionNames = append(tempConnectionNames, connectionNames[index])
		}
		playerIndex := slices.Index(playerLabels, label)
		zones = append(zones, this.CreateClusterZone(
			configuration, label, tempConnectionNames, playerIndex, playerIndex >= 0,
			label == holdCityNeutralLabel, tuning, neutralZones))
	}
	return zones
}

func (this *ChainTopologyService) createConnections(
	playerLabels, orderedLabels []string,
	tuning models.GenerationTuning,
	neutralZones neutral_zone.Plans,
	connectionNames []string) []entities.Connection {
	labelCount := len(orderedLabels)

	var connections []entities.Connection
	for i := range labelCount - 1 {
		if connectionNames[i] == "" {
			continue
		}

		labelFrom := orderedLabels[i]
		labelTo := orderedLabels[i+1]
		zoneFrom := this.ZoneLabelProvider.CreateZoneName(labelFrom, playerLabels)
		zoneTo := this.ZoneLabelProvider.CreateZoneName(labelTo, playerLabels)
		connections = append(connections, variant_content.NewConnectionBuilder().
			WithName(connectionNames[i]).
			WithFrom(zoneFrom).
			WithTo(zoneTo).
			WithConnectionTypeDirect().
			WithGuardZone(zoneFrom).
			WithSimTurnSquad().
			WithGuardValue(this.GetBorderGuardValue(labelFrom, labelTo, playerLabels, neutralZones, tuning)).
			WithGuardWeeklyIncrement(common_connections.GetGuardWeeklyIncrements().Standard).
			WithGuardMatchGroup(fmt.Sprintf("chain_guard_%s_%s", labelFrom, labelTo)).
			Build())
	}
	return connections
}
