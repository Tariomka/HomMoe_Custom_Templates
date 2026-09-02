package template_content_model

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
)

type ContentList map[string]any

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
