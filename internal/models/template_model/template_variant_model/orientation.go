package template_variant_model

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template"
)

type Orientation struct{ template.Orientation }

func ToOrientationModel(entity template.Orientation) Orientation {
	return Orientation{Orientation: entity}
}

func ToOrientationEntity(model Orientation) template.Orientation {
	return model.Orientation
}
