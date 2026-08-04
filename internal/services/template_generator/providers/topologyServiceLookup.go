package providers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology"
)

// TopologyServiceLookup resolves a map topology to the service that builds it.
// Every service is supplied once and shared, so resolving allocates nothing on
// the auto-regeneration path.
type TopologyServiceLookup struct {
	tournament TopologyVariantCreator
	ring       TopologyVariantCreator
	byTopology map[config.MapTopology]TopologyVariantCreator
}

func NewTopologyServiceLookup(
	tournament *topology.TournamentTopologyService,
	ring *topology.RingTopologyService,
	hub *topology.HubTopologyService,
	geometricHub *topology.GeometricHubTopologyService,
	chain *topology.ChainTopologyService,
	sharedWeb *topology.SharedWebTopologyService,
	random *topology.RandomTopologyService,
	circles *topology.CirclesTopologyService,
	square *topology.SquareTopologyService,
	geometric *topology.GeometricTopologyService,
	cross *topology.CrossTopologyService,
	fractal *topology.FractalTopologyService,
) *TopologyServiceLookup {
	return &TopologyServiceLookup{
		tournament: tournament.CreateTopologyVariant,
		ring:       ring.CreateTopologyVariant,
		byTopology: map[config.MapTopology]TopologyVariantCreator{
			config.TopologyHubAndSpoke:  hub.CreateTopologyVariant,
			config.TopologyGeometricHub: geometricHub.CreateTopologyVariant,
			config.TopologyChain:        chain.CreateTopologyVariant,
			config.TopologySharedWeb:    sharedWeb.CreateTopologyVariant,
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
func (this *TopologyServiceLookup) Tournament() TopologyVariantCreator {
	return this.tournament
}

// Resolve falls back to the ring topology for any topology without an entry,
// which is what config.TopologyRing ("Default") and unknown values get.
func (this *TopologyServiceLookup) Resolve(mapTopology config.MapTopology) TopologyVariantCreator {
	if creator, found := this.byTopology[mapTopology]; found {
		return creator
	}
	return this.ring
}
