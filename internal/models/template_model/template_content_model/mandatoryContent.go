package template_content_model

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
)

type MandatoryContent struct {
	Name    string
	Content []MandatoryContentItem
}

func ToMandatoryContentModel(entity template.MandatoryContent) MandatoryContent {
	return MandatoryContent{
		Name:    entity.Name,
		Content: ToMandatoryContentItemModels(entity.Content),
	}
}

func ToMandatoryContentEntity(model MandatoryContent) template.MandatoryContent {
	return template.MandatoryContent{
		Name:    model.Name,
		Content: ToMandatoryContentItemEntities(model.Content),
	}
}

func ToMandatoryContentModels(entities []template.MandatoryContent) []MandatoryContent {
	return helpers.MapSlice(entities, ToMandatoryContentModel)
}

func ToMandatoryContentEntities(models []MandatoryContent) []template.MandatoryContent {
	return helpers.MapSlice(models, ToMandatoryContentEntity)
}
