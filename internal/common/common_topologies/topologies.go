package common_topologies

import (
	"iter"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
)

var descriptorValues = models.TopologyDescriptors{
	Default: models.TopologyDescriptor{
		Type:        config.TopologyRing,
		Label:       "Ring",
		Description: "Ring: each player borders two neighbors in a closed loop.",
	},
	Circles: models.TopologyDescriptor{
		Type:        config.TopologyCircles,
		Label:       "Circles",
		Description: "Circles: layered concentric rings sorted by zone tier.",
	},
	Random: models.TopologyDescriptor{
		Type:        config.TopologyRandom,
		Label:       "Random",
		Description: "Random: layout decided by the generator.",
	},
	HubAndSpoke: models.TopologyDescriptor{
		Type:        config.TopologyHubAndSpoke,
		Label:       "Hub",
		Description: "Hub: central neutral hub connects all player zones.",
	},
	GeometricHub: models.TopologyDescriptor{
		Type:        config.TopologyGeometricHub,
		Label:       "Geometric Hub",
		Description: "Geometric Hub: each player forms a hexagon around a shared central hub; extra zones fill the hexagons.",
	},
	Chain: models.TopologyDescriptor{
		Type:        config.TopologyChain,
		Label:       "Chain",
		Description: "Chain: linear series, harder for outer players to interact.",
	},
	SharedWeb: models.TopologyDescriptor{
		Type:        config.TopologySharedWeb,
		Label:       "Shared Web",
		Description: "Shared web: heavy interconnection through central neutral mesh.",
	},
	Square: models.TopologyDescriptor{
		Type:        config.TopologySquare,
		Label:       "Square",
		Description: "Square: players line the edges of a square loop with neutral zones on the edges and inside.",
	},
	Geometric: models.TopologyDescriptor{
		Type:        config.TopologyGeometric,
		Label:       "Geometric",
		Description: "Geometric: zones and connections form symmetric geometric shapes around a center.",
	},
	Cross: models.TopologyDescriptor{
		Type:        config.TopologyCross,
		Label:       "Cross",
		Description: "Cross: zones and connections radiate from a central hub into cross-shaped arms.",
	},
	Fractal: models.TopologyDescriptor{
		Type:        config.TopologyFractal,
		Label:       "Fractal",
		Description: "Fractal: each player is the base of a fractal that branches inward through low, then high neutral tiers, weaving into a shared center.",
	},
}

var topologies = []models.TopologyDescriptor{
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

func GetTopologyDescriptors() models.TopologyDescriptors {
	return descriptorValues
}

func GetTopologyDescriptorSeq() iter.Seq[models.TopologyDescriptor] {
	return func(yield func(models.TopologyDescriptor) bool) {
		for _, value := range topologies {
			if !yield(value) {
				return
			}
		}
	}
}

func GetTopologyDescriptorFromType(topology config.MapTopology) models.TopologyDescriptor {
	for value := range GetTopologyDescriptorSeq() {
		if value.Type == topology {
			return value
		}
	}

	return topologies[0]
}

func GetTopologyDescriptorFromIndex(index int) models.TopologyDescriptor {
	if index >= 0 && index < len(topologies) {
		return topologies[index]
	}

	return topologies[0]
}
