package variant_content

import "github.com/Tariomka/hommoe_custom_templates/internal/models/template"

type VariantBuilder struct {
	item template.Variant
}

func NewVariantBuilder() *VariantBuilder { return &VariantBuilder{item: template.Variant{}} }

func (this *VariantBuilder) WithOrientation(orientation template.Orientation) *VariantBuilder {
	this.item.Orientation = orientation
	return this
}
func (this *VariantBuilder) WithBorder(border template.Border) *VariantBuilder {
	this.item.Border = border
	return this
}
func (this *VariantBuilder) WithZones(zones ...template.Zone) *VariantBuilder {
	this.item.Zones = append(this.item.Zones, zones...)
	return this
}
func (this *VariantBuilder) WithConnections(connections ...template.Connection) *VariantBuilder {
	this.item.Connections = append(this.item.Connections, connections...)
	return this
}
func (this *VariantBuilder) Build() template.Variant { return this.item }
