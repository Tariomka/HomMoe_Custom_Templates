package zone

// Variant represents a single layout variant of the map.
// The template root has a `variants` array - each variant has its own zones/connections,
// orientation, and border. Connections live inside the variant (not the template root).
type Variant struct {
	Orientation Orientation  `json:"orientation"`
	Border      Border       `json:"border"`
	Zones       []Zone       `json:"zones"`
	Connections []Connection `json:"connections"`
}
