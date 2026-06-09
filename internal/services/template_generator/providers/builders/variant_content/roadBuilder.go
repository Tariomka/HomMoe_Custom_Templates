package variant_content

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
)

var roadTypes = registry.GetRoadTypeValues()

type RoadBuilder struct {
	item template.Road
}

func NewRoadBuilder() *RoadBuilder { return &RoadBuilder{item: template.Road{}} }

func (this *RoadBuilder) WithStoneType() *RoadBuilder {
	this.item.Type = roadTypes.Stone
	return this
}
func (this *RoadBuilder) WithDirtType() *RoadBuilder {
	this.item.Type = roadTypes.Dirt
	return this
}
func (this *RoadBuilder) WithFrom(from template.TypedRef) *RoadBuilder {
	this.item.From = from
	return this
}
func (this *RoadBuilder) WithTo(to template.TypedRef) *RoadBuilder {
	this.item.To = to
	return this
}
func (this *RoadBuilder) Build() template.Road { return this.item }
