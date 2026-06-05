package topology

import (
	"fmt"
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/builders/variant_content"
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
	neutralZones models.NeutralZonePlans,
	tuning models.GenerationTuning,
	holdCityNeutralLetter string) template.Variant {
	orderedLabels := this.zoneLabelProvider.CreateOrderedZoneLabels(configuration, playerLetters, neutralZones, false)
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
	playerLetters, orderedLabels []string,
	tuning models.GenerationTuning,
	neutralZones models.NeutralZonePlans,
	holdCityNeutralLabel string,
	connectionNames []string) []template.Zone {
	labelCount := len(orderedLabels)

	var zones []template.Zone
	for i, label := range orderedLabels {
		var tempConnectionNames []string
		if i > 0 && connectionNames[i-1] != "" {
			tempConnectionNames = append(tempConnectionNames, connectionNames[i-1])
		}
		if i < labelCount-1 && connectionNames[i] != "" {
			tempConnectionNames = append(tempConnectionNames, connectionNames[i])
		}
		if playerIndex := slices.Index(playerLetters, label); playerIndex >= 0 {
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
	connectionNames []string) []template.Connection {
	labelCount := len(orderedLabels)

	var conns []template.Connection
	for i := range labelCount - 1 {
		if connectionNames[i] == "" {
			continue
		}

		from := orderedLabels[i]
		to := orderedLabels[i+1]
		fromZone := createZoneName(from, playerLabels)
		toZone := createZoneName(to, playerLabels)
		conns = append(conns, variant_content.NewConnectionBuilder().
			WithName(connectionNames[i]).
			WithFrom(fromZone).
			WithTo(toZone).
			WithConnectionTypeDirect().
			WithGuardZone(fromZone).
			WithSimTurnSquad().
			WithGuardValue(this.GetBorderGuardValue(from, to, playerLabels, neutralZones, tuning)).
			WithGuardWeeklyIncrement(0.15).
			WithGuardMatchGroup(fmt.Sprintf("chain_guard_%s_%s", from, to)).
			Build())
	}
	return conns
}
