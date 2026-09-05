package template_content_model

import (
	"maps"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
)

type ContentList map[string]any

func (this ContentList) Clone() ContentList {
	return maps.Clone(this)
}

func ToContentListModel(entity template.ContentList) ContentList {
	return ContentList(entity)
}

func ToContentListEntity(model ContentList) template.ContentList {
	return template.ContentList(model)
}

func ToContentListModels(entities []template.ContentList) []ContentList {
	return helpers.MapSlice(entities, ToContentListModel)
}

func ToContentListEntities(models []ContentList) []template.ContentList {
	return helpers.MapSlice(models, ToContentListEntity)
}
