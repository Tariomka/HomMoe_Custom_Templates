package template_rule_model

import (
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template"
)

type WinConditions struct{ template.WinConditions }

func (this WinConditions) Clone() WinConditions {
	clone := this
	clone.TournamentDays = slices.Clone(this.TournamentDays)
	clone.TournamentAnnounceDays = slices.Clone(this.TournamentAnnounceDays)
	return clone
}

func ToWinConditionsModel(entity template.WinConditions) WinConditions {
	return WinConditions{WinConditions: entity}
}

func ToWinConditionsEntity(model WinConditions) template.WinConditions {
	return model.WinConditions
}
