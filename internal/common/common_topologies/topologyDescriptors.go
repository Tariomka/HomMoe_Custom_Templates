package common_topologies

import (
	"iter"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
)

var descriptorValues = models.TopologyDescriptors{ //nolint:gochecknoglobals // Immutable topology catalog.
	Circles: models.TopologyDescriptor{
		Type:        config.TopologyCircles,
		Label:       "Circles",
		Description: "Circles: layered concentric rings sorted by zone tier.",
		Capabilities: models.TopologyCapabilities{
			LayoutKind:            models.TopologyLayoutScatter,
			UsesGeneratorPosition: true,
			UsesGeneratorRing:     true,
		},
	},
	Random: models.TopologyDescriptor{
		Type:        config.TopologyRandom,
		Label:       "Random",
		Description: "Random: layout decided by the generator.",
		Capabilities: models.TopologyCapabilities{
			LayoutKind:            models.TopologyLayoutScatter,
			UsesGeneratorPosition: true,
		},
	},
	GeometricHub: models.TopologyDescriptor{
		Type:        config.TopologyGeometricHub,
		Label:       "Geometric Hub",
		Description: "Geometric Hub: each player forms a hexagon around a shared central hub; extra zones fill the hexagons.",
		Capabilities: models.TopologyCapabilities{
			LayoutKind:            models.TopologyLayoutFixedGeometry,
			UsesHub:               true,
			UsesGeneratorPosition: true,
		},
	},
	Square: models.TopologyDescriptor{
		Type:        config.TopologySquare,
		Label:       "Square",
		Description: "Square: players line the edges of a square loop with neutral zones on the edges and inside.",
		Capabilities: models.TopologyCapabilities{
			LayoutKind:            models.TopologyLayoutFixedGeometry,
			UsesGeneratorPosition: true,
		},
	},
	Geometric: models.TopologyDescriptor{
		Type:        config.TopologyGeometric,
		Label:       "Geometric",
		Description: "Geometric: zones and connections form symmetric geometric shapes around a center.",
		Capabilities: models.TopologyCapabilities{
			LayoutKind:            models.TopologyLayoutFixedGeometry,
			UsesGeneratorPosition: true,
		},
	},
	Cross: models.TopologyDescriptor{
		Type:        config.TopologyCross,
		Label:       "Cross",
		Description: "Cross: zones and connections radiate from a central hub into cross-shaped arms.",
		Capabilities: models.TopologyCapabilities{
			LayoutKind:            models.TopologyLayoutFixedGeometry,
			UsesGeneratorPosition: true,
		},
	},
	Fractal: models.TopologyDescriptor{
		Type:        config.TopologyFractal,
		Label:       "Fractal",
		Description: "Fractal: each player is the base of a fractal that branches inward through low, then high neutral tiers, weaving into a shared center.",
		Capabilities: models.TopologyCapabilities{
			LayoutKind:            models.TopologyLayoutFixedGeometry,
			UsesGeneratorPosition: true,
		},
	},
}

var topologies = []models.TopologyDescriptor{ //nolint:gochecknoglobals // Immutable display order.
	descriptorValues.Random,
	descriptorValues.Circles,
	descriptorValues.GeometricHub,
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

	return descriptorValues.Random
}

func GetTopologyCapabilities(topology config.MapTopology) models.TopologyCapabilities {
	return GetTopologyDescriptorFromType(topology).Capabilities
}

func GetTopologyDescriptorFromIndex(index int) models.TopologyDescriptor {
	if index >= 0 && index < len(topologies) {
		return topologies[index]
	}

	return topologies[0]
}
