package template_common_model

import (
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
)

type PlacementRule struct {
	template_entity.PlacementRule
}

func (this PlacementRule) Clone() PlacementRule {
	clone := this
	clone.Args = slices.Clone(this.Args)
	return clone
}

func ToPlacementRuleModels(entities []template_entity.PlacementRule) []PlacementRule {
	return helpers.MapSlice(entities, func(entity template_entity.PlacementRule) PlacementRule {
		return PlacementRule{PlacementRule: entity}
	})
}

func ToPlacementRuleEntities(models []PlacementRule) []template_entity.PlacementRule {
	return helpers.MapSlice(models, func(model PlacementRule) template_entity.PlacementRule {
		return model.PlacementRule
	})
}
