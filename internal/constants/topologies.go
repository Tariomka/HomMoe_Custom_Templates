package constants

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/generator"
)

type TopologyDescriptor struct {
	Type        generator.MapTopology
	Label       string
	Description string
}

var Topology = struct {
	Balanced    TopologyDescriptor
	Random      TopologyDescriptor
	Default     TopologyDescriptor
	HubAndSpoke TopologyDescriptor
	Chain       TopologyDescriptor
	SharedWeb   TopologyDescriptor
}{
	Balanced:    TopologyDescriptor{Type: generator.TopologyBalanced, Label: "Balanced", Description: "Balanced: layered rings sorted by zone tier."},
	Random:      TopologyDescriptor{Type: generator.TopologyRandom, Label: "Random", Description: "Random: layout decided by the generator."},
	Default:     TopologyDescriptor{Type: generator.TopologyDefault, Label: "Ring", Description: "Ring: each player borders two neighbors in a closed loop."},
	HubAndSpoke: TopologyDescriptor{Type: generator.TopologyHubAndSpoke, Label: "Hub", Description: "Hub: central neutral hub connects all player zones."},
	Chain:       TopologyDescriptor{Type: generator.TopologyChain, Label: "Chain", Description: "Chain: linear series, harder for outer players to interact."},
	SharedWeb:   TopologyDescriptor{Type: generator.TopologySharedWeb, Label: "Shared Web", Description: "Shared web: heavy interconnection through central neutral mesh."},
}

var Topologies = []TopologyDescriptor{
	Topology.Balanced,
	Topology.Random,
	Topology.Default,
	Topology.HubAndSpoke,
	Topology.Chain,
	Topology.SharedWeb,
}

func GetTopologyDescriptor(topology generator.MapTopology) TopologyDescriptor {
	for _, value := range Topologies {
		if value.Type == topology {
			return value
		}
	}
	return Topologies[0]
}
