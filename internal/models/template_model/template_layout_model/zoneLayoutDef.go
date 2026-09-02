package template_layout_model

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
)

type ZoneLayoutDef struct {
	Name string

	ObstaclesFill     float64
	ObstaclesFillVoid float64
	LakesFill         float64
	MinLakeArea       int

	ElevationClusterScale float64
	ElevationModes        []ElevationMode

	RoadClusterArea int

	GuardedEncounterResourceFractions GuardedEncounterResourceFractions
	AmbientPickupDistribution         AmbientPickupDistribution
}

func ToZoneLayoutDefModel(entity template.ZoneLayoutDef) ZoneLayoutDef {
	return ZoneLayoutDef{
		Name:                  entity.Name,
		ObstaclesFill:         entity.ObstaclesFill,
		ObstaclesFillVoid:     entity.ObstaclesFillVoid,
		LakesFill:             entity.LakesFill,
		MinLakeArea:           entity.MinLakeArea,
		ElevationClusterScale: entity.ElevationClusterScale,
		ElevationModes:        ToElevationModeModels(entity.ElevationModes),
		RoadClusterArea:       entity.RoadClusterArea,
		GuardedEncounterResourceFractions: ToGuardedEncounterResourceFractionsModel(
			entity.GuardedEncounterResourceFractions,
		),
		AmbientPickupDistribution: ToAmbientPickupDistributionModel(entity.AmbientPickupDistribution),
	}
}

func ToZoneLayoutDefEntity(model ZoneLayoutDef) template.ZoneLayoutDef {
	return template.ZoneLayoutDef{
		Name:                  model.Name,
		ObstaclesFill:         model.ObstaclesFill,
		ObstaclesFillVoid:     model.ObstaclesFillVoid,
		LakesFill:             model.LakesFill,
		MinLakeArea:           model.MinLakeArea,
		ElevationClusterScale: model.ElevationClusterScale,
		ElevationModes:        ToElevationModeEntities(model.ElevationModes),
		RoadClusterArea:       model.RoadClusterArea,
		GuardedEncounterResourceFractions: ToGuardedEncounterResourceFractionsEntity(
			model.GuardedEncounterResourceFractions,
		),
		AmbientPickupDistribution: ToAmbientPickupDistributionEntity(model.AmbientPickupDistribution),
	}
}

func ToZoneLayoutDefModels(entities []template.ZoneLayoutDef) []ZoneLayoutDef {
	return helpers.MapSlice(entities, ToZoneLayoutDefModel)
}

func ToZoneLayoutDefEntities(models []ZoneLayoutDef) []template.ZoneLayoutDef {
	return helpers.MapSlice(models, ToZoneLayoutDefEntity)
}
