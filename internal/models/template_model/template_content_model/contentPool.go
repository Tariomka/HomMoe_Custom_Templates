package template_content_model

import (
	"maps"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
)

type ContentPool map[string]any

func (this ContentPool) Clone() ContentPool {
	return maps.Clone(this)
}

func ToContentPoolModel(entity template.ContentPool) ContentPool {
	return ContentPool(entity)
}

func ToContentPoolEntity(model ContentPool) template.ContentPool {
	return template.ContentPool(model)
}

func ToContentPoolModels(entities []template.ContentPool) []ContentPool {
	return helpers.MapSlice(entities, ToContentPoolModel)
}

func ToContentPoolEntities(models []ContentPool) []template.ContentPool {
	return helpers.MapSlice(models, ToContentPoolEntity)
}
