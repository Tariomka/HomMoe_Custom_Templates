package template_layout_model

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
)

type ElevationMode struct {
	template_entity.ElevationMode
}

func ToElevationModeModel(entity template_entity.ElevationMode) ElevationMode {
	return ElevationMode{ElevationMode: entity}
}

func ToElevationModeEntity(model ElevationMode) template_entity.ElevationMode {
	return model.ElevationMode
}

func ToElevationModeModels(entities []template_entity.ElevationMode) []ElevationMode {
	return helpers.MapSlice(entities, ToElevationModeModel)
}

func ToElevationModeEntities(models []ElevationMode) []template_entity.ElevationMode {
	return helpers.MapSlice(models, ToElevationModeEntity)
}
