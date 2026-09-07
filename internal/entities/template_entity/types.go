package template_entity

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity/template_common_entity"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity/template_content_entity"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity/template_layout_entity"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity/template_override_entity"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity/template_rule_entity"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity/template_variant_entity"
)

type (
	ValueOverride = template_override_entity.ValueOverride

	PlacementRule = template_common_entity.PlacementRule

	ContentCountLimit    = template_content_entity.ContentCountLimit
	ContentLimit         = template_content_entity.ContentLimit
	ContentList          = template_content_entity.ContentList
	ContentPool          = template_content_entity.ContentPool
	MandatoryContent     = template_content_entity.MandatoryContent
	MandatoryContentItem = template_content_entity.MandatoryContentItem
	WeightedContent      = template_content_entity.WeightedContent

	Bonus         = template_rule_entity.Bonus
	BonusList     = template_rule_entity.BonusList
	GameRules     = template_rule_entity.GameRules
	GlobalBans    = template_rule_entity.GlobalBans
	WinConditions = template_rule_entity.WinConditions

	Border                 = template_variant_entity.Border
	Connection             = template_variant_entity.Connection
	EncounterHolesSettings = template_variant_entity.EncounterHolesSettings
	MainObject             = template_variant_entity.MainObject
	Noise                  = template_variant_entity.Noise
	Orientation            = template_variant_entity.Orientation
	Road                   = template_variant_entity.Road
	StringList             = template_variant_entity.StringList
	TypedRef               = template_variant_entity.TypedRef
	Variant                = template_variant_entity.Variant
	Zone                   = template_variant_entity.Zone

	AmbientPickupDistribution         = template_layout_entity.AmbientPickupDistribution
	ElevationMode                     = template_layout_entity.ElevationMode
	GuardedEncounterResourceFractions = template_layout_entity.GuardedEncounterResourceFractions
	ZoneLayoutDef                     = template_layout_entity.ZoneLayoutDef
)
