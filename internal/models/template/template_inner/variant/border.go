package variant

// Border describes the outer obstacles/water boundary of a variant.
type Border struct {
	CornerRadius   float64 `json:"cornerRadius"`
	ObstaclesWidth int     `json:"obstaclesWidth"`
	ObstaclesNoise []Noise `json:"obstaclesNoise"`
	WaterWidth     int     `json:"waterWidth"`
	WaterNoise     []Noise `json:"waterNoise"`
	WaterType      string  `json:"waterType"`
}
