package template_variant_model

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
)

type Orientation struct{ template_entity.Orientation }

func ToOrientationModel(entity template_entity.Orientation) Orientation {
	return Orientation{Orientation: entity}
}

func ToOrientationEntity(model Orientation) template_entity.Orientation {
	return model.Orientation
}
