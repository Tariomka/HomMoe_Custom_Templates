package tournament_variant

import (
	"fmt"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_connections"
	"github.com/Tariomka/hommoe_custom_templates/internal/common/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones/zone_interfaces"
)

type ChainClusterService struct {
	base.TopologyBase
}

func NewChainClusterService(
	zoneFactory zone_interfaces.IZoneFactory,
	roadFactory zone_interfaces.IRoadFactory,
	zoneLabelProvider zone_interfaces.IZoneLabelProvider,
	connectionService base.ITopologyConnectionService,
) *ChainClusterService {
	return &ChainClusterService{
		TopologyBase: base.NewTopologyBase(zoneFactory, roadFactory, zoneLabelProvider, connectionService),
	}
}

func (this *ChainClusterService) CreateClusterVariant(
	configuration config.GeneratorConfig,
	tuning models.GenerationTuning,
	allNeutralZonePlans, playerNeutralZonePlans neutral_zone.Plans,
	playerIndex int,
	playerLabel string) ([]entities.Zone, []entities.Connection) {
	chainLabels := append([]string{playerLabel},
		linq.FromSlice(playerNeutralZonePlans).
			SelectString(func(x neutral_zone.Plan) string { return x.Label }).
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
	allNeutralZonePlans neutral_zone.Plans,
	playerIndex int) []entities.Zone {
	var zones []entities.Zone
	for index, label := range chainLabels {
		var myConns []string
		if index > 0 {
			myConns = append(myConns, connectionNames[index-1])
		}
		if index < len(connectionNames) {
			myConns = append(myConns, connectionNames[index])
		}
		zones = append(zones, this.CreateClusterZone(
			configuration, label, myConns, playerIndex, index == 0, false, tuning, allNeutralZonePlans))
	}
	return zones
}

func (this *ChainClusterService) createConnections(
	chainLabels, connectionNames []string,
	tuning models.GenerationTuning,
	allNeutralZonePlans neutral_zone.Plans,
	playerLabel string) []entities.Connection {
	var connections []entities.Connection
	for index, name := range connectionNames {
		labelFrom := chainLabels[index]
		labelTo := chainLabels[index+1]
		connectionBuilder := variant_content.NewConnectionBuilder().
			WithName(name).
			WithTo(constants.NeutralZonePrefix + labelTo).
			WithConnectionTypeDirect().
			WithGuardValue(this.GetBorderGuardValue(labelFrom, labelTo, []string{playerLabel}, allNeutralZonePlans, tuning)).
			WithGuardWeeklyIncrement(common_connections.GetGuardWeeklyIncrements().Standard).
			WithGuardMatchGroup(fmt.Sprintf("tourney_guard_%s_%s", labelFrom, labelTo))

		if index > 0 {
			connectionBuilder = connectionBuilder.WithFrom(constants.NeutralZonePrefix + labelFrom)
		} else {
			connectionBuilder = connectionBuilder.WithFrom(constants.PlayerZonePrefix + labelFrom)
		}

		connections = append(connections, connectionBuilder.Build())
	}
	return connections
}
