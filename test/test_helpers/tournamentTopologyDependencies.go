package test_helpers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/tournament_variant"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
)

// NewTournamentTopologyDependencies builds the collaborators
// NewTournamentTopologyService takes, in the order it declares them, so it can
// be spread directly into the call:
//
//	topology.NewTournamentTopologyService(test_helpers.NewTournamentTopologyDependencies())
func NewTournamentTopologyDependencies() (
	*zones.ZoneFactory,
	*zones.RoadFactory,
	zones.IZoneLabelProvider,
	*base.TopologyConnectionService,
	tournament_variant.IClusterService,
	tournament_variant.IClusterService,
	tournament_variant.IClusterService,
	tournament_variant.IClusterService) {
	zoneFactory, roadFactory, zoneLabelProvider, connectionService := NewZoneFactories()
	return zoneFactory,
		roadFactory,
		zoneLabelProvider,
		connectionService,
		tournament_variant.NewHubClusterService(zoneFactory, roadFactory, zoneLabelProvider, connectionService),
		tournament_variant.NewBalancedClusterService(zoneFactory, roadFactory, zoneLabelProvider, connectionService),
		tournament_variant.NewRingClusterService(zoneFactory, roadFactory, zoneLabelProvider, connectionService),
		tournament_variant.NewChainClusterService(zoneFactory, roadFactory, zoneLabelProvider, connectionService)
}
