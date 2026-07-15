package common

import (
	"iter"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
)

type TopologyDescriptor struct {
	Type        config.MapTopology
	Label       string
	Description string
}

type TopologyDescriptors struct {
	Default      TopologyDescriptor
	Circles      TopologyDescriptor
	Random       TopologyDescriptor
	HubAndSpoke  TopologyDescriptor
	GeometricHub TopologyDescriptor
	Chain        TopologyDescriptor
	SharedWeb    TopologyDescriptor
	Square       TopologyDescriptor
	Geometric    TopologyDescriptor
	Cross        TopologyDescriptor
	Fractal      TopologyDescriptor
}

var descriptorValues = TopologyDescriptors{
	Default: TopologyDescriptor{
		Type:        config.TopologyRing,
		Label:       "Ring",
		Description: "Ring: each player borders two neighbors in a closed loop.",
	},
	Circles: TopologyDescriptor{
		Type:        config.TopologyCircles,
		Label:       "Circles",
		Description: "Circles: layered concentric rings sorted by zone tier.",
	},
	Random: TopologyDescriptor{
		Type:        config.TopologyRandom,
		Label:       "Random",
		Description: "Random: layout decided by the generator.",
	},
	HubAndSpoke: TopologyDescriptor{
		Type:        config.TopologyHubAndSpoke,
		Label:       "Hub",
		Description: "Hub: central neutral hub connects all player zones.",
	},
	GeometricHub: TopologyDescriptor{
		Type:        config.TopologyGeometricHub,
		Label:       "Geometric Hub",
		Description: "Geometric Hub: each player forms a hexagon around a shared central hub; extra zones fill the hexagons.",
	},
	Chain: TopologyDescriptor{
		Type:        config.TopologyChain,
		Label:       "Chain",
		Description: "Chain: linear series, harder for outer players to interact.",
	},
	SharedWeb: TopologyDescriptor{
		Type:        config.TopologySharedWeb,
		Label:       "Shared Web",
		Description: "Shared web: heavy interconnection through central neutral mesh.",
	},
	Square: TopologyDescriptor{
		Type:        config.TopologySquare,
		Label:       "Square",
		Description: "Square: players line the edges of a square loop with neutral zones on the edges and inside.",
	},
	Geometric: TopologyDescriptor{
		Type:        config.TopologyGeometric,
		Label:       "Geometric",
		Description: "Geometric: zones and connections form symmetric geometric shapes around a center.",
	},
	Cross: TopologyDescriptor{
		Type:        config.TopologyCross,
		Label:       "Cross",
		Description: "Cross: zones and connections radiate from a central hub into cross-shaped arms.",
	},
	Fractal: TopologyDescriptor{
		Type:        config.TopologyFractal,
		Label:       "Fractal",
		Description: "Fractal: each player is the base of a fractal that branches inward through low, then high neutral tiers, weaving into a shared center.",
	},
}

var topologies = []TopologyDescriptor{
	descriptorValues.Random,
	descriptorValues.Default,
	descriptorValues.Circles,
	descriptorValues.HubAndSpoke,
	descriptorValues.GeometricHub,
	descriptorValues.Chain,
	descriptorValues.SharedWeb,
	descriptorValues.Square,
	descriptorValues.Geometric,
	descriptorValues.Cross,
	descriptorValues.Fractal,
}

func GetTopologyDescriptors() TopologyDescriptors {
	return descriptorValues
}

func GetTopologyDescriptorSeq() iter.Seq[TopologyDescriptor] {
	return func(yield func(TopologyDescriptor) bool) {
		for _, value := range topologies {
			if !yield(value) {
				return
			}
		}
	}
}

func GetTopologyDescriptorFromType(topology config.MapTopology) TopologyDescriptor {
	for value := range GetTopologyDescriptorSeq() {
		if value.Type == topology {
			return value
		}
	}

	return topologies[0]
}

func GetTopologyDescriptorFromIndex(index int) TopologyDescriptor {
	if index >= 0 && index < len(topologies) {
		return topologies[index]
	}

	return topologies[0]
}
