package template_content_model

import (
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
)

type ContentLimit struct {
	SID          string
	IncludeLists []string
	Content      []WeightedContent
	Variant      *int
	MaxCount     int
}

func (this ContentLimit) Clone() ContentLimit {
	clone := this
	clone.IncludeLists = slices.Clone(this.IncludeLists)
	clone.Content = slices.Clone(this.Content)
	clone.Variant = helpers.ClonePointer(this.Variant)
	return clone
}

func ToContentLimitModel(entity template.ContentLimit) ContentLimit {
	return ContentLimit{
		SID:          entity.SID,
		IncludeLists: entity.IncludeLists,
		Content:      ToWeightedContentModels(entity.Content),
		Variant:      entity.Variant,
		MaxCount:     entity.MaxCount,
	}
}

func ToContentLimitEntity(model ContentLimit) template.ContentLimit {
	return template.ContentLimit{
		SID:          model.SID,
		IncludeLists: model.IncludeLists,
		Content:      ToWeightedContentEntities(model.Content),
		Variant:      model.Variant,
		MaxCount:     model.MaxCount,
	}
}

func ToContentLimitModels(entities []template.ContentLimit) []ContentLimit {
	return helpers.MapSlice(entities, ToContentLimitModel)
}

func ToContentLimitEntities(models []ContentLimit) []template.ContentLimit {
	return helpers.MapSlice(models, ToContentLimitEntity)
}
