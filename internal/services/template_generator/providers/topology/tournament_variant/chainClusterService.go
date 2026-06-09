package tournament_variant

import (
	"fmt"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/builders/variant_content"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
)

type ChainClusterService struct {
	base.TopologyBase
}

func NewChainClusterService() *ChainClusterService {
	return &ChainClusterService{
		TopologyBase: base.NewTopologyBase(),
	}
}

func (this *ChainClusterService) CreateClusterVariant(
	configuration config.GeneratorConfig,
	tuning models.GenerationTuning,
	allNeutralZonePlans, playerNeutralZonePlans models.NeutralZonePlans,
	playerIndex int,
	playerLabel string) ([]template.Zone, []template.Connection) {
	chainLabels := append([]string{playerLabel},
		linq.FromSlice(playerNeutralZonePlans).
			SelectString(func(x models.NeutralZonePlan) string { return x.Label }).
			ToSlice()...)
	connectionNames := make([]string, len(chainLabels)-1)
	for index := range connectionNames {
		connectionNames[index] = fmt.Sprintf("Tourney-%s-%s", chainLabels[index], chainLabels[index+1])
	}

	zones := this.createZones(configuration, chainLabels, connectionNames, tuning, allNeutralZonePlans, playerIndex)
	connections := this.createConnections(chainLabels, connectionNames, tuning, allNeutralZonePlans, playerLabel)
	return zones, connections
}

func (this *ChainClusterService) createZones(
	configuration config.GeneratorConfig,
	chainLabels, connectionNames []string,
	tuning models.GenerationTuning,
	allNeutralZonePlans models.NeutralZonePlans,
	playerIndex int) []template.Zone {
	var zones []template.Zone
	for index, label := range chainLabels {
		var myConns []string
		if index > 0 {
			myConns = append(myConns, connectionNames[index-1])
		}
		if index < len(connectionNames) {
			myConns = append(myConns, connectionNames[index])
		}
		if index == 0 {
			zones = append(zones,
				this.CreateSpawnZone(
					label, fmt.Sprintf("Player%d", playerIndex+1), myConns,
					configuration.ZoneConfiguration.PlayerZoneCastles, configuration.MatchPlayerCastleFactions,
					configuration.ZoneConfiguration.Advanced.PlayerZoneSize, configuration.SpawnRemoteFootholds,
					configuration.GenerateRoads, tuning))
		} else {
			zones = append(zones,
				this.CreateNeutralZone(
					linq.FromSlice(allNeutralZonePlans).FirstOrDefault(func(x models.NeutralZonePlan) bool { return x.Label == label }),
					myConns, configuration.ZoneConfiguration.Advanced.NeutralZoneSize, configuration.SpawnRemoteFootholds,
					configuration.GenerateRoads, tuning, false))
		}
	}
	return zones
}

func (this *ChainClusterService) createConnections(
	chainLabels, connectionNames []string,
	tuning models.GenerationTuning,
	allNeutralZonePlans models.NeutralZonePlans,
	playerLabel string) []template.Connection {
	var connections []template.Connection
	for index, name := range connectionNames {
		labelFrom := chainLabels[index]
		labelTo := chainLabels[index+1]
		connectionBuilder := variant_content.NewConnectionBuilder().
			WithName(name).
			WithTo("Neutral-" + labelTo).
			WithConnectionTypeDirect().
			WithGuardValue(this.GetBorderGuardValue(labelFrom, labelTo, []string{playerLabel}, allNeutralZonePlans, tuning)).
			WithGuardWeeklyIncrement(0.15).
			WithGuardMatchGroup(fmt.Sprintf("tourney_guard_%s_%s", labelFrom, labelTo))

		if index > 0 {
			connectionBuilder = connectionBuilder.WithFrom("Neutral-" + labelFrom)
		} else {
			connectionBuilder = connectionBuilder.WithFrom("Spawn-" + labelFrom)
		}

		connections = append(connections, connectionBuilder.Build())
	}
	return connections
}
