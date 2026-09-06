package template_override_model

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
)

type ValueOverride struct{ template_entity.ValueOverride }

func ToValueOverrideModel(entity template_entity.ValueOverride) ValueOverride {
	return ValueOverride{ValueOverride: entity}
}

func ToValueOverrideEntity(model ValueOverride) template_entity.ValueOverride {
	return model.ValueOverride
}

func ToValueOverrideModels(entities []template_entity.ValueOverride) []ValueOverride {
	return helpers.MapSlice(entities, ToValueOverrideModel)
}

func ToValueOverrideEntities(models []ValueOverride) []template_entity.ValueOverride {
	return helpers.MapSlice(models, ToValueOverrideEntity)
}
