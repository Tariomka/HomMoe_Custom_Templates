package constants

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/generator"
)

var Topology = struct {
	Balanced    models.TopologyDescriptor
	Random      models.TopologyDescriptor
	Default     models.TopologyDescriptor
	HubAndSpoke models.TopologyDescriptor
	Chain       models.TopologyDescriptor
	SharedWeb   models.TopologyDescriptor
}{
	Balanced:    models.TopologyDescriptor{Type: generator.TopologyBalanced, Label: "Balanced", Description: "NOT ADDED YET"},
	Random:      models.TopologyDescriptor{Type: generator.TopologyRandom, Label: "Random", Description: "Random: layout decided by the generator."},
	Default:     models.TopologyDescriptor{Type: generator.TopologyDefault, Label: "Ring", Description: "Ring: each player borders two neighbors in a closed loop."},
	HubAndSpoke: models.TopologyDescriptor{Type: generator.TopologyHubAndSpoke, Label: "Hub", Description: "Hub: central neutral hub connects all player zones."},
	Chain:       models.TopologyDescriptor{Type: generator.TopologyChain, Label: "Chain", Description: "Chain: linear series, harder for outer players to interact."},
	SharedWeb:   models.TopologyDescriptor{Type: generator.TopologySharedWeb, Label: "Shared Web", Description: "Shared web: heavy interconnection through central neutral mesh."},
}

var Topologies = []models.TopologyDescriptor{
	Topology.Balanced,
	Topology.Random,
	Topology.Default,
	Topology.HubAndSpoke,
	Topology.Chain,
	Topology.SharedWeb,
}

func GetTopologyDescriptor(topology models.MapTopology) models.TopologyDescriptor {
	for _, value := range Topologies {
		if value.Type == topology {
			return value
		}
	}
	return Topologies[0]
}
