package template_model

import (
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model/template_content_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model/template_layout_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model/template_override_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model/template_rule_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model/template_variant_model"
)

type Template struct {
	Name string

	GameMode            string
	Description         string
	DisplayWinCondition string

	SizeX int
	SizeZ int

	ValueOverrides []ValueOverride

	Orientation *Orientation
	Border      *Border

	GameRules  GameRules
	GlobalBans *GlobalBans

	Variants []Variant

	ZoneLayouts        []ZoneLayoutDef
	MandatoryContent   []MandatoryContent
	ContentCountLimits []ContentCountLimit
	ContentPools       []ContentPool
	ContentLists       []ContentList
}

func (this Template) Clone() Template {
	clone := this
	clone.ValueOverrides = slices.Clone(this.ValueOverrides)
	clone.Orientation = helpers.ClonePointer(this.Orientation)
	clone.Border = helpers.MapPointer(this.Border, Border.Clone)
	clone.GameRules = this.GameRules.Clone()
	clone.GlobalBans = helpers.MapPointer(this.GlobalBans, GlobalBans.Clone)
	clone.Variants = helpers.MapSlice(this.Variants, Variant.Clone)
	clone.ZoneLayouts = helpers.MapSlice(this.ZoneLayouts, ZoneLayoutDef.Clone)
	clone.MandatoryContent = helpers.MapSlice(this.MandatoryContent, MandatoryContent.Clone)
	clone.ContentCountLimits = helpers.MapSlice(this.ContentCountLimits, ContentCountLimit.Clone)
	clone.ContentPools = helpers.MapSlice(this.ContentPools, ContentPool.Clone)
	clone.ContentLists = helpers.MapSlice(this.ContentLists, ContentList.Clone)
	return clone
}

func ToTemplateModel(entity template.RmgTemplate) Template {
	return Template{
		Name:                entity.Name,
		GameMode:            entity.GameMode,
		Description:         entity.Description,
		DisplayWinCondition: entity.DisplayWinCondition,
		SizeX:               entity.SizeX,
		SizeZ:               entity.SizeZ,
		ValueOverrides:      template_override_model.ToValueOverrideModels(entity.ValueOverrides),
		Orientation:         helpers.MapPointer(entity.Orientation, template_variant_model.ToOrientationModel),
		Border:              helpers.MapPointer(entity.Border, template_variant_model.ToBorderModel),
		GameRules:           template_rule_model.ToGameRulesModel(entity.GameRules),
		GlobalBans:          helpers.MapPointer(entity.GlobalBans, template_rule_model.ToGlobalBansModel),
		Variants:            template_variant_model.ToVariantModels(entity.Variants),
		ZoneLayouts:         template_layout_model.ToZoneLayoutDefModels(entity.ZoneLayouts),
		MandatoryContent:    template_content_model.ToMandatoryContentModels(entity.MandatoryContent),
		ContentCountLimits:  template_content_model.ToContentCountLimitModels(entity.ContentCountLimits),
		ContentPools:        template_content_model.ToContentPoolModels(entity.ContentPools),
		ContentLists:        template_content_model.ToContentListModels(entity.ContentLists),
	}
}

func ToTemplateEntity(model Template) template.RmgTemplate {
	return template.RmgTemplate{
		Name:                model.Name,
		GameMode:            model.GameMode,
		Description:         model.Description,
		DisplayWinCondition: model.DisplayWinCondition,
		SizeX:               model.SizeX,
		SizeZ:               model.SizeZ,
		ValueOverrides:      template_override_model.ToValueOverrideEntities(model.ValueOverrides),
		Orientation:         helpers.MapPointer(model.Orientation, template_variant_model.ToOrientationEntity),
		Border:              helpers.MapPointer(model.Border, template_variant_model.ToBorderEntity),
		GameRules:           template_rule_model.ToGameRulesEntity(model.GameRules),
		GlobalBans:          helpers.MapPointer(model.GlobalBans, template_rule_model.ToGlobalBansEntity),
		Variants:            template_variant_model.ToVariantEntities(model.Variants),
		ZoneLayouts:         template_layout_model.ToZoneLayoutDefEntities(model.ZoneLayouts),
		MandatoryContent:    template_content_model.ToMandatoryContentEntities(model.MandatoryContent),
		ContentCountLimits:  template_content_model.ToContentCountLimitEntities(model.ContentCountLimits),
		ContentPools:        template_content_model.ToContentPoolEntities(model.ContentPools),
		ContentLists:        template_content_model.ToContentListEntities(model.ContentLists),
	}
}
