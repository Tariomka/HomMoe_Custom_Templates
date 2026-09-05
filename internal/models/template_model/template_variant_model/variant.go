package template_variant_model

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
)

type Variant struct {
	Orientation Orientation
	Border      Border
	Zones       []Zone
	Connections []Connection
}

func (this Variant) Clone() Variant {
	return Variant{
		Orientation: this.Orientation,
		Border:      this.Border.Clone(),
		Zones:       helpers.MapSlice(this.Zones, Zone.Clone),
		Connections: helpers.MapSlice(this.Connections, Connection.Clone),
	}
}

func ToVariantModel(entity template.Variant) Variant {
	return Variant{
		Orientation: ToOrientationModel(entity.Orientation),
		Border:      ToBorderModel(entity.Border),
		Zones:       ToZoneModels(entity.Zones),
		Connections: ToConnectionModels(entity.Connections),
	}
}

func ToVariantEntity(model Variant) template.Variant {
	return template.Variant{
		Orientation: ToOrientationEntity(model.Orientation),
		Border:      ToBorderEntity(model.Border),
		Zones:       ToZoneEntities(model.Zones),
		Connections: ToConnectionEntities(model.Connections),
	}
}

func ToVariantModels(entities []template.Variant) []Variant {
	return helpers.MapSlice(entities, ToVariantModel)
}

func ToVariantEntities(models []Variant) []template.Variant {
	return helpers.MapSlice(models, ToVariantEntity)
}
