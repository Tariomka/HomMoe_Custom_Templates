package test_helpers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/provider_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/tournament_variant"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones/zone_interfaces"
)

// NewTopologyServiceLookup mirrors the lookup that internal/composition builds
// for production, so tests exercise the same set of topology services.
func NewTopologyServiceLookup(
	zoneFactory zone_interfaces.IZoneFactory,
	roadFactory zone_interfaces.IRoadFactory,
	zoneLabelProvider zone_interfaces.IZoneLabelProvider,
	connectionService base.ITopologyConnectionService,
) provider_interfaces.ITopologyServiceLookup {
	return providers.NewTopologyServiceLookup(
		topology.NewTournamentTopologyService(
			zoneFactory,
			roadFactory,
			zoneLabelProvider,
			connectionService,
			tournament_variant.NewHubClusterService(zoneFactory, roadFactory, zoneLabelProvider, connectionService),
			tournament_variant.NewBalancedClusterService(
				zoneFactory, roadFactory, zoneLabelProvider, connectionService),
			tournament_variant.NewRingClusterService(zoneFactory, roadFactory, zoneLabelProvider, connectionService),
			tournament_variant.NewChainClusterService(zoneFactory, roadFactory, zoneLabelProvider, connectionService)),
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
		topology.NewFractalTopologyService(zoneFactory, roadFactory, zoneLabelProvider, connectionService))
}
