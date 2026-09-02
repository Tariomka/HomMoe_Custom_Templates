package template_variant_model

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
)

type Noise struct{ template.Noise }

func ToNoiseModel(entity template.Noise) Noise {
	return Noise{Noise: entity}
}

func ToNoiseEntity(model Noise) template.Noise {
	return model.Noise
}

func ToNoiseModels(entities []template.Noise) []Noise {
	return helpers.MapSlice(entities, ToNoiseModel)
}

func ToNoiseEntities(models []Noise) []template.Noise {
	return helpers.MapSlice(models, ToNoiseEntity)
}
