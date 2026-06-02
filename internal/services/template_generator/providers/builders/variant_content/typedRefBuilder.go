package variant_content

import "github.com/Tariomka/hommoe_custom_templates/internal/models/template"

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
