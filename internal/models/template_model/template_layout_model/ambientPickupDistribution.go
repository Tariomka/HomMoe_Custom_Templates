package template_layout_model

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
)

type AmbientPickupDistribution struct {
	template.AmbientPickupDistribution
}

func ToAmbientPickupDistributionModel(entity template.AmbientPickupDistribution) AmbientPickupDistribution {
	return AmbientPickupDistribution{AmbientPickupDistribution: entity}
}

func ToAmbientPickupDistributionEntity(model AmbientPickupDistribution) template.AmbientPickupDistribution {
	return model.AmbientPickupDistribution
}

func ToAmbientPickupDistributionModels(
	entities []template.AmbientPickupDistribution) []AmbientPickupDistribution {
	return helpers.MapSlice(entities, ToAmbientPickupDistributionModel)
}

func ToAmbientPickupDistributionEntities(
	models []AmbientPickupDistribution) []template.AmbientPickupDistribution {
	return helpers.MapSlice(models, ToAmbientPickupDistributionEntity)
}
