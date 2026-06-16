package constants

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
)

type TopologyDescriptor struct {
	Type        config.MapTopology
	Label       string
	Description string
}

var Topology = struct {
	Circles     TopologyDescriptor
	Random      TopologyDescriptor
	Default     TopologyDescriptor
	HubAndSpoke TopologyDescriptor
	Chain       TopologyDescriptor
	SharedWeb   TopologyDescriptor
	Square      TopologyDescriptor
	Geometric   TopologyDescriptor
	Cross       TopologyDescriptor
	Fractal     TopologyDescriptor
}{
	Circles:     TopologyDescriptor{Type: config.TopologyCircles, Label: "Circles", Description: "Circles: layered concentric rings sorted by zone tier."},
	Random:      TopologyDescriptor{Type: config.TopologyRandom, Label: "Random", Description: "Random: layout decided by the generator."},
	Default:     TopologyDescriptor{Type: config.TopologyDefault, Label: "Ring", Description: "Ring: each player borders two neighbors in a closed loop."},
	HubAndSpoke: TopologyDescriptor{Type: config.TopologyHubAndSpoke, Label: "Hub", Description: "Hub: central neutral hub connects all player zones."},
	Chain:       TopologyDescriptor{Type: config.TopologyChain, Label: "Chain", Description: "Chain: linear series, harder for outer players to interact."},
	SharedWeb:   TopologyDescriptor{Type: config.TopologySharedWeb, Label: "Shared Web", Description: "Shared web: heavy interconnection through central neutral mesh."},
	Square:      TopologyDescriptor{Type: config.TopologySquare, Label: "Square", Description: "Square: players line the edges of a square loop with neutral zones on the edges and inside."},
	Geometric:   TopologyDescriptor{Type: config.TopologyGeometric, Label: "Geometric", Description: "Geometric: zones and connections form symmetric geometric shapes around a centre."},
	Cross:       TopologyDescriptor{Type: config.TopologyCross, Label: "Cross", Description: "Cross: zones and connections radiate from a central hub into cross-shaped arms."},
	Fractal:     TopologyDescriptor{Type: config.TopologyFractal, Label: "Fractal", Description: "Fractal: each player is the base of a fractal that branches inward through low, then high neutral tiers, weaving into a shared centre."},
}

var Topologies = []TopologyDescriptor{
	Topology.Circles,
	Topology.Random,
	Topology.Default,
	Topology.HubAndSpoke,
	Topology.Chain,
	Topology.SharedWeb,
	Topology.Square,
	Topology.Geometric,
	Topology.Cross,
	Topology.Fractal,
}

func GetTopologyDescriptor(topology config.MapTopology) TopologyDescriptor {
	for _, value := range Topologies {
		if value.Type == topology {
			return value
		}
	}
	return Topologies[0]
}
