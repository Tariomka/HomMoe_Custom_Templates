package template_layout_model

import "slices"

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
