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

type RingClusterService struct{ base.TopologyBase }

func NewRingClusterService(
	zoneFactory zone_interfaces.IZoneFactory,
	roadFactory zone_interfaces.IRoadFactory,
	zoneLabelProvider zone_interfaces.IZoneLabelProvider,
	connectionService base.ITopologyConnectionService) *RingClusterService {
	return &RingClusterService{
		TopologyBase: base.NewTopologyBase(zoneFactory, roadFactory, zoneLabelProvider, connectionService),
	}
}

func (this *RingClusterService) CreateClusterVariant(
	configuration config.GeneratorConfig,
	tuning models.GenerationTuning,
	allNeutralZonePlans, playerNeutralZonePlans neutral_zone.Plans,
	playerIndex int,
	playerLabel string) ([]entities.Zone, []entities.Connection) {
	ringLabels := this.createLabels(playerNeutralZonePlans, playerLabel)
	ringCount := len(ringLabels)
	if ringCount < 2 {
		return []entities.Zone{this.createSinglePlayerZone(configuration, playerLabel, playerIndex, tuning)}, nil
	}

	connNames := make([]string, ringCount)
	for index, label := range ringLabels {
		nextIndex := (index + 1) % ringCount
		connNames[index] = constants.GetTournamentRingConnectionNameFor(label, ringLabels[nextIndex])
	}

	zones := this.createZones(configuration, ringLabels, connNames, tuning, allNeutralZonePlans, playerIndex)
	connections := this.createConnections(ringLabels, connNames, tuning, allNeutralZonePlans, playerLabel)
	return zones, connections
}

func (this *RingClusterService) createLabels(playerNeutralZonePlans neutral_zone.Plans, playerLabel string) []string {
	sortedNeutralZonePlans := neutral_zone.Plans{}
	sortedNeutralZonePlans.AddPlans(playerNeutralZonePlans...)
	sortedNeutralZonePlans.SortByBalanceScoreAscending()

	zoneCount := len(sortedNeutralZonePlans)
	orderedNeutralZonePlans := make(neutral_zone.Plans, zoneCount)
	lowIndex, highIndex := 0, zoneCount-1
	for index, zonePlan := range sortedNeutralZonePlans {
		if index%2 == 0 {
			orderedNeutralZonePlans[lowIndex] = zonePlan
			lowIndex++
		} else {
			orderedNeutralZonePlans[highIndex] = zonePlan
			highIndex--
		}
	}

	return append([]string{playerLabel},
		linq.FromSlice(orderedNeutralZonePlans).
			Select(func(x neutral_zone.Plan) string { return x.Label }).
			ToSlice()...)
}

func (this *RingClusterService) createZones(
	configuration config.GeneratorConfig,
	ringLabels, connectionNames []string,
	tuning models.GenerationTuning,
	allNeutralZonePlans neutral_zone.Plans,
	playerIndex int) []entities.Zone {
	count := len(ringLabels)

	var zones []entities.Zone
	for index, label := range ringLabels {
		prev := (index - 1 + count) % count
		seen := map[string]bool{}
		var myConns []string
		for _, name := range []string{connectionNames[prev], connectionNames[index]} {
			if !seen[name] {
				seen[name] = true
				myConns = append(myConns, name)
			}
		}
		zones = append(zones, this.CreateClusterZone(
			configuration, label, myConns, playerIndex, index == 0, false, tuning, allNeutralZonePlans))
	}
	return zones
}

func (this *RingClusterService) createSinglePlayerZone(
	configuration config.GeneratorConfig,
	playerLabel string,
	playerIndex int,
	tuning models.GenerationTuning) entities.Zone {
	return this.CreateClusterZone(configuration, playerLabel, nil, playerIndex, true, false, tuning, nil)
}

func (this *RingClusterService) createConnections(
	ringLabels, connectionNames []string,
	tuning models.GenerationTuning,
	allNeutralZonePlans neutral_zone.Plans,
	playerLabel string) []entities.Connection {
	ringCount := len(ringLabels)

	var connections []entities.Connection
	for currentIndex := range ringCount {
		nextIndex := (currentIndex + 1) % ringCount
		labelFrom := ringLabels[currentIndex]
		labelTo := ringLabels[nextIndex]

		connectionBuilder := variant_content.NewConnectionBuilder().
			WithName(connectionNames[currentIndex]).
			WithConnectionTypeDirect().
			WithSimTurnSquad().
			WithGuardValue(this.GetBorderGuardValue(labelFrom, labelTo, []string{playerLabel}, allNeutralZonePlans, tuning)).
			WithGuardWeeklyIncrement(common_connections.GetGuardWeeklyIncrements().Standard).
			WithGuardMatchGroup(fmt.Sprintf("tourney_ring_guard_%s_%s", labelFrom, labelTo))

		if currentIndex != 0 {
			zoneFrom := constants.GetNeutralZoneNameFor(labelFrom)
			connectionBuilder.WithFrom(zoneFrom).WithGuardZone(zoneFrom)
		} else {
			zoneFrom := constants.GetPlayerZoneNameFor(labelFrom)
			connectionBuilder.WithFrom(zoneFrom).WithGuardZone(zoneFrom)
		}

		if nextIndex != 0 {
			zoneTo := constants.GetNeutralZoneNameFor(labelTo)
			connectionBuilder.WithTo(zoneTo)
		} else {
			zoneTo := constants.GetPlayerZoneNameFor(labelTo)
			connectionBuilder.WithTo(zoneTo)
		}

		connections = append(connections, connectionBuilder.Build())
	}
	return connections
}
