package composition

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/tournament_variant"
	zone_services "github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
)

// provideTopologyServices builds every topology service exactly once. All of
// them are stateless, so the lookup they are registered in is shared and the
// auto-regeneration loop resolves instead of allocating.
func provideTopologyServices(
	zoneFactory *zone_services.ZoneFactory,
	roadFactory *zone_services.RoadFactory,
	zoneLabelProvider zone_services.IZoneLabelProvider,
	connectionService *base.TopologyConnectionService,
) *providers.TopologyServiceLookup {
	return providers.NewTopologyServiceLookup(
		topology.NewTournamentTopologyService(
			zoneFactory,
			roadFactory,
			zoneLabelProvider,
			connectionService,
			tournament_variant.NewHubClusterService(zoneFactory, roadFactory, zoneLabelProvider, connectionService),
			tournament_variant.NewBalancedClusterService(
				zoneFactory,
				roadFactory,
				zoneLabelProvider,
				connectionService,
			),
			tournament_variant.NewRingClusterService(zoneFactory, roadFactory, zoneLabelProvider, connectionService),
			tournament_variant.NewChainClusterService(zoneFactory, roadFactory, zoneLabelProvider, connectionService),
		),
		topology.NewRingTopologyService(zoneFactory, roadFactory, zoneLabelProvider, connectionService),
		topology.NewHubTopologyService(zoneFactory, roadFactory, zoneLabelProvider, connectionService),
		topology.NewGeometricHubTopologyService(zoneFactory, roadFactory, zoneLabelProvider, connectionService),
		topology.NewChainTopologyService(zoneFactory, roadFactory, zoneLabelProvider, connectionService),
		topology.NewSharedWebTopologyService(zoneFactory, roadFactory, zoneLabelProvider, connectionService),
		topology.NewRandomTopologyService(zoneFactory, roadFactory, zoneLabelProvider, connectionService),
		topology.NewCirclesTopologyService(zoneFactory, roadFactory, zoneLabelProvider, connectionService),
		topology.NewSquareTopologyService(zoneFactory, roadFactory, zoneLabelProvider, connectionService),
		topology.NewGeometricTopologyService(zoneFactory, roadFactory, zoneLabelProvider, connectionService),
		topology.NewCrossTopologyService(zoneFactory, roadFactory, zoneLabelProvider, connectionService),
		topology.NewFractalTopologyService(zoneFactory, roadFactory, zoneLabelProvider, connectionService),
	)
}
