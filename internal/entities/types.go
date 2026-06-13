package entities

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template/template_common"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template/template_content"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template/template_layout"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template/template_override"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template/template_rule"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template/template_variant"
)

// Re-exported aliases so external callers can keep using template.Border, template.Variant, etc.
type (
	RmgTemplateModel = template.RmgTemplateModel

	ValueOverride = template_override.ValueOverride

	PlacementRule = template_common.PlacementRule

	ContentCountLimit    = template_content.ContentCountLimit
	ContentLimit         = template_content.ContentLimit
	ContentList          = template_content.ContentList
	ContentPool          = template_content.ContentPool
	MandatoryContent     = template_content.MandatoryContent
	MandatoryContentItem = template_content.MandatoryContentItem
	WeightedContent      = template_content.WeightedContent

	Bonus         = template_rule.Bonus
	BonusList     = template_rule.BonusList
	GameRules     = template_rule.GameRules
	GlobalBans    = template_rule.GlobalBans
	WinConditions = template_rule.WinConditions

	Border                 = template_variant.Border
	Connection             = template_variant.Connection
	EncounterHolesSettings = template_variant.EncounterHolesSettings
	MainObject             = template_variant.MainObject
	Noise                  = template_variant.Noise
	Orientation            = template_variant.Orientation
	Road                   = template_variant.Road
	StringList             = template_variant.StringList
	TypedRef               = template_variant.TypedRef
	Variant                = template_variant.Variant
	Zone                   = template_variant.Zone

	AmbientPickupDistribution         = template_layout.AmbientPickupDistribution
	ElevationMode                     = template_layout.ElevationMode
	GuardedEncounterResourceFractions = template_layout.GuardedEncounterResourceFractions
	ZoneLayoutDef                     = template_layout.ZoneLayoutDef
)
