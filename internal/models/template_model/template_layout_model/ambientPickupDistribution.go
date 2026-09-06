package template_layout_model

import (
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
)

type AmbientPickupDistribution struct {
	template_entity.AmbientPickupDistribution
}

func (this AmbientPickupDistribution) Clone() AmbientPickupDistribution {
	clone := this
	clone.GroupSizeWeights = slices.Clone(this.GroupSizeWeights)
	return clone
}

func ToAmbientPickupDistributionModel(entity template_entity.AmbientPickupDistribution) AmbientPickupDistribution {
	return AmbientPickupDistribution{AmbientPickupDistribution: entity}
}

func ToAmbientPickupDistributionEntity(model AmbientPickupDistribution) template_entity.AmbientPickupDistribution {
	return model.AmbientPickupDistribution
}

func ToAmbientPickupDistributionModels(
	entities []template_entity.AmbientPickupDistribution) []AmbientPickupDistribution {
	return helpers.MapSlice(entities, ToAmbientPickupDistributionModel)
}

func ToAmbientPickupDistributionEntities(
	models []AmbientPickupDistribution) []template_entity.AmbientPickupDistribution {
	return helpers.MapSlice(models, ToAmbientPickupDistributionEntity)
}
