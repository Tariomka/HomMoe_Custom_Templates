package template_variant_model

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template"
)

type EncounterHolesSettings struct {
	template.EncounterHolesSettings
}

func ToEncounterHolesSettingsModel(entity template.EncounterHolesSettings) EncounterHolesSettings {
	return EncounterHolesSettings{EncounterHolesSettings: entity}
}

func ToEncounterHolesSettingsEntity(model EncounterHolesSettings) template.EncounterHolesSettings {
	return model.EncounterHolesSettings
}
