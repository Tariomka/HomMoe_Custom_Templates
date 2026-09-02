package template_content_model

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model/template_common_model"
)

type MandatoryContentItem struct {
	SID       string
	Name      string
	IsMine    bool
	IsGuarded bool
	Rules     []template_common_model.PlacementRule

	Variant *int

	Owner string

	GuardValue int

	IncludeLists []string

	Content []WeightedContent

	DesignatedEncounter *bool

	SoloEncounter bool

	Road *bool
}

func ToMandatoryContentItemModel(entity template.MandatoryContentItem) MandatoryContentItem {
	return MandatoryContentItem{
		SID:                 entity.SID,
		Name:                entity.Name,
		IsMine:              entity.IsMine,
		IsGuarded:           entity.IsGuarded,
		Rules:               template_common_model.ToPlacementRuleModels(entity.Rules),
		Variant:             entity.Variant,
		Owner:               entity.Owner,
		GuardValue:          entity.GuardValue,
		IncludeLists:        entity.IncludeLists,
		Content:             ToWeightedContentModels(entity.Content),
		DesignatedEncounter: entity.DesignatedEncounter,
		SoloEncounter:       entity.SoloEncounter,
		Road:                entity.Road,
	}
}

func ToMandatoryContentItemEntity(model MandatoryContentItem) template.MandatoryContentItem {
	return template.MandatoryContentItem{
		SID:                 model.SID,
		Name:                model.Name,
		IsMine:              model.IsMine,
		IsGuarded:           model.IsGuarded,
		Rules:               template_common_model.ToPlacementRuleEntities(model.Rules),
		Variant:             model.Variant,
		Owner:               model.Owner,
		GuardValue:          model.GuardValue,
		IncludeLists:        model.IncludeLists,
		Content:             ToWeightedContentEntities(model.Content),
		DesignatedEncounter: model.DesignatedEncounter,
		SoloEncounter:       model.SoloEncounter,
		Road:                model.Road,
	}
}

func ToMandatoryContentItemModels(entities []template.MandatoryContentItem) []MandatoryContentItem {
	return helpers.MapSlice(entities, ToMandatoryContentItemModel)
}

func ToMandatoryContentItemEntities(models []MandatoryContentItem) []template.MandatoryContentItem {
	return helpers.MapSlice(models, ToMandatoryContentItemEntity)
}
