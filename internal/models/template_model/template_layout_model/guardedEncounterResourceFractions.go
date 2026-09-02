package template_layout_model

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
)

type GuardedEncounterResourceFractions struct {
	template.GuardedEncounterResourceFractions
}

func ToGuardedEncounterResourceFractionsModel(
	entity template.GuardedEncounterResourceFractions) GuardedEncounterResourceFractions {
	return GuardedEncounterResourceFractions{GuardedEncounterResourceFractions: entity}
}

func ToGuardedEncounterResourceFractionsEntity(
	model GuardedEncounterResourceFractions) template.GuardedEncounterResourceFractions {
	return model.GuardedEncounterResourceFractions
}

func ToGuardedEncounterResourceFractionsModels(
	entities []template.GuardedEncounterResourceFractions) []GuardedEncounterResourceFractions {
	return helpers.MapSlice(entities, ToGuardedEncounterResourceFractionsModel)
}

func ToGuardedEncounterResourceFractionsEntities(
	models []GuardedEncounterResourceFractions) []template.GuardedEncounterResourceFractions {
	return helpers.MapSlice(models, ToGuardedEncounterResourceFractionsEntity)
}
