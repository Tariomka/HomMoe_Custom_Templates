package topology

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/tournament_variant"
	zone_services "github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
)

type TournamentTopologyService struct {
	base.TopologyBase

	clusterService tournament_variant.IClusterService
	zoneFactory    *zone_services.ZoneFactory
	roadFactory    *zone_services.RoadFactory
}

func NewTournamentTopologyService(
	zoneFactory *zone_services.ZoneFactory,
	roadFactory *zone_services.RoadFactory,
) *TournamentTopologyService {
	return &TournamentTopologyService{
		TopologyBase: base.NewTopologyBase(zoneFactory, roadFactory),
		zoneFactory:  zoneFactory,
		roadFactory:  roadFactory,
	}
}

func (this *TournamentTopologyService) CreateTopologyVariant(
	configuration config.GeneratorConfig,
	playerLabels []string,
	neutralZones neutral_zone.Plans,
	tuning models.GenerationTuning) entities.Variant {
	perPlayerNeutralZones := this.createPerPlayerNeutralZonePlans(neutralZones)

	switch configuration.Topology {
	case config.TopologyHubAndSpoke:
		this.clusterService = tournament_variant.NewHubClusterService(this.zoneFactory, this.roadFactory)
	case config.TopologyCircles:
		this.clusterService = tournament_variant.NewBalancedClusterService(this.zoneFactory, this.roadFactory)
	case config.TopologyRing:
		this.clusterService = tournament_variant.NewRingClusterService(this.zoneFactory, this.roadFactory)
	default:
		// Chain, SharedWeb, Random → chain-per-cluster fallback.
		this.clusterService = tournament_variant.NewChainClusterService(this.zoneFactory, this.roadFactory)
	}

	var zones []entities.Zone
	var conns []entities.Connection
	for playerIndex := range 2 {
		perPlayerZones, perPlayerConns := this.clusterService.CreateClusterVariant(
			configuration,
			tuning,
			neutralZones,
			perPlayerNeutralZones[playerIndex],
			playerIndex,
			playerLabels[playerIndex])
		zones = append(zones, perPlayerZones...)
		conns = append(conns, perPlayerConns...)
	}

	if configuration.RandomPortals {
		for playerIndex := range 2 {
			clusterLabels := linq.FromSlice(perPlayerNeutralZones[playerIndex]).
				SelectString(func(x neutral_zone.Plan) string { return x.Label }).
				ToSlice()
			conns = append(conns,
				this.CreateRandomPortalConnections(
					playerLabels, clusterLabels, tuning, configuration.MaxPortalConnections)...)
		}
	}
	return this.CreateVariant(playerLabels, playerLabels[0], len(zones), zones, conns)
}

// createPerPlayerNeutralZonePlans splits neutral zone plans into two per-player lists
// in a balanced round-robin so that quality tiers are split evenly across the two players.
func (this *TournamentTopologyService) createPerPlayerNeutralZonePlans(
	neutralZones neutral_zone.Plans) [2]neutral_zone.Plans {
	perPlayerNeutralZones := [2]neutral_zone.Plans{}

	sorted := neutral_zone.NewNeutralZonePlansSorted(neutralZones)
	for index, zonePlan := range *sorted {
		perPlayerNeutralZones[index%2].AddPlans(zonePlan)
	}
	for index := range perPlayerNeutralZones {
		perPlayerNeutralZones[index].SortByBalanceScoreAscending()
	}

	return perPlayerNeutralZones
}
