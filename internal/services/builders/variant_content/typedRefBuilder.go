package variant_content

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
)

type TypedRefBuilder struct {
	item entities.TypedRef
}

func NewRefBuilder() *TypedRefBuilder { return &TypedRefBuilder{item: entities.TypedRef{}} }

func (this *TypedRefBuilder) WithType(refType string) *TypedRefBuilder {
	this.item.Type = refType
	return this
}

func (this *TypedRefBuilder) WithArgs(args ...string) *TypedRefBuilder {
	this.item.Args = append(this.item.Args, args...)
	return this
}

func (this *TypedRefBuilder) Build() entities.TypedRef { return this.item }

func (this *TypedRefBuilder) BuildMainObjectType(args ...string) entities.TypedRef {
	return this.WithType(registry.GetRoadConnectionTypeValues().MainObject).WithArgs(args...).Build()
}

func (this *TypedRefBuilder) BuildConnectionType(args ...string) entities.TypedRef {
	return this.WithType(registry.GetRoadConnectionTypeValues().Connection).WithArgs(args...).Build()
}

func (this *TypedRefBuilder) BuildMandatoryContentType(args ...string) entities.TypedRef {
	return this.WithType(registry.GetRoadConnectionTypeValues().MandatoryContent).WithArgs(args...).Build()
}

func (this *TypedRefBuilder) BuildBiomeMatchZoneType(args ...string) entities.TypedRef {
	return this.WithType(registry.GetBiomeTypeValues().MatchZone).WithArgs(args...).Build()
}

func (this *TypedRefBuilder) BuildBiomeMatchMainObjectType(args ...string) entities.TypedRef {
	return this.WithType(registry.GetBiomeTypeValues().MatchMainObject).WithArgs(args...).Build()
}

func (this *TypedRefBuilder) BuildBiomeFromListType(args ...string) entities.TypedRef {
	return this.WithType(registry.GetBiomeTypeValues().FromList).WithArgs(args...).Build()
}
