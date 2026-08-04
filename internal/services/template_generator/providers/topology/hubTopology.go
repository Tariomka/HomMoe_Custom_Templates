package topology

import (
	"fmt"
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
	zone_services "github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
)

type HubTopologyService struct {
	base.TopologyBase
}

func NewHubTopologyService(
	zoneFactory *zone_services.ZoneFactory,
	roadFactory *zone_services.RoadFactory,
) *HubTopologyService {
	return &HubTopologyService{
		TopologyBase: base.NewTopologyBase(zoneFactory, roadFactory),
	}
}

func (this *HubTopologyService) CreateTopologyVariant(
	configuration config.GeneratorConfig,
	playerLabels []string,
	neutralZones neutral_zone.Plans,
	tuning models.GenerationTuning,
	_ string) entities.Variant {
	outerLabels := this.createOuterLabels(configuration, playerLabels, neutralZones)

	zones := this.createZones(configuration, playerLabels, outerLabels, tuning, neutralZones)
	conns := this.createConnections(
		playerLabels, outerLabels, tuning, configuration.NoDirectPlayerConnections, neutralZones)
	if configuration.RandomPortals {
		conns = append(conns,
			this.CreateRandomPortalConnections(
				playerLabels, outerLabels, tuning, configuration.MaxPortalConnections)...)
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
		linq.FromSlice(neutralZones).
			SelectString(func(nz neutral_zone.Plan) string { return nz.Label }).
			ToSlice()...)
}

func (this *HubTopologyService) createZones(
	configuration config.GeneratorConfig,
	playerLabels, outerLabels []string,
	tuning models.GenerationTuning,
	neutralZones neutral_zone.Plans) []entities.Zone {
	hubConns := make([]string, len(outerLabels))
	for index, label := range outerLabels {
		hubConns[index] = constants.HubZonePrefix + label
	}
	hubContentName := ""
	if len(configuration.HubZoneMandatoryContent) > 0 {
		hubContentName = "mandatory_content_hub"
	}
	zones := []entities.Zone{
		this.CreateHubZone(
			constants.HubZoneName, hubConns, tuning, configuration.IsHubCityToHold(),
			configuration.ZoneConfiguration.HubZoneSize,
			configuration.ZoneConfiguration.Advanced.HubZoneCastles,
			configuration.GenerateRoads, hubContentName),
	}

	for _, label := range outerLabels {
		spokeConnectionNames := []string{constants.HubZonePrefix + label}
		if playerIndex := slices.Index(playerLabels, label); playerIndex >= 0 {
			zones = append(zones,
				this.CreateSpawnZone(
					label, fmt.Sprintf("Player%d", playerIndex+1), spokeConnectionNames,
					configuration.ZoneConfiguration.PlayerZoneCastles, configuration.MatchPlayerCastleFactions,
					configuration.ZoneConfiguration.PlayerZoneSize, tuning.RemoteFootholdCount,
					configuration.GenerateRoads, tuning))
		} else {
			zones = append(zones,
				this.CreateNeutralZone(
					linq.FromSlice(neutralZones).
						FirstOrDefault(func(x neutral_zone.Plan) bool { return x.Label == label }),
					spokeConnectionNames, configuration.ZoneConfiguration.NeutralZoneSize,
					tuning.RemoteFootholdCount, configuration.GenerateRoads, tuning, false))
		}
	}
	return zones
}

func (this *HubTopologyService) createConnections(
	playerLabels, outerLabels []string,
	tuning models.GenerationTuning,
	isIsolated bool,
	neutralZones neutral_zone.Plans) []entities.Connection {
	var connections []entities.Connection
	for index, label := range outerLabels {
		hubAnchor := label
		if len(playerLabels) > 0 {
			hubAnchor = playerLabels[0]
		}
		hubGuard := this.GetBorderGuardValue(hubAnchor, label, playerLabels, neutralZones, tuning)
		outerZone := this.ZoneLabelProvider.CreateZoneName(label, playerLabels)
		connections = append(connections,
			variant_content.NewConnectionBuilder().
				WithName(constants.HubZonePrefix+label).
				WithFrom(constants.HubZoneName).
				WithTo(outerZone).
				WithConnectionTypeDirect().
				WithGuardZone(constants.HubZoneName).
				WithSimTurnSquad().
				WithGuardValue(hubGuard).
				WithGuardWeeklyIncrement(0.15).
				WithGuardMatchGroup("hub_guard_"+label).
				Build(),
			variant_content.NewConnectionBuilder().
				WithFrom(constants.HubZoneName).
				WithTo(outerZone).
				WithConnectionTypeDirect().
				WithGuardZone(constants.HubZoneName).
				WithSimTurnSquad().
				WithGuardValue(hubGuard).
				WithGuardWeeklyIncrement(0.15).
				WithGuardMatchGroup(fmt.Sprintf("hub_guard_%s_%d", label, 1)).
				Build())

		labelTo := outerLabels[(index+1)%len(outerLabels)]
		if isIsolated && slices.Contains(playerLabels, label) && slices.Contains(playerLabels, labelTo) {
			continue
		}
		connections = append(connections, variant_content.NewConnectionBuilder().
			WithName(fmt.Sprintf("Pseudo-%s-%s", label, labelTo)).
			WithFrom(this.ZoneLabelProvider.CreateZoneName(label, playerLabels)).
			WithTo(this.ZoneLabelProvider.CreateZoneName(labelTo, playerLabels)).
			WithConnectionTypeProximity().
			Build())
	}
	return connections
}
