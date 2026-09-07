package template_variant_model

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
)

type EncounterHolesSettings struct {
	template_entity.EncounterHolesSettings
}

func ToEncounterHolesSettingsModel(entity template_entity.EncounterHolesSettings) EncounterHolesSettings {
	return EncounterHolesSettings{EncounterHolesSettings: entity}
}

func ToEncounterHolesSettingsEntity(model EncounterHolesSettings) template_entity.EncounterHolesSettings {
	return model.EncounterHolesSettings
}
