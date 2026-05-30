package zone

// EncounterHolesSettings overrides the global encounter-holes parameters for a single zone.
type EncounterHolesSettings struct {
	AffectedEncounters float64 `json:"affectedEncounters"`
	TwoHoleEncounters  float64 `json:"twoHoleEncounters"`
}
