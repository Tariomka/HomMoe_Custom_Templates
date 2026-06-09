package variant

// Noise describes a single noise layer applied to borders/water.
// `amp` may be integer or fractional across the example templates.
type Noise struct {
	Amplitude float64 `json:"amp"`
	Frequency int     `json:"freq"`
}
