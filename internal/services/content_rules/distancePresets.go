package content_rules

import "strings"

// DistanceVariation mirrors the C# DistancePresets.DistanceVariation struct:
// a named distance band with fractional min/max bounds.
type DistanceVariation struct {
	Name string
	Min  float64
	Max  float64
}

// Distance presets, matching the C# DistancePresets values exactly.
var (
	DistanceNextTo  = DistanceVariation{Name: "Next To", Min: 0.05, Max: 0.1}
	DistanceNear    = DistanceVariation{Name: "Near", Min: 0.1, Max: 0.25}
	DistanceMedium  = DistanceVariation{Name: "Medium", Min: 0.25, Max: 0.5}
	DistanceFar     = DistanceVariation{Name: "Far", Min: 0.5, Max: 0.75}
	DistanceVeryFar = DistanceVariation{Name: "Very Far", Min: 0.75, Max: 0.9}
)

// allDistanceVariations preserves declaration order to mirror the field
// ordering the C# reflection-based GetDisplayNames relies on.
var allDistanceVariations = []DistanceVariation{
	DistanceNextTo,
	DistanceNear,
	DistanceMedium,
	DistanceFar,
	DistanceVeryFar,
}

// GetDistanceDisplayNames returns every variation's display name in preset order.
func GetDistanceDisplayNames() []string {
	names := make([]string, len(allDistanceVariations))
	for i, variation := range allDistanceVariations {
		names[i] = variation.Name
	}
	return names
}

// GetDistanceVariationByName resolves a display name to its variation.
func GetDistanceVariationByName(name string) (DistanceVariation, bool) {
	trimmed := strings.TrimSpace(name)
	for _, variation := range allDistanceVariations {
		if strings.EqualFold(variation.Name, trimmed) {
			return variation, true
		}
	}
	return DistanceVariation{}, false
}
