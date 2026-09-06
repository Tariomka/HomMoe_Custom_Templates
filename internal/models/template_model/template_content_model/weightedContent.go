package template_content_model

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
)

type WeightedContent struct {
	template_entity.WeightedContent
}

func ToWeightedContentModel(entity template_entity.WeightedContent) WeightedContent {
	return WeightedContent{WeightedContent: entity}
}

func ToWeightedContentEntity(model WeightedContent) template_entity.WeightedContent {
	return model.WeightedContent
}

func ToWeightedContentModels(entities []template_entity.WeightedContent) []WeightedContent {
	return helpers.MapSlice(entities, ToWeightedContentModel)
}

func ToWeightedContentEntities(models []WeightedContent) []template_entity.WeightedContent {
	return helpers.MapSlice(models, ToWeightedContentEntity)
}
