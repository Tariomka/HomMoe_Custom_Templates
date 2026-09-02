package template_override_model

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
)

type ValueOverride struct{ template.ValueOverride }

func ToValueOverrideModel(entity template.ValueOverride) ValueOverride {
	return ValueOverride{ValueOverride: entity}
}

func ToValueOverrideEntity(model ValueOverride) template.ValueOverride {
	return model.ValueOverride
}

func ToValueOverrideModels(entities []template.ValueOverride) []ValueOverride {
	return helpers.MapSlice(entities, ToValueOverrideModel)
}

func ToValueOverrideEntities(models []ValueOverride) []template.ValueOverride {
	return helpers.MapSlice(models, ToValueOverrideEntity)
}
