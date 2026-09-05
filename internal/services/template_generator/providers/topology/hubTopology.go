package topology

import (
	"fmt"
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_connections"
	"github.com/Tariomka/hommoe_custom_templates/internal/common/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones/zone_interfaces"
)

type HubTopologyService struct {
	base.TopologyBase
}

func NewHubTopologyService(
	zoneFactory zone_interfaces.IZoneFactory,
	roadFactory zone_interfaces.IRoadFactory,
	zoneLabelProvider zone_interfaces.IZoneLabelProvider,
	connectionService base.ITopologyConnectionService) *HubTopologyService {
	return &HubTopologyService{
		TopologyBase: base.NewTopologyBase(zoneFactory, roadFactory, zoneLabelProvider, connectionService),
	}
}

func (this *HubTopologyService) CreateTopologyVariant(
	configuration config.GeneratorConfig,
	playerLabels []string,
	neutralZones neutral_zone.Plans,
	tuning models.GenerationTuning,
	_ string) template_model.Variant {
	outerLabels := this.createOuterLabels(configuration, playerLabels, neutralZones)

	zones := this.createZones(configuration, playerLabels, outerLabels, tuning, neutralZones)
	conns := this.createConnections(
		playerLabels, outerLabels, tuning, configuration.NoDirectPlayerConnections, neutralZones)
	if configuration.RandomPortals {
		conns = append(conns,
			this.CreateRandomPortalConnections(
				playerLabels, outerLabels, tuning, configuration.MaxPortalConnections, neutralZones)...)
	}
	return this.CreateVariant(playerLabels, outerLabels[0], len(outerLabels)+1, zones, conns)
}

func (this *HubTopologyService) createOuterLabels(
	configuration config.GeneratorConfig,
	playerLabels []string,
	neutralZones neutral_zone.Plans) []string {
	if configuration.Topology == config.TopologyCircles {
		return this.ZoneLabelProvider.CreateBalancedChainZoneLabels(playerLabels, neutralZones)
	}

	return append(playerLabels,
		linq.FromSlice(neutralZones).Select(func(nz neutral_zone.Plan) string { return nz.Label }).ToSlice()...)
}

func (this *HubTopologyService) createZones(
	configuration config.GeneratorConfig,
	playerLabels, outerLabels []string,
	tuning models.GenerationTuning,
	neutralZones neutral_zone.Plans) []template_model.Zone {
	hubConns := make([]string, len(outerLabels))
	for index, label := range outerLabels {
		hubConns[index] = constants.GetHubSpokeConnectionNameFor(label)
	}
	hubContentName := ""
	if len(configuration.HubZoneMandatoryContent) > 0 {
		hubContentName = constants.HubContentName
	}
	zones := []template_model.Zone{
		this.CreateHubZone(
			constants.HubZoneName, hubConns, tuning, configuration.IsHubCityToHold(),
			configuration.ZoneConfiguration.HubZoneSize,
			configuration.ZoneConfiguration.Advanced.HubZoneCastles,
			configuration.GenerateRoads, hubContentName),
	}

	for _, label := range outerLabels {
		spokeConnectionNames := []string{constants.GetHubSpokeConnectionNameFor(label)}
		playerIndex := slices.Index(playerLabels, label)
		zones = append(zones, this.CreateClusterZone(
			configuration, label, spokeConnectionNames, playerIndex, playerIndex >= 0, false, tuning, neutralZones))
	}
	return zones
}

func (this *HubTopologyService) createConnections(
	playerLabels, outerLabels []string,
	tuning models.GenerationTuning,
	isIsolated bool,
	neutralZones neutral_zone.Plans) []template_model.Connection {
	var connections []template_model.Connection
	for index, label := range outerLabels {
		hubGuard := this.GetBorderGuardValue(
			constants.HubZoneName, label, playerLabels, neutralZones, tuning)
		outerZone := this.ZoneLabelProvider.CreateZoneName(label, playerLabels)
		connections = append(connections,
			variant_content.NewConnectionBuilder().
				WithName(constants.GetHubSpokeConnectionNameFor(label)).
				WithFrom(constants.HubZoneName).
				WithTo(outerZone).
				WithConnectionTypeDirect().
				WithGuardZone(constants.HubZoneName).
				WithSimTurnSquad().
				WithGuardValue(hubGuard).
				WithGuardWeeklyIncrement(common_connections.GetGuardWeeklyIncrements().Standard).
				WithGuardMatchGroup("hub_guard_"+label).
				Build(),
			variant_content.NewConnectionBuilder().
				WithFrom(constants.HubZoneName).
				WithTo(outerZone).
				WithConnectionTypeDirect().
				WithGuardZone(constants.HubZoneName).
				WithSimTurnSquad().
				WithGuardValue(hubGuard).
				WithGuardWeeklyIncrement(common_connections.GetGuardWeeklyIncrements().Standard).
				WithGuardMatchGroup(fmt.Sprintf("hub_guard_%s_%d", label, 1)).
				Build())

		labelTo := outerLabels[(index+1)%len(outerLabels)]
		if isIsolated && slices.Contains(playerLabels, label) && slices.Contains(playerLabels, labelTo) {
			continue
		}
		connections = append(connections, variant_content.NewConnectionBuilder().
			WithName(constants.GetPseudoConnectionNameFor(label, labelTo)).
			WithFrom(this.ZoneLabelProvider.CreateZoneName(label, playerLabels)).
			WithTo(this.ZoneLabelProvider.CreateZoneName(labelTo, playerLabels)).
			WithConnectionTypeProximity().
			Build())
	}
	return connections
}
