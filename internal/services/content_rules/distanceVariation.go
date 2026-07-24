package content_rules

// DistanceVariation mirrors the C# DistancePresets.DistanceVariation struct:
// a named distance band with fractional min/max bounds.
type DistanceVariation struct {
	Name string
	Min  float64
	Max  float64
}

func defaultDistanceVariation() DistanceVariation {
	return DistanceVariation{Name: "Medium", Min: 0.25, Max: 0.5}
}
