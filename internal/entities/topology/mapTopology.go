package topology

// MapTopology enumerates the supported map shapes.
type MapTopology string

const (
	TopologyRandom       MapTopology = "Random"
	TopologyCircles      MapTopology = "Circles"
	TopologySquare       MapTopology = "Square"
	TopologyGeometric    MapTopology = "Geometric"
	TopologyCross        MapTopology = "Cross"
	TopologyFractal      MapTopology = "Fractal"
	TopologyGeometricHub MapTopology = "GeometricHub"
)
