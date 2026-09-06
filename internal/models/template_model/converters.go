package template_model

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model/template_variant_model"
)

func ToZoneEntity(model Zone) template_entity.Zone {
	return template_variant_model.ToZoneEntity(model)
}

func ToZoneModel(entity template_entity.Zone) Zone {
	return template_variant_model.ToZoneModel(entity)
}

func ToConnectionEntity(model Connection) template_entity.Connection {
	return template_variant_model.ToConnectionEntity(model)
}

func ToConnectionModel(entity template_entity.Connection) Connection {
	return template_variant_model.ToConnectionModel(entity)
}
