package variant_content

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
)

type VariantBuilder struct {
	item entities.Variant
}

func NewVariantBuilder() *VariantBuilder { return &VariantBuilder{item: entities.Variant{}} }

func (this *VariantBuilder) WithOrientation(orientation entities.Orientation) *VariantBuilder {
	this.item.Orientation = orientation
	return this
}

func (this *VariantBuilder) WithBorder(border entities.Border) *VariantBuilder {
	this.item.Border = border
	return this
}

func (this *VariantBuilder) WithZones(zones ...entities.Zone) *VariantBuilder {
	this.item.Zones = append(this.item.Zones, zones...)
	return this
}

func (this *VariantBuilder) WithConnections(connections ...entities.Connection) *VariantBuilder {
	this.item.Connections = append(this.item.Connections, connections...)
	return this
}

func (this *VariantBuilder) Build() entities.Variant { return this.item }
