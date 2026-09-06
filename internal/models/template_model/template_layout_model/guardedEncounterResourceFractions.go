package template_layout_model

import (
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
)

type GuardedEncounterResourceFractions struct {
	template_entity.GuardedEncounterResourceFractions
}

func (this GuardedEncounterResourceFractions) Clone() GuardedEncounterResourceFractions {
	clone := this
	clone.CountBounds = slices.Clone(this.CountBounds)
	clone.Fractions = slices.Clone(this.Fractions)
	return clone
}

func ToGuardedEncounterResourceFractionsModel(
	entity template_entity.GuardedEncounterResourceFractions) GuardedEncounterResourceFractions {
	return GuardedEncounterResourceFractions{GuardedEncounterResourceFractions: entity}
}

func ToGuardedEncounterResourceFractionsEntity(
	model GuardedEncounterResourceFractions) template_entity.GuardedEncounterResourceFractions {
	return model.GuardedEncounterResourceFractions
}

func ToGuardedEncounterResourceFractionsModels(
	entities []template_entity.GuardedEncounterResourceFractions) []GuardedEncounterResourceFractions {
	return helpers.MapSlice(entities, ToGuardedEncounterResourceFractionsModel)
}

func ToGuardedEncounterResourceFractionsEntities(
	models []GuardedEncounterResourceFractions) []template_entity.GuardedEncounterResourceFractions {
	return helpers.MapSlice(models, ToGuardedEncounterResourceFractionsEntity)
}
