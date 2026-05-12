package zonelayout

// GuardedEncounterResourceFractions configures resource splits for guarded encounters.
type GuardedEncounterResourceFractions struct {
	CountBounds []int     `json:"countBounds"`
	Fractions   []float64 `json:"fractions"`
}
