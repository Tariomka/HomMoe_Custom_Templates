package template_model

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model/template_common_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model/template_content_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model/template_layout_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model/template_override_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model/template_rule_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model/template_variant_model"
)

type (
	ValueOverride = template_override_model.ValueOverride

	PlacementRule = template_common_model.PlacementRule

	ContentCountLimit    = template_content_model.ContentCountLimit
	ContentLimit         = template_content_model.ContentLimit
	ContentList          = template_content_model.ContentList
	ContentPool          = template_content_model.ContentPool
	MandatoryContent     = template_content_model.MandatoryContent
	MandatoryContentItem = template_content_model.MandatoryContentItem
	WeightedContent      = template_content_model.WeightedContent

	Bonus         = template_rule_model.Bonus
	BonusList     = template_rule_model.BonusList
	GameRules     = template_rule_model.GameRules
	GlobalBans    = template_rule_model.GlobalBans
	WinConditions = template_rule_model.WinConditions

	Border                 = template_variant_model.Border
	Connection             = template_variant_model.Connection
	EncounterHolesSettings = template_variant_model.EncounterHolesSettings
	MainObject             = template_variant_model.MainObject
	Noise                  = template_variant_model.Noise
	Orientation            = template_variant_model.Orientation
	Road                   = template_variant_model.Road
	StringList             = template_variant_model.StringList
	TypedRef               = template_variant_model.TypedRef
	Variant                = template_variant_model.Variant
	Zone                   = template_variant_model.Zone

	AmbientPickupDistribution         = template_layout_model.AmbientPickupDistribution
	ElevationMode                     = template_layout_model.ElevationMode
	GuardedEncounterResourceFractions = template_layout_model.GuardedEncounterResourceFractions
	ZoneLayoutDef                     = template_layout_model.ZoneLayoutDef
)
