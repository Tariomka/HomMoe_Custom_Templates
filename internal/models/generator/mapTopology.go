package generator

// MapTopology enumerates the supported map shapes.
type MapTopology string

const (
	TopologyDefault     MapTopology = "Default" // Ring
	TopologyHubAndSpoke MapTopology = "HubAndSpoke"
	TopologyChain       MapTopology = "Chain"
	TopologySharedWeb   MapTopology = "SharedWeb"
	TopologyRandom      MapTopology = "Random"
	// TopologyBalanced is the ring layout with evenly distributed castle
	// counts and balanced player zone placement (see C# v0.7).
	TopologyBalanced MapTopology = "Balanced"
)
