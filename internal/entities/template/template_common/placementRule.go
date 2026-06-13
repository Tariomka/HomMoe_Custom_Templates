package template_common

// PlacementRule biases placement of a MandatoryContentItem (or portal endpoint) towards a reference point.
// Observed `type` values include: "Crossroads", "Road", "MainObject", "Sid".
// `Args` is loosely-typed (most rules use strings; "Pyramid" / "Sand Clover" pass numeric thresholds).
// A rule may specify either an inclusive [targetMin, targetMax] range OR a single `target` value.
type PlacementRule struct {
	Type      string  `json:"type"`
	Args      []any   `json:"args"`
	TargetMin float64 `json:"targetMin,omitempty"`
	TargetMax float64 `json:"targetMax,omitempty"`
	Target    float64 `json:"target,omitempty"`
	Weight    int     `json:"weight"`
}
