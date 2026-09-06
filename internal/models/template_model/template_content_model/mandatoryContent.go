package template_content_model

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
)

type MandatoryContent struct {
	Name    string
	Content []MandatoryContentItem
}

func (this MandatoryContent) Clone() MandatoryContent {
	clone := this
	clone.Content = helpers.MapSlice(this.Content, MandatoryContentItem.Clone)
	return clone
}

func ToMandatoryContentModel(entity template_entity.MandatoryContent) MandatoryContent {
	return MandatoryContent{
		Name:    entity.Name,
		Content: ToMandatoryContentItemModels(entity.Content),
	}
}

func ToMandatoryContentEntity(model MandatoryContent) template_entity.MandatoryContent {
	return template_entity.MandatoryContent{
		Name:    model.Name,
		Content: ToMandatoryContentItemEntities(model.Content),
	}
}

func ToMandatoryContentModels(entities []template_entity.MandatoryContent) []MandatoryContent {
	return helpers.MapSlice(entities, ToMandatoryContentModel)
}

func ToMandatoryContentEntities(models []MandatoryContent) []template_entity.MandatoryContent {
	return helpers.MapSlice(models, ToMandatoryContentEntity)
}
