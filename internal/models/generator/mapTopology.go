package generator

// MapTopology enumerates the supported map shapes.
type MapTopology string

const (
	TopologyDefault     MapTopology = "Default" // Ring
	TopologyHubAndSpoke MapTopology = "HubAndSpoke"
	TopologyChain       MapTopology = "Chain"
	TopologySharedWeb   MapTopology = "SharedWeb"
	TopologyRandom      MapTopology = "Random"
)
