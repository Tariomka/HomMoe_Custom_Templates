package mandatory_content

import "strings"

type Distance struct{ Min, Max float64 }

var (
	DistanceNextTo  = Distance{Min: 0.05, Max: 0.1}
	DistanceNear    = Distance{Min: 0.0, Max: 0.35}
	DistanceMedium  = Distance{Min: 0.25, Max: 0.5}
	DistanceFar     = Distance{Min: 0.5, Max: 0.75}
	DistanceVeryFar = Distance{Min: 0.75, Max: 0.9}
)

// TryGetDistanceFrom resolves the road-distance label persisted in a
// ZoneContentRowSave to the matching distancePreset. An empty or unknown
// label means "Any" (no constraint added).
func TryGetDistanceFrom(label string) (Distance, bool) {
	switch strings.TrimSpace(label) {
	case "Next To":
		return DistanceNextTo, true
	case "Near":
		return DistanceNear, true
	case "Medium":
		return DistanceMedium, true
	case "Far":
		return DistanceFar, true
	case "Very Far":
		return DistanceVeryFar, true
	}
	return Distance{}, false
}
