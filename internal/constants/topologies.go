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
	Balanced    TopologyDescriptor
	Random      TopologyDescriptor
	Default     TopologyDescriptor
	HubAndSpoke TopologyDescriptor
	Chain       TopologyDescriptor
	SharedWeb   TopologyDescriptor
}{
	Balanced:    TopologyDescriptor{Type: config.TopologyBalanced, Label: "Balanced", Description: "Balanced: layered rings sorted by zone tier."},
	Random:      TopologyDescriptor{Type: config.TopologyRandom, Label: "Random", Description: "Random: layout decided by the generator."},
	Default:     TopologyDescriptor{Type: config.TopologyDefault, Label: "Ring", Description: "Ring: each player borders two neighbors in a closed loop."},
	HubAndSpoke: TopologyDescriptor{Type: config.TopologyHubAndSpoke, Label: "Hub", Description: "Hub: central neutral hub connects all player zones."},
	Chain:       TopologyDescriptor{Type: config.TopologyChain, Label: "Chain", Description: "Chain: linear series, harder for outer players to interact."},
	SharedWeb:   TopologyDescriptor{Type: config.TopologySharedWeb, Label: "Shared Web", Description: "Shared web: heavy interconnection through central neutral mesh."},
}

var Topologies = []TopologyDescriptor{
	Topology.Balanced,
	Topology.Random,
	Topology.Default,
	Topology.HubAndSpoke,
	Topology.Chain,
	Topology.SharedWeb,
}

func GetTopologyDescriptor(topology config.MapTopology) TopologyDescriptor {
	for _, value := range Topologies {
		if value.Type == topology {
			return value
		}
	}
	return Topologies[0]
}
