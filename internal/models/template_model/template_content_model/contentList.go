package template_content_model

import (
	"maps"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
)

type ContentList map[string]any

func (this ContentList) Clone() ContentList {
	return maps.Clone(this)
}

func ToContentListModel(entity template_entity.ContentList) ContentList {
	return ContentList(entity)
}

func ToContentListEntity(model ContentList) template_entity.ContentList {
	return template_entity.ContentList(model)
}

func ToContentListModels(entities []template_entity.ContentList) []ContentList {
	return helpers.MapSlice(entities, ToContentListModel)
}

func ToContentListEntities(models []ContentList) []template_entity.ContentList {
	return helpers.MapSlice(models, ToContentListEntity)
}
