package variant_content

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
)

type TypedRefBuilder struct {
	item template_model.TypedRef
}

func NewRefBuilder() *TypedRefBuilder { return &TypedRefBuilder{item: template_model.TypedRef{}} }

func (this *TypedRefBuilder) WithType(refType string) *TypedRefBuilder {
	this.item.Type = refType
	return this
}

func (this *TypedRefBuilder) WithArgs(args ...string) *TypedRefBuilder {
	this.item.Args = append(this.item.Args, args...)
	return this
}

func (this *TypedRefBuilder) Build() template_model.TypedRef { return this.item }

func (this *TypedRefBuilder) BuildMainObjectType(args ...string) template_model.TypedRef {
	return this.WithType(registry.GetRoadConnectionTypeValues().MainObject).WithArgs(args...).Build()
}

func (this *TypedRefBuilder) BuildConnectionType(args ...string) template_model.TypedRef {
	return this.WithType(registry.GetRoadConnectionTypeValues().Connection).WithArgs(args...).Build()
}

func (this *TypedRefBuilder) BuildMandatoryContentType(args ...string) template_model.TypedRef {
	return this.WithType(registry.GetRoadConnectionTypeValues().MandatoryContent).WithArgs(args...).Build()
}

func (this *TypedRefBuilder) BuildBiomeMatchZoneType(args ...string) template_model.TypedRef {
	return this.WithType(registry.GetBiomeTypeValues().MatchZone).WithArgs(args...).Build()
}

func (this *TypedRefBuilder) BuildBiomeMatchMainObjectType(args ...string) template_model.TypedRef {
	return this.WithType(registry.GetBiomeTypeValues().MatchMainObject).WithArgs(args...).Build()
}

func (this *TypedRefBuilder) BuildBiomeFromListType(args ...string) template_model.TypedRef {
	return this.WithType(registry.GetBiomeTypeValues().FromList).WithArgs(args...).Build()
}
