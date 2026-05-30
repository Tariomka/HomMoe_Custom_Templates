package zone_layout

// AmbientPickupDistribution configures distribution of ambient pickups in a zone.
type AmbientPickupDistribution struct {
	Repulsion          float64 `json:"repulsion"`
	Noise              float64 `json:"noise"`
	RoadAttraction     float64 `json:"roadAttraction"`
	ObstacleAttraction float64 `json:"obstacleAttraction"`
	GroupSizeWeights   []int   `json:"groupSizeWeights"`
}
