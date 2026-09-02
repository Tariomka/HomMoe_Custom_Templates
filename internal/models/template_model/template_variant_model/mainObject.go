package template_variant_model

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
)

type MainObject struct {
	Type string

	Spawn string
	Owner string

	RemoveGuardIfHasOwner bool

	GuardChance          float64
	GuardValue           int
	GuardRandomization   float64
	GuardWeeklyIncrement float64

	BuildingsConstructionSid string

	Faction  *TypedRef
	Factions []string

	Placement     string
	PlacementArgs []string

	HoldCityWinCon bool

	IsKeyObject bool

	EnableWeeklyUnitIncrement bool
	InitialUnitIncrement      int
}

func ToMainObjectModel(entity template.MainObject) MainObject {
	return MainObject{
		Type:                      entity.Type,
		Spawn:                     entity.Spawn,
		Owner:                     entity.Owner,
		RemoveGuardIfHasOwner:     entity.RemoveGuardIfHasOwner,
		GuardChance:               entity.GuardChance,
		GuardValue:                entity.GuardValue,
		GuardRandomization:        entity.GuardRandomization,
		GuardWeeklyIncrement:      entity.GuardWeeklyIncrement,
		BuildingsConstructionSid:  entity.BuildingsConstructionSid,
		Faction:                   helpers.MapPointer(entity.Faction, ToTypedRefModel),
		Factions:                  entity.Factions,
		Placement:                 entity.Placement,
		PlacementArgs:             entity.PlacementArgs,
		HoldCityWinCon:            entity.HoldCityWinCon,
		IsKeyObject:               entity.IsKeyObject,
		EnableWeeklyUnitIncrement: entity.EnableWeeklyUnitIncrement,
		InitialUnitIncrement:      entity.InitialUnitIncrement,
	}
}

func ToMainObjectEntity(model MainObject) template.MainObject {
	return template.MainObject{
		Type:                      model.Type,
		Spawn:                     model.Spawn,
		Owner:                     model.Owner,
		RemoveGuardIfHasOwner:     model.RemoveGuardIfHasOwner,
		GuardChance:               model.GuardChance,
		GuardValue:                model.GuardValue,
		GuardRandomization:        model.GuardRandomization,
		GuardWeeklyIncrement:      model.GuardWeeklyIncrement,
		BuildingsConstructionSid:  model.BuildingsConstructionSid,
		Faction:                   helpers.MapPointer(model.Faction, ToTypedRefEntity),
		Factions:                  model.Factions,
		Placement:                 model.Placement,
		PlacementArgs:             model.PlacementArgs,
		HoldCityWinCon:            model.HoldCityWinCon,
		IsKeyObject:               model.IsKeyObject,
		EnableWeeklyUnitIncrement: model.EnableWeeklyUnitIncrement,
		InitialUnitIncrement:      model.InitialUnitIncrement,
	}
}

func ToMainObjectModels(entities []template.MainObject) []MainObject {
	return helpers.MapSlice(entities, ToMainObjectModel)
}

func ToMainObjectEntities(models []MainObject) []template.MainObject {
	return helpers.MapSlice(models, ToMainObjectEntity)
}
