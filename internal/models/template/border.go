package template

type Border struct {
	CornerRadius   float32 `json:"cornerRadius"`
	ObstaclesWidth int     `json:"obstaclesWidth"`
	ObstaclesNoise []Noise `json:"obstaclesNoise"`
	WaterWidth     int     `json:"waterWidth"`
	WaterNoise     []Noise `json:"waterNoise"`
	WaterType      string  `json:"waterType"`
}
