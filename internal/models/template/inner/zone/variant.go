package zone

import "github.com/Tariomka/hommoe_custom_templates/internal/models/template/inner"

// Variant represents a single layout variant of the map.
// The template root has a `variants` array - each variant has its own zones/connections,
// orientation, and border. Connections live inside the variant (not the template root).
type Variant struct {
	Orientation inner.Orientation  `json:"orientation"`
	Border      inner.Border       `json:"border"`
	Zones       []Zone             `json:"zones"`
	Connections []inner.Connection `json:"connections"`
}
