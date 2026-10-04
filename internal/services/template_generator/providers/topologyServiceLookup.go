package providers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/provider_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology"
)

// TopologyServiceLookup resolves a map topology to the service that builds it.
type TopologyServiceLookup struct {
	tournament provider_interfaces.TopologyVariantCreator
	byTopology map[config.MapTopology]provider_interfaces.TopologyVariantCreator
}

func NewTopologyServiceLookup(
	tournament *topology.TournamentTopologyService,
	geometricHub *topology.GeometricHubTopologyService,
	random *topology.RandomTopologyService,
	circles *topology.CirclesTopologyService,
	square *topology.SquareTopologyService,
	geometric *topology.GeometricTopologyService,
	cross *topology.CrossTopologyService,
	fractal *topology.FractalTopologyService) provider_interfaces.ITopologyServiceLookup {
	return &TopologyServiceLookup{
		tournament: tournament.CreateTopologyVariant,
		byTopology: map[config.MapTopology]provider_interfaces.TopologyVariantCreator{
			config.TopologyGeometricHub: geometricHub.CreateTopologyVariant,
			config.TopologyRandom:       random.CreateTopologyVariant,
			config.TopologyCircles:      circles.CreateTopologyVariant,
			config.TopologySquare:       square.CreateTopologyVariant,
			config.TopologyGeometric:    geometric.CreateTopologyVariant,
			config.TopologyCross:        cross.CreateTopologyVariant,
			config.TopologyFractal:      fractal.CreateTopologyVariant,
		},
	}
}

// Tournament returns the creator used for the two-player tournament variant,
// which is selected by generation mode rather than by map topology.
func (this *TopologyServiceLookup) Tournament() provider_interfaces.TopologyVariantCreator {
	return this.tournament
}

// Resolve returns a registered topology variant creator or reports false for a topology without a service.
func (this *TopologyServiceLookup) Resolve(
	mapTopology config.MapTopology) (provider_interfaces.TopologyVariantCreator, bool) {
	creator, found := this.byTopology[mapTopology]
	return creator, found
}
