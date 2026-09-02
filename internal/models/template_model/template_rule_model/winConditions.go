package template_rule_model

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template"
)

type WinConditions struct{ template.WinConditions }

func ToWinConditionsModel(entity template.WinConditions) WinConditions {
	return WinConditions{WinConditions: entity}
}

func ToWinConditionsEntity(model WinConditions) template.WinConditions {
	return model.WinConditions
}
