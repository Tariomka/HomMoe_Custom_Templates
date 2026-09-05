package template_content_model

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template"
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

func ToContentCountLimitModel(entity template.ContentCountLimit) ContentCountLimit {
	return ContentCountLimit{
		Name:   entity.Name,
		Limits: ToContentLimitModels(entity.Limits),
	}
}

func ToContentCountLimitEntity(model ContentCountLimit) template.ContentCountLimit {
	return template.ContentCountLimit{
		Name:   model.Name,
		Limits: ToContentLimitEntities(model.Limits),
	}
}

func ToContentCountLimitModels(entities []template.ContentCountLimit) []ContentCountLimit {
	return helpers.MapSlice(entities, ToContentCountLimitModel)
}

func ToContentCountLimitEntities(models []ContentCountLimit) []template.ContentCountLimit {
	return helpers.MapSlice(models, ToContentCountLimitEntity)
}
