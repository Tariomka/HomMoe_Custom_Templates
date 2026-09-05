package template_content_model

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
)

type WeightedContent struct{ template.WeightedContent }

func ToWeightedContentModel(entity template.WeightedContent) WeightedContent {
	return WeightedContent{WeightedContent: entity}
}

func ToWeightedContentEntity(model WeightedContent) template.WeightedContent {
	return model.WeightedContent
}

func ToWeightedContentModels(entities []template.WeightedContent) []WeightedContent {
	return helpers.MapSlice(entities, ToWeightedContentModel)
}

func ToWeightedContentEntities(models []WeightedContent) []template.WeightedContent {
	return helpers.MapSlice(models, ToWeightedContentEntity)
}
