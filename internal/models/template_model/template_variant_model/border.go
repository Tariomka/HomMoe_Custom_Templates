package template_variant_model

import "slices"

type Border struct {
	CornerRadius   float64
	ObstaclesWidth int
	ObstaclesNoise []Noise
	WaterWidth     int
	WaterNoise     []Noise
	WaterType      string
}

func (this Border) Clone() Border {
	clone := this
	clone.ObstaclesNoise = slices.Clone(this.ObstaclesNoise)
	clone.WaterNoise = slices.Clone(this.WaterNoise)
	return clone
}
