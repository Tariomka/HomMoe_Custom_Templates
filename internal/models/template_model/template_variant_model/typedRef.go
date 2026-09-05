package template_variant_model

import (
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template"
)

type TypedRef struct{ template.TypedRef }

func (this TypedRef) Clone() TypedRef {
	clone := this
	clone.Args = slices.Clone(this.Args)
	return clone
}

func ToTypedRefModel(entity template.TypedRef) TypedRef {
	return TypedRef{TypedRef: entity}
}

func ToTypedRefEntity(model TypedRef) template.TypedRef {
	return model.TypedRef
}
