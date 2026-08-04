package test_helpers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
)

// NewTopologyServiceLookup mirrors the lookup that internal/composition builds
// for production, so tests exercise the same set of topology services.
func NewTopologyServiceLookup(
	zoneFactory *zones.ZoneFactory,
	roadFactory *zones.RoadFactory,
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
