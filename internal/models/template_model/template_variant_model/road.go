package template_variant_model

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
)

type Road struct {
	Type string
	From TypedRef
	To   TypedRef

	Road                 *bool
	SimTurnSquad         bool
	GuardValue           int
	GuardWeeklyIncrement float64
}

func (this Road) Clone() Road {
	clone := this
	clone.From = this.From.Clone()
	clone.To = this.To.Clone()
	clone.Road = helpers.ClonePointer(this.Road)
	return clone
}

func ToRoadModel(entity template.Road) Road {
	return Road{
		Type:                 entity.Type,
		From:                 ToTypedRefModel(entity.From),
		To:                   ToTypedRefModel(entity.To),
		Road:                 entity.Road,
		SimTurnSquad:         entity.SimTurnSquad,
		GuardValue:           entity.GuardValue,
		GuardWeeklyIncrement: entity.GuardWeeklyIncrement,
	}
}

func ToRoadEntity(model Road) template.Road {
	return template.Road{
		Type:                 model.Type,
		From:                 ToTypedRefEntity(model.From),
		To:                   ToTypedRefEntity(model.To),
		Road:                 model.Road,
		SimTurnSquad:         model.SimTurnSquad,
		GuardValue:           model.GuardValue,
		GuardWeeklyIncrement: model.GuardWeeklyIncrement,
	}
}

func ToRoadModels(entities []template.Road) []Road {
	return helpers.MapSlice(entities, ToRoadModel)
}

func ToRoadEntities(models []Road) []template.Road {
	return helpers.MapSlice(models, ToRoadEntity)
}
