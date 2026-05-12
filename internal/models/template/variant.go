package template

// Variant represents a single map variant
type Variant struct {
	Orientation Orientation  `json:"orientation"`
	Border      Border       `json:"border,omitempty"`
	Zones       []Zone       `json:"zones"`
	Connections []Connection `json:"connections"`
}
