package template

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

// ElevationMode is one weighted elevation band used by zone generation.
type ElevationMode struct {
	Weight              int     `json:"weight"`
	MinElevatedFraction float64 `json:"minElevatedFraction"`
	MaxElevatedFraction float64 `json:"maxElevatedFraction"`
}

// GuardedEncounterResourceFractions configures resource splits for guarded encounters.
type GuardedEncounterResourceFractions struct {
	CountBounds []int     `json:"countBounds"`
	Fractions   []float64 `json:"fractions"`
}

// AmbientPickupDistribution configures distribution of ambient pickups in a zone.
type AmbientPickupDistribution struct {
	Repulsion          float64 `json:"repulsion"`
	Noise              float64 `json:"noise"`
	RoadAttraction     float64 `json:"roadAttraction"`
	ObstacleAttraction float64 `json:"obstacleAttraction"`
	GroupSizeWeights   []int   `json:"groupSizeWeights"`
}
