package template_content_model

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
)

type ContentCountLimit struct {
	Name   string
	Limits []ContentLimit
}

func (this ContentCountLimit) Clone() ContentCountLimit {
	clone := this
	clone.Limits = helpers.MapSlice(this.Limits, ContentLimit.Clone)
	return clone
}

func ToContentCountLimitModel(entity template_entity.ContentCountLimit) ContentCountLimit {
	return ContentCountLimit{
		Name:   entity.Name,
		Limits: ToContentLimitModels(entity.Limits),
	}
}

func ToContentCountLimitEntity(model ContentCountLimit) template_entity.ContentCountLimit {
	return template_entity.ContentCountLimit{
		Name:   model.Name,
		Limits: ToContentLimitEntities(model.Limits),
	}
}

func ToContentCountLimitModels(entities []template_entity.ContentCountLimit) []ContentCountLimit {
	return helpers.MapSlice(entities, ToContentCountLimitModel)
}

func ToContentCountLimitEntities(models []ContentCountLimit) []template_entity.ContentCountLimit {
	return helpers.MapSlice(models, ToContentCountLimitEntity)
}
