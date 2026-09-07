package template_layout_model

import (
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
)

type AmbientPickupDistribution struct {
	template_entity.AmbientPickupDistribution
}

func (this AmbientPickupDistribution) Clone() AmbientPickupDistribution {
	clone := this
	clone.GroupSizeWeights = slices.Clone(this.GroupSizeWeights)
	return clone
}
