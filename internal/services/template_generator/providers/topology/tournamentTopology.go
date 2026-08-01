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

	clusterService   tournament_variant.IClusterService
	creationServices *zone_services.CreationServices
}

func NewTournamentTopologyService() *TournamentTopologyService {
	return NewTournamentTopologyServiceWithCreationServices(zone_services.NewCreationServices(nil, nil))
}

func NewTournamentTopologyServiceWithCreationServices(
	creationServices *zone_services.CreationServices,
) *TournamentTopologyService {
	if creationServices == nil {
		creationServices = zone_services.NewCreationServices(nil, nil)
	}
	return &TournamentTopologyService{
		TopologyBase:     base.NewTopologyBaseWithCreationServices(creationServices),
		creationServices: creationServices,
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
		this.clusterService = tournament_variant.NewHubClusterServiceWithCreationServices(this.creationServices)
	case config.TopologyCircles:
		this.clusterService = tournament_variant.NewBalancedClusterServiceWithCreationServices(this.creationServices)
	case config.TopologyRing:
		this.clusterService = tournament_variant.NewRingClusterServiceWithCreationServices(this.creationServices)
	default:
		// Chain, SharedWeb, Random → chain-per-cluster fallback.
		this.clusterService = tournament_variant.NewChainClusterServiceWithCreationServices(this.creationServices)
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
