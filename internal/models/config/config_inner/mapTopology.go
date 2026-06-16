package config_inner

// MapTopology enumerates the supported map shapes.
type MapTopology string

const (
	TopologyDefault     MapTopology = "Default" // Ring
	TopologyHubAndSpoke MapTopology = "HubAndSpoke"
	TopologyChain       MapTopology = "Chain"
	TopologySharedWeb   MapTopology = "SharedWeb"
	TopologyRandom      MapTopology = "Random"
	TopologyCircles     MapTopology = "Circles"
	TopologySquare      MapTopology = "Square"
	TopologyGeometric   MapTopology = "Geometric"
	TopologyCross       MapTopology = "Cross"
	TopologyFractal     MapTopology = "Fractal"
)
