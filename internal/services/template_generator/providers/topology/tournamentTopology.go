package topology

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/tournament_variant"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones/zone_interfaces"
)

type TournamentTopologyService struct {
	base.TopologyBase

	clusterService tournament_variant.IClusterService
}

func NewTournamentTopologyService(
	zoneFactory zone_interfaces.IZoneFactory,
	roadFactory zone_interfaces.IRoadFactory,
	zoneLabelProvider zone_interfaces.IZoneLabelProvider,
	connectionService base.ITopologyConnectionService,
	clusterService tournament_variant.IClusterService) *TournamentTopologyService {
	return &TournamentTopologyService{
		TopologyBase:   base.NewTopologyBase(zoneFactory, roadFactory, zoneLabelProvider, connectionService),
		clusterService: clusterService,
	}
}

func (this *TournamentTopologyService) CreateTopologyVariant(
	configuration config.GeneratorConfig,
	playerLabels []string,
	neutralZones neutral_zone.Plans,
	tuning models.GenerationTuning,
	_ string) template_model.Variant {
	perPlayerNeutralZones := this.createPerPlayerNeutralZonePlans(neutralZones)

	var zones []template_model.Zone
	var conns []template_model.Connection
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
				Select(func(x neutral_zone.Plan) string { return x.Label }).
				ToSlice()
			conns = append(conns,
				this.CreateRandomPortalConnections(
					playerLabels, clusterLabels, tuning, configuration.MaxPortalConnections, neutralZones)...)
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
