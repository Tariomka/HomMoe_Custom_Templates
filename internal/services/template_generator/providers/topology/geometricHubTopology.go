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
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/placement_rule"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
	zone_services "github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
)

// GeometricHubTopologyService builds the "Geometric Hub" topology: every
// player forms a hexagon that shares a central Hub zone, adjacent hexagons
// share corner zones, and additional zones fill the hexagon interiors. All
// connections that touch the Hub are portals. See plans/geometric-hub-topology.md
// for the full slot-growth specification derived from the output/ PoCs.
type GeometricHubTopologyService struct {
	base.TopologyBase
}

func NewGeometricHubTopologyService(
	zoneFactory *zone_services.ZoneFactory,
	roadFactory *zone_services.RoadFactory,
) *GeometricHubTopologyService {
	return &GeometricHubTopologyService{
		TopologyBase: base.NewTopologyBase(zoneFactory, roadFactory),
	}
}

func (this *GeometricHubTopologyService) CreateTopologyVariant(
	configuration config.GeneratorConfig,
	playerLabels []string,
	neutralZones neutral_zone.Plans,
	tuning models.GenerationTuning,
	_ string) entities.Variant {
	layout := newGeometricHubLayout(playerLabels, neutralZones)
	allLabels := append(append([]string{}, playerLabels...),
		linq.FromSlice(neutralZones).
			SelectString(func(plan neutral_zone.Plan) string { return plan.Label }).
			ToSlice()...)

	connectionNames := this.createConnectionNameIndex(layout)
	zones := this.createZones(configuration, playerLabels, neutralZones, layout, connectionNames, tuning)
	conns := this.createConnections(playerLabels, neutralZones, layout, tuning)
	if configuration.RandomPortals {
		conns = append(conns,
			this.CreateRandomPortalConnections(
				playerLabels, allLabels, tuning, configuration.MaxPortalConnections)...)
	}
	return this.CreateVariant(playerLabels, playerLabels[0], len(allLabels)+1, zones, conns)
}

// createConnectionNameIndex maps every label (and "Hub") to the names of its
// named connections, so each zone's roads can reference them.
func (this *GeometricHubTopologyService) createConnectionNameIndex(
	layout *geometricHubLayout) map[string][]string {
	names := map[string][]string{}
	for _, edge := range layout.directEdges {
		name := geometricHubEdgeName(edge)
		names[edge[0]] = append(names[edge[0]], name)
		names[edge[1]] = append(names[edge[1]], name)
	}
	for _, label := range layout.hubPortalLabels {
		name := "Portal-Hub-" + label
		names[label] = append(names[label], name)
		names[constants.HubZoneName] = append(names[constants.HubZoneName], name)
	}
	return names
}

func geometricHubEdgeName(edge [2]string) string {
	return fmt.Sprintf("GeoHub-%s-%s", edge[0], edge[1])
}

func (this *GeometricHubTopologyService) createZones(
	configuration config.GeneratorConfig,
	playerLabels []string,
	neutralZones neutral_zone.Plans,
	layout *geometricHubLayout,
	connectionNames map[string][]string,
	tuning models.GenerationTuning) []entities.Zone {
	hubContentName := ""
	if len(configuration.HubZoneMandatoryContent) > 0 {
		hubContentName = "mandatory_content_hub"
	}
	zones := []entities.Zone{
		this.CreateHubZone(
			constants.HubZoneName, connectionNames[constants.HubZoneName], tuning, configuration.IsHubCityToHold(),
			configuration.ZoneConfiguration.HubZoneSize, configuration.ZoneConfiguration.Advanced.HubZoneCastles,
			configuration.GenerateRoads, hubContentName),
	}
	zones[0].GeneratorPosition = &[2]float64{layoutCenter, layoutCenter}

	for index, label := range playerLabels {
		zone := this.CreateSpawnZone(
			label, fmt.Sprintf("Player%d", index+1), connectionNames[label],
			configuration.ZoneConfiguration.PlayerZoneCastles, configuration.MatchPlayerCastleFactions,
			configuration.ZoneConfiguration.PlayerZoneSize, tuning.RemoteFootholdCount,
			configuration.GenerateRoads, tuning)
		position := layout.positions[label]
		zone.GeneratorPosition = &[2]float64{position.X, position.Y}
		zones = append(zones, zone)
	}
	for _, plan := range neutralZones {
		zone := this.CreateNeutralZone(
			plan, connectionNames[plan.Label], configuration.ZoneConfiguration.NeutralZoneSize,
			tuning.RemoteFootholdCount, configuration.GenerateRoads, tuning, false)
		position := layout.positions[plan.Label]
		zone.GeneratorPosition = &[2]float64{position.X, position.Y}
		zones = append(zones, zone)
	}
	return zones
}

func (this *GeometricHubTopologyService) createConnections(
	playerLabels []string,
	neutralZones neutral_zone.Plans,
	layout *geometricHubLayout,
	tuning models.GenerationTuning) []entities.Connection {
	var connections []entities.Connection
	for _, edge := range layout.directEdges {
		zoneFrom := this.ZoneLabelProvider.CreateZoneName(edge[0], playerLabels)
		connections = append(connections, variant_content.NewConnectionBuilder().
			WithName(geometricHubEdgeName(edge)).
			WithFrom(zoneFrom).
			WithTo(this.ZoneLabelProvider.CreateZoneName(edge[1], playerLabels)).
			WithConnectionTypeDirect().
			WithGuardZone(zoneFrom).
			WithSimTurnSquad().
			WithGuardValue(this.GetBorderGuardValue(edge[0], edge[1], playerLabels, neutralZones, tuning)).
			WithGuardWeeklyIncrement(0.15).
			WithGuardMatchGroup(fmt.Sprintf("geohub_guard_%s_%s", edge[0], edge[1])).
			Build())
	}

	// Rule 11: every connection that touches the Hub is a portal.
	portalRule := placement_rule.NewPlacementRuleBuilder().BuildNearCrossroadsRule(2)
	guardAnchor := playerLabels[0]
	for _, label := range layout.hubPortalLabels {
		guardLabel := label
		if slices.Contains(playerLabels, label) {
			guardLabel = guardAnchor
		}
		connections = append(connections, variant_content.NewConnectionBuilder().
			WithName("Portal-Hub-"+label).
			WithFrom(constants.HubZoneName).
			WithTo(this.ZoneLabelProvider.CreateZoneName(label, playerLabels)).
			WithConnectionTypePortal().
			WithPortalPlacementRulesFrom(portalRule).
			WithPortalPlacementRulesTo(portalRule).
			WithRoad(true).
			WithGuardValue(this.GetBorderGuardValue(guardAnchor, guardLabel, playerLabels, neutralZones, tuning)).
			WithGuardWeeklyIncrement(0.15).
			Build())
	}
	return connections
}
