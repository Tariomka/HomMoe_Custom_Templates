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

type HubTopologyService struct {
	base.TopologyBase
}

func NewHubTopologyService() *HubTopologyService {
	return &HubTopologyService{
		TopologyBase: base.NewTopologyBase(),
	}
}

func (this *HubTopologyService) CreateTopologyVariant(
	configuration config.GeneratorConfig,
	playerLabels []string,
	neutralZones models.NeutralZonePlans,
	tuning models.GenerationTuning,
	hubIsHoldCity bool) entities.Variant {
	outerLabels := this.createOuterLabels(configuration, playerLabels, neutralZones)

	zones := this.createZones(configuration, playerLabels, outerLabels, tuning, neutralZones, hubIsHoldCity)
	conns := this.createConnections(playerLabels, outerLabels, tuning, configuration.NoDirectPlayerConnections, neutralZones)
	if configuration.RandomPortals {
		conns = append(conns, this.CreateRandomPortalConnections(playerLabels, outerLabels, tuning, configuration.MaxPortalConnections)...)
	}
	return this.CreateVariant(playerLabels, outerLabels[0], len(outerLabels)+1, zones, conns)
}

func (this *HubTopologyService) createOuterLabels(
	configuration config.GeneratorConfig,
	playerLabels []string,
	neutralZones models.NeutralZonePlans) []string {
	if configuration.Topology == config.TopologyBalanced {
		sep := 0
		if configuration.MinNeutralZonesBetweenPlayers > 0 && configuration.CanHonorNeutralSeparation() {
			sep = configuration.MinNeutralZonesBetweenPlayers
		}
		return this.ZoneLabelProvider.CreateBalancedChainZoneLabels(playerLabels, neutralZones, sep)
	}

	return append(playerLabels,
		linq.FromSlice(neutralZones).
			SelectString(func(nz models.NeutralZonePlan) string { return nz.Label }).
			ToSlice()...)
}

func (this *HubTopologyService) createZones(
	configuration config.GeneratorConfig,
	playerLabels, outerLabels []string,
	tuning models.GenerationTuning,
	neutralZones models.NeutralZonePlans,
	hubIsHoldCity bool) []entities.Zone {
	hubConns := make([]string, len(outerLabels))
	for index, label := range outerLabels {
		hubConns[index] = "Hub-" + label
	}
	zones := []entities.Zone{this.CreateHubZone(
		hubConns, tuning, hubIsHoldCity, configuration.ZoneConfiguration.HubZoneSize,
		configuration.ZoneConfiguration.HubZoneCastles, configuration.GenerateRoads)}

	for _, label := range outerLabels {
		spokeConnectionNames := []string{"Hub-" + label}
		if playerIndex := slices.Index(playerLabels, label); playerIndex >= 0 {
			zones = append(zones,
				this.CreateSpawnZone(
					label, fmt.Sprintf("Player%d", playerIndex+1), spokeConnectionNames,
					configuration.ZoneConfiguration.PlayerZoneCastles, configuration.MatchPlayerCastleFactions,
					configuration.ZoneConfiguration.Advanced.PlayerZoneSize, configuration.SpawnRemoteFootholds,
					configuration.GenerateRoads, tuning))
		} else {
			zones = append(zones,
				this.CreateNeutralZone(
					linq.FromSlice(neutralZones).FirstOrDefault(func(x models.NeutralZonePlan) bool { return x.Label == label }),
					spokeConnectionNames, configuration.ZoneConfiguration.Advanced.NeutralZoneSize,
					configuration.SpawnRemoteFootholds, configuration.GenerateRoads, tuning, false))
		}
	}
	return zones
}

func (this *HubTopologyService) createConnections(
	playerLabels, outerLabels []string,
	tuning models.GenerationTuning,
	isIsolated bool,
	neutralZones models.NeutralZonePlans) []entities.Connection {
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
				WithName("Hub-"+label).
				WithFrom("Hub").
				WithTo(outerZone).
				WithConnectionTypeDirect().
				WithGuardZone("Hub").
				WithSimTurnSquad().
				WithGuardValue(hubGuard).
				WithGuardWeeklyIncrement(0.15).
				WithGuardMatchGroup("hub_guard_"+label).
				Build(),
			variant_content.NewConnectionBuilder().
				WithFrom("Hub").
				WithTo(outerZone).
				WithConnectionTypeDirect().
				WithGuardZone("Hub").
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
