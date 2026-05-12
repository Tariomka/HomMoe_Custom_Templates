package zonelayout

// ZoneLayoutDef is a named layout preset referenced by Zone.Layout.
// Templates declare these at the root level under `zoneLayouts`.
type ZoneLayoutDef struct {
	Name string `json:"name"`

	ObstaclesFill     float64 `json:"obstaclesFill"`
	ObstaclesFillVoid float64 `json:"obstaclesFillVoid"`
	LakesFill         float64 `json:"lakesFill"`
	MinLakeArea       int     `json:"minLakeArea,omitempty"`

	ElevationClusterScale float64         `json:"elevationClusterScale"`
	ElevationModes        []ElevationMode `json:"elevationModes"`

	RoadClusterArea int `json:"roadClusterArea"`

	GuardedEncounterResourceFractions GuardedEncounterResourceFractions `json:"guardedEncounterResourceFractions"`
	AmbientPickupDistribution         AmbientPickupDistribution         `json:"ambientPickupDistribution"`
}
