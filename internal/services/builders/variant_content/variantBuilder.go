package variant_content

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

type VariantBuilder struct {
	item template_model.Variant
}

func NewVariantBuilder() *VariantBuilder { return &VariantBuilder{item: template_model.Variant{}} }

func (this *VariantBuilder) WithOrientation(orientation template_model.Orientation) *VariantBuilder {
	this.item.Orientation = orientation
	return this
}

func (this *VariantBuilder) WithBorder(border template_model.Border) *VariantBuilder {
	this.item.Border = border
	return this
}

func (this *VariantBuilder) WithZones(zones ...template_model.Zone) *VariantBuilder {
	this.item.Zones = append(this.item.Zones, zones...)
	return this
}

func (this *VariantBuilder) WithConnections(connections ...template_model.Connection) *VariantBuilder {
	this.item.Connections = append(this.item.Connections, connections...)
	return this
}

func (this *VariantBuilder) Build() template_model.Variant { return this.item }
