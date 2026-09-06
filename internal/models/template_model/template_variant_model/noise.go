package template_variant_model

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
)

type Noise struct{ template_entity.Noise }

func ToNoiseModel(entity template_entity.Noise) Noise {
	return Noise{Noise: entity}
}

func ToNoiseEntity(model Noise) template_entity.Noise {
	return model.Noise
}

func ToNoiseModels(entities []template_entity.Noise) []Noise {
	return helpers.MapSlice(entities, ToNoiseModel)
}

func ToNoiseEntities(models []Noise) []template_entity.Noise {
	return helpers.MapSlice(models, ToNoiseEntity)
}
