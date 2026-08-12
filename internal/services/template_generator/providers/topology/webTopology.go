package topology

import (
	"fmt"
	"strings"

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

type SharedWebTopologyService struct {
	base.TopologyBase
}

func NewSharedWebTopologyService(
	zoneFactory zone_interfaces.IZoneFactory,
	roadFactory zone_interfaces.IRoadFactory,
	zoneLabelProvider zone_interfaces.IZoneLabelProvider,
	connectionService base.ITopologyConnectionService,
) *SharedWebTopologyService {
	return &SharedWebTopologyService{
		TopologyBase: base.NewTopologyBase(zoneFactory, roadFactory, zoneLabelProvider, connectionService),
	}
}

func (this *SharedWebTopologyService) CreateTopologyVariant(
	configuration config.GeneratorConfig,
	playerLabels []string,
	neutralZones neutral_zone.Plans,
	tuning models.GenerationTuning,
	holdCityNeutralLetter string) entities.Variant {
	neutralLabels := this.createLabels(playerLabels, neutralZones, configuration.Topology == config.TopologyCircles)
	playerSpokes, neutralSpokes := this.createSpokes(playerLabels, neutralLabels)
	neutralConnNames := this.createRingConnectionNames(neutralLabels)

	zones := this.createZones(
		configuration, playerLabels, neutralLabels, tuning, neutralZones,
		holdCityNeutralLetter, playerSpokes, neutralSpokes, neutralConnNames)
	conns := this.createConnections(playerLabels, neutralLabels, tuning, neutralZones, playerSpokes, neutralConnNames)
	if configuration.RandomPortals {
		allLabels := append(append([]string{}, playerLabels...), neutralLabels...)
		conns = append(conns,
			this.CreateRandomPortalConnections(
				playerLabels, allLabels, tuning, configuration.MaxPortalConnections, neutralZones)...)
	}
	if configuration.NoDirectPlayerConnections && len(playerLabels) > 1 {
		conns = append(conns, this.CreateMissingPlayerConnections(playerLabels, zones, conns, tuning)...)
	}
	return this.CreateVariant(playerLabels, playerLabels[0], len(zones), zones, conns)
}

func (this *SharedWebTopologyService) createLabels(
	playerLabels []string,
	neutralZones neutral_zone.Plans,
	isBalanced bool) []string {
	var neutrals []string
	if isBalanced {
		neutrals = this.ZoneLabelProvider.CreateBalancedNeutralRingZoneLabels(neutralZones, len(playerLabels))
	} else {
		for _, zonePlan := range neutralZones {
			neutrals = append(neutrals, zonePlan.Label)
		}
	}
	return neutrals
}

func (this *SharedWebTopologyService) createSpokes(
	playerLabels, neutralLabels []string) (playerSpokes, neutralSpokes map[string][]string) {
	playerCount := len(playerLabels)
	neutralCount := len(neutralLabels)

	addSpoke := func(playerLabel, neutralLabel string) {
		connectionName := constants.GetWebConnectionNameFor(playerLabel, neutralLabel)
		playerSpokes[playerLabel] = append(playerSpokes[playerLabel], connectionName)
		neutralSpokes[neutralLabel] = append(neutralSpokes[neutralLabel], connectionName)
	}

	playerSpokes = make(map[string][]string, playerCount)
	for _, label := range playerLabels {
		playerSpokes[label] = nil
	}
	neutralSpokes = make(map[string][]string, neutralCount)
	for _, label := range neutralLabels {
		neutralSpokes[label] = nil
	}

	for index, label := range playerLabels {
		n1 := (index * neutralCount / playerCount) % neutralCount
		n2 := ((index * neutralCount / playerCount) + 1) % neutralCount
		addSpoke(label, neutralLabels[n1])
		if n1 != n2 {
			addSpoke(label, neutralLabels[n2])
		}
	}

	return playerSpokes, neutralSpokes
}

func (this *SharedWebTopologyService) createRingConnectionNames(neutralLabels []string) []string {
	neutralCount := len(neutralLabels)
	neutralRingConnNames := make([]string, neutralCount)
	for i, label := range neutralLabels {
		neutralRingConnNames[i] = constants.GetNeutralRingConnectionNameFor(label, neutralLabels[(i+1)%neutralCount])
	}
	return neutralRingConnNames
}

func (this *SharedWebTopologyService) createZones(
	configuration config.GeneratorConfig,
	playerLabels, neutralLabels []string,
	tuning models.GenerationTuning,
	neutralZones neutral_zone.Plans,
	holdCityNeutralLabel string,
	playerSpokes, neutralSpokes map[string][]string,
	connectionNames []string) []entities.Zone {
	neutralCount := len(neutralLabels)

	var zones []entities.Zone
	for i, label := range neutralLabels {
		var neutralConnNames []string
		if neutralCount > 1 {
			neutralConnNames = append(neutralConnNames,
				connectionNames[(i-1+neutralCount)%neutralCount],
				connectionNames[i])
		}
		neutralConnNames = linq.FromSlice(append(neutralConnNames, neutralSpokes[label]...)).Distinct().ToSlice()
		zonePlan := linq.FromSlice(neutralZones).
			FirstOrDefault(func(x neutral_zone.Plan) bool { return x.Label == label })
		zone := this.CreateClusterZone(
			configuration, label, neutralConnNames, 0, false, label == holdCityNeutralLabel, tuning, neutralZones)
		if zonePlan.CastleCount == 0 && tuning.AbandonedOutpostCount == 0 {
			zone.Roads = this.CreateConnectorZoneRoads(neutralConnNames, configuration.GenerateRoads)
		}
		zones = append(zones, zone)
	}

	for i, label := range playerLabels {
		zones = append(zones, this.CreateClusterZone(
			configuration, label, playerSpokes[label], i, true, false, tuning, neutralZones))
	}

	return zones
}

func (this *SharedWebTopologyService) createConnections(
	playerLabels, neutralLabels []string,
	tuning models.GenerationTuning,
	neutralZones neutral_zone.Plans,
	playerSpokes map[string][]string,
	connectionNames []string) []entities.Connection {
	neutralCount := len(neutralLabels)

	var connections []entities.Connection
	for _, label := range playerLabels {
		for _, connectionName := range playerSpokes[label] {
			nextLabel := strings.Split(connectionName, "-")[2]
			connections = append(connections, variant_content.NewConnectionBuilder().
				WithName(connectionName).
				WithFrom(constants.GetPlayerZoneNameFor(label)).
				WithTo(constants.GetNeutralZoneNameFor(nextLabel)).
				WithConnectionTypeDirect().
				WithGuardZone(constants.GetNeutralZoneNameFor(nextLabel)).
				WithSimTurnSquad().
				WithGuardValue(this.GetBorderGuardValue(label, nextLabel, playerLabels, neutralZones, tuning)).
				WithGuardWeeklyIncrement(common_connections.GetGuardWeeklyIncrements().Standard).
				WithGuardMatchGroup(fmt.Sprintf("web_guard_%s_%s", label, nextLabel)).
				Build())
		}
	}

	if neutralCount < 2 {
		return connections
	}

	for i, label := range neutralLabels {
		next := (i + 1) % neutralCount
		nextLabel := neutralLabels[next]
		connections = append(connections, variant_content.NewConnectionBuilder().
			WithName(connectionNames[i]).
			WithFrom(constants.GetNeutralZoneNameFor(label)).
			WithTo(constants.GetNeutralZoneNameFor(nextLabel)).
			WithConnectionTypeDirect().
			WithGuardZone(constants.GetNeutralZoneNameFor(label)).
			WithSimTurnSquad().
			WithGuardValue(this.GetBorderGuardValue(label, nextLabel, playerLabels, neutralZones, tuning)).
			WithGuardWeeklyIncrement(common_connections.GetGuardWeeklyIncrements().Standard).
			WithGuardMatchGroup(fmt.Sprintf("nring_guard_%s_%s", label, nextLabel)).
			Build())
	}
	return connections
}
