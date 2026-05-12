package zone_layout

// ElevationMode is one weighted elevation band used by zone generation.
type ElevationMode struct {
	Weight              int     `json:"weight"`
	MinElevatedFraction float64 `json:"minElevatedFraction"`
	MaxElevatedFraction float64 `json:"maxElevatedFraction"`
}
