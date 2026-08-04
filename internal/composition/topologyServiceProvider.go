package composition

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology"
	zone_services "github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
)

// provideTopologyServices builds every topology service exactly once. All of
// them are stateless, so the lookup they are registered in is shared and the
// auto-regeneration loop resolves instead of allocating.
func provideTopologyServices(
	zoneFactory *zone_services.ZoneFactory,
	roadFactory *zone_services.RoadFactory,
) *providers.TopologyServiceLookup {
	return providers.NewTopologyServiceLookup(
		topology.NewTournamentTopologyService(zoneFactory, roadFactory),
		topology.NewRingTopologyService(zoneFactory, roadFactory),
		topology.NewHubTopologyService(zoneFactory, roadFactory),
		topology.NewGeometricHubTopologyService(zoneFactory, roadFactory),
		topology.NewChainTopologyService(zoneFactory, roadFactory),
		topology.NewSharedWebTopologyService(zoneFactory, roadFactory),
		topology.NewRandomTopologyService(zoneFactory, roadFactory),
		topology.NewCirclesTopologyService(zoneFactory, roadFactory),
		topology.NewSquareTopologyService(zoneFactory, roadFactory),
		topology.NewGeometricTopologyService(zoneFactory, roadFactory),
		topology.NewCrossTopologyService(zoneFactory, roadFactory),
		topology.NewFractalTopologyService(zoneFactory, roadFactory),
	)
}
