package template_variant_model

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template"
)

type Border struct {
	CornerRadius   float64
	ObstaclesWidth int
	ObstaclesNoise []Noise
	WaterWidth     int
	WaterNoise     []Noise
	WaterType      string
}

func ToBorderModel(entity template.Border) Border {
	return Border{
		CornerRadius:   entity.CornerRadius,
		ObstaclesWidth: entity.ObstaclesWidth,
		ObstaclesNoise: ToNoiseModels(entity.ObstaclesNoise),
		WaterWidth:     entity.WaterWidth,
		WaterNoise:     ToNoiseModels(entity.WaterNoise),
		WaterType:      entity.WaterType,
	}
}

func ToBorderEntity(model Border) template.Border {
	return template.Border{
		CornerRadius:   model.CornerRadius,
		ObstaclesWidth: model.ObstaclesWidth,
		ObstaclesNoise: ToNoiseEntities(model.ObstaclesNoise),
		WaterWidth:     model.WaterWidth,
		WaterNoise:     ToNoiseEntities(model.WaterNoise),
		WaterType:      model.WaterType,
	}
}
