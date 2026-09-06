package template_layout_model

import (
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
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

func (this ZoneLayoutDef) Clone() ZoneLayoutDef {
	clone := this
	clone.ElevationModes = slices.Clone(this.ElevationModes)
	clone.GuardedEncounterResourceFractions = this.GuardedEncounterResourceFractions.Clone()
	clone.AmbientPickupDistribution = this.AmbientPickupDistribution.Clone()
	return clone
}

func ToZoneLayoutDefModel(entity template_entity.ZoneLayoutDef) ZoneLayoutDef {
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

func ToZoneLayoutDefEntity(model ZoneLayoutDef) template_entity.ZoneLayoutDef {
	return template_entity.ZoneLayoutDef{
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

func ToZoneLayoutDefModels(entities []template_entity.ZoneLayoutDef) []ZoneLayoutDef {
	return helpers.MapSlice(entities, ToZoneLayoutDefModel)
}

func ToZoneLayoutDefEntities(models []ZoneLayoutDef) []template_entity.ZoneLayoutDef {
	return helpers.MapSlice(models, ToZoneLayoutDefEntity)
}
