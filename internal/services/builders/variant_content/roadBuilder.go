package variant_content

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
)

type RoadBuilder struct{ item template_model.Road }

func NewRoadBuilder() *RoadBuilder { return &RoadBuilder{item: template_model.Road{}} }

func (this *RoadBuilder) WithStoneType() *RoadBuilder {
	this.item.Type = registry.GetRoadTypeValues().Stone
	return this
}
func (this *RoadBuilder) WithDirtType() *RoadBuilder {
	this.item.Type = registry.GetRoadTypeValues().Dirt
	return this
}

func (this *RoadBuilder) WithFrom(from template_model.TypedRef) *RoadBuilder {
	this.item.From = from
	return this
}

func (this *RoadBuilder) WithTo(to template_model.TypedRef) *RoadBuilder {
	this.item.To = to
	return this
}

func (this *RoadBuilder) Build() template_model.Road { return this.item }
