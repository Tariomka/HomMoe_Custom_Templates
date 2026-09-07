package template_variant_model

import (
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
)

type TypedRef struct{ template_entity.TypedRef }

func (this TypedRef) Clone() TypedRef {
	clone := this
	clone.Args = slices.Clone(this.Args)
	return clone
}

func ToTypedRefModel(entity template_entity.TypedRef) TypedRef {
	return TypedRef{TypedRef: entity}
}

func ToTypedRefEntity(model TypedRef) template_entity.TypedRef {
	return model.TypedRef
}
