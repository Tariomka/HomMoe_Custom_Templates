package topology

import (
	"fmt"
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/builders/variant_content"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
)

type ChainTopologyService struct {
	base.TopologyBase
}

func NewChainTopologyService() *ChainTopologyService {
	return &ChainTopologyService{
		TopologyBase: base.NewTopologyBase(),
	}
}

func (this *ChainTopologyService) CreateTopologyVariant(
	configuration config.GeneratorConfig,
	playerLetters []string,
	neutralZones models.NeutralZonePlans,
	tuning models.GenerationTuning,
	holdCityNeutralLetter string) entities.Variant {
	orderedLabels := this.ZoneLabelProvider.CreateOrderedZoneLabels(configuration, playerLetters, neutralZones, false)
	isIsolated := configuration.NoDirectPlayerConnections && len(playerLetters) > 1
	connNames := this.createConnectionNames(playerLetters, orderedLabels, isIsolated)

	zones := this.createZones(configuration, playerLetters, orderedLabels, tuning, neutralZones, holdCityNeutralLetter, connNames)
	conns := this.createConnections(playerLetters, orderedLabels, tuning, neutralZones, connNames)
	if configuration.RandomPortals {
		conns = append(conns, this.CreateRandomPortalConnections(playerLetters, orderedLabels, tuning, configuration.MaxPortalConnections)...)
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
		if isIsolated && slices.Contains(playerLetters, orderedLabels[i]) && slices.Contains(playerLetters, orderedLabels[i+1]) {
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
	neutralZones models.NeutralZonePlans,
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
		if playerIndex := slices.Index(playerLabels, label); playerIndex >= 0 {
			zones = append(zones,
				this.CreateSpawnZone(
					label, fmt.Sprintf("Player%d", playerIndex+1), tempConnectionNames, configuration.ZoneConfiguration.PlayerZoneCastles,
					configuration.MatchPlayerCastleFactions, configuration.ZoneConfiguration.Advanced.PlayerZoneSize,
					configuration.SpawnRemoteFootholds, configuration.GenerateRoads, tuning))
		} else {
			zones = append(zones,
				this.CreateNeutralZone(
					linq.FromSlice(neutralZones).FirstOrDefault(func(x models.NeutralZonePlan) bool { return x.Label == label }),
					tempConnectionNames, configuration.ZoneConfiguration.Advanced.NeutralZoneSize,
					configuration.SpawnRemoteFootholds, configuration.GenerateRoads, tuning, label == holdCityNeutralLabel))
		}
	}
	return zones
}

func (this *ChainTopologyService) createConnections(
	playerLabels, orderedLabels []string,
	tuning models.GenerationTuning,
	neutralZones models.NeutralZonePlans,
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
			WithGuardWeeklyIncrement(0.15).
			WithGuardMatchGroup(fmt.Sprintf("chain_guard_%s_%s", labelFrom, labelTo)).
			Build())
	}
	return connections
}
