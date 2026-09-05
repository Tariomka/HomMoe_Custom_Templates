package template_common_model

import (
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
)

type PlacementRule struct {
	template.PlacementRule
}

func (this PlacementRule) Clone() PlacementRule {
	clone := this
	clone.Args = slices.Clone(this.Args)
	return clone
}

func ToPlacementRuleModel(entity template.PlacementRule) PlacementRule {
	return PlacementRule{PlacementRule: entity}
}

func ToPlacementRuleEntity(model PlacementRule) template.PlacementRule {
	return model.PlacementRule
}

func ToPlacementRuleModels(entities []template.PlacementRule) []PlacementRule {
	return helpers.MapSlice(entities, ToPlacementRuleModel)
}

func ToPlacementRuleEntities(models []PlacementRule) []template.PlacementRule {
	return helpers.MapSlice(models, ToPlacementRuleEntity)
}
