package template_rule_model

import (
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
)

type WinConditions struct{ template_entity.WinConditions }

func (this WinConditions) Clone() WinConditions {
	clone := this
	clone.TournamentDays = slices.Clone(this.TournamentDays)
	clone.TournamentAnnounceDays = slices.Clone(this.TournamentAnnounceDays)
	return clone
}

func ToWinConditionsModel(entity template_entity.WinConditions) WinConditions {
	return WinConditions{WinConditions: entity}
}

func ToWinConditionsEntity(model WinConditions) template_entity.WinConditions {
	return model.WinConditions
}
