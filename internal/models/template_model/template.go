// Package template_model is the service-layer mirror of the .rmg.json schema in
// internal/entities/template. The entity is the wire format and nothing else;
// this package owns the structure the application actually works with, which is
// what lets a zone carry its tier without touching the protected schema.
//
// Every type is reachable from here - see types.go - so no caller ever names a
// *_model subpackage. The depguard rule template-model-inner-private enforces
// that.
package template_model

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model/template_content_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model/template_layout_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model/template_override_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model/template_rule_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model/template_variant_model"
)

// Template mirrors template.RmgTemplate; it deliberately does not embed it,
// because every composite field is re-typed.
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
