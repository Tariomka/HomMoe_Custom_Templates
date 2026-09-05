package template_layout_model

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
)

type ElevationMode struct {
	template.ElevationMode
}

func ToElevationModeModel(entity template.ElevationMode) ElevationMode {
	return ElevationMode{ElevationMode: entity}
}

func ToElevationModeEntity(model ElevationMode) template.ElevationMode {
	return model.ElevationMode
}

func ToElevationModeModels(entities []template.ElevationMode) []ElevationMode {
	return helpers.MapSlice(entities, ToElevationModeModel)
}

func ToElevationModeEntities(models []ElevationMode) []template.ElevationMode {
	return helpers.MapSlice(models, ToElevationModeEntity)
}
