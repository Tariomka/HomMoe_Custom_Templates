package variant_content

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
)

var (
	roadConnTypes = registry.GetRoadConnectionTypeValues()
	biomeTypes    = registry.GetBiomeTypeValues()
)

type TypedRefBuilder struct {
	item template.TypedRef
}

func NewRefBuilder() *TypedRefBuilder { return &TypedRefBuilder{item: template.TypedRef{}} }

func (this *TypedRefBuilder) WithType(refType string) *TypedRefBuilder {
	this.item.Type = refType
	return this
}
func (this *TypedRefBuilder) WithArgs(args ...string) *TypedRefBuilder {
	this.item.Args = append(this.item.Args, args...)
	return this
}
func (this *TypedRefBuilder) Build() template.TypedRef { return this.item }

func (this *TypedRefBuilder) BuildMainObjectType(args ...string) template.TypedRef {
	return this.WithType(roadConnTypes.MainObject).WithArgs(args...).Build()
}

func (this *TypedRefBuilder) BuildConnectionType(args ...string) template.TypedRef {
	return this.WithType(roadConnTypes.Connection).WithArgs(args...).Build()
}

func (this *TypedRefBuilder) BuildMandatoryContentType(args ...string) template.TypedRef {
	return this.WithType(roadConnTypes.MandatoryContent).WithArgs(args...).Build()
}

func (this *TypedRefBuilder) BuildBiomeMatchZoneType(args ...string) template.TypedRef {
	return this.WithType(biomeTypes.MatchZone).WithArgs(args...).Build()
}

func (this *TypedRefBuilder) BuildBiomeMatchMainObjectType(args ...string) template.TypedRef {
	return this.WithType(biomeTypes.MatchMainObject).WithArgs(args...).Build()
}

func (this *TypedRefBuilder) BuildBiomeFromListType(args ...string) template.TypedRef {
	return this.WithType(biomeTypes.FromList).WithArgs(args...).Build()
}
