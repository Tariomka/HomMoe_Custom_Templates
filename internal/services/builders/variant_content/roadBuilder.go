package variant_content

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
)

var roadTypes = registry.GetRoadTypeValues()

type RoadBuilder struct {
	item entities.Road
}

func NewRoadBuilder() *RoadBuilder { return &RoadBuilder{item: entities.Road{}} }

func (this *RoadBuilder) WithStoneType() *RoadBuilder {
	this.item.Type = roadTypes.Stone
	return this
}
func (this *RoadBuilder) WithDirtType() *RoadBuilder {
	this.item.Type = roadTypes.Dirt
	return this
}
func (this *RoadBuilder) WithFrom(from entities.TypedRef) *RoadBuilder {
	this.item.From = from
	return this
}
func (this *RoadBuilder) WithTo(to entities.TypedRef) *RoadBuilder {
	this.item.To = to
	return this
}
func (this *RoadBuilder) Build() entities.Road { return this.item }
