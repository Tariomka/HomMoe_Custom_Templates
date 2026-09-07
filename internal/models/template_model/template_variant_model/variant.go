package template_variant_model

import "github.com/Tariomka/hommoe_custom_templates/internal/helpers"

type Variant struct {
	Orientation Orientation
	Border      Border
	Zones       []Zone
	Connections []Connection
}

func (this Variant) Clone() Variant {
	return Variant{
		Orientation: this.Orientation,
		Border:      this.Border.Clone(),
		Zones:       helpers.MapSlice(this.Zones, Zone.Clone),
		Connections: helpers.MapSlice(this.Connections, Connection.Clone),
	}
}
