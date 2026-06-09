package topology

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/tournament_variant"
)

type TournamentTopologyService struct {
	base.TopologyBase
	clusterService tournament_variant.IClusterService
}

func NewTournamentTopologyService() *TournamentTopologyService {
	return &TournamentTopologyService{
		TopologyBase: base.NewTopologyBase(),
	}
}

func (this *TournamentTopologyService) CreateTopologyVariant(
	configuration config.GeneratorConfig,
	playerLabels []string,
	neutralZones models.NeutralZonePlans,
	tuning models.GenerationTuning) template.Variant {
	perPlayerNeutralZones := this.createPerPlayerNeutralZonePlans(neutralZones)

	switch configuration.Topology {
	case config.TopologyHubAndSpoke:
		this.clusterService = tournament_variant.NewHubClusterService()
	case config.TopologyBalanced:
		this.clusterService = tournament_variant.NewBalancedClusterService()
	case config.TopologyDefault:
		this.clusterService = tournament_variant.NewRingClusterService()
	default:
		// Chain, SharedWeb, Random → chain-per-cluster fallback.
		this.clusterService = tournament_variant.NewChainClusterService()
	}

	var zones []template.Zone
	var conns []template.Connection
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
				SelectString(func(x models.NeutralZonePlan) string { return x.Label }).
				ToSlice()
			conns = append(conns, this.CreateRandomPortalConnections(playerLabels, clusterLabels, tuning, configuration.MaxPortalConnections)...)
		}
	}
	return this.CreateVariant(playerLabels, playerLabels[0], len(zones), zones, conns)
}

// createPerPlayerNeutralZonePlans splits neutral zone plans into two per-player lists
// in a balanced round-robin so that quality tiers are split evenly across the two players.
func (this *TournamentTopologyService) createPerPlayerNeutralZonePlans(neutralZones models.NeutralZonePlans) [2]models.NeutralZonePlans {
	perPlayerNeutralZones := [2]models.NeutralZonePlans{}

	sorted := models.NewNeutralZonePlansSorted(neutralZones)
	for index, zonePlan := range *sorted {
		perPlayerNeutralZones[index%2].AddPlans(zonePlan)
	}
	for index := range perPlayerNeutralZones {
		perPlayerNeutralZones[index].SortByBalanceScoreAscending()
	}

	return perPlayerNeutralZones
}
