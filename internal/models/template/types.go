package template

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template/template_inner"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template/template_inner/common"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template/template_inner/content"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template/template_inner/game_rules"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template/template_inner/variant"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template/template_inner/zone_layout"
)

// Re-exported aliases so external callers can keep using template.Border, template.Variant, etc.
type (
	ValueOverride = template_inner.ValueOverride

	PlacementRule = common.PlacementRule

	ContentCountLimit    = content.ContentCountLimit
	ContentLimit         = content.ContentLimit
	ContentList          = content.ContentList
	ContentPool          = content.ContentPool
	MandatoryContent     = content.MandatoryContent
	MandatoryContentItem = content.MandatoryContentItem
	WeightedContent      = content.WeightedContent

	Bonus         = game_rules.Bonus
	BonusList     = game_rules.BonusList
	GameRules     = game_rules.GameRules
	GlobalBans    = game_rules.GlobalBans
	WinConditions = game_rules.WinConditions

	Border                 = variant.Border
	Connection             = variant.Connection
	EncounterHolesSettings = variant.EncounterHolesSettings
	MainObject             = variant.MainObject
	Noise                  = variant.Noise
	Orientation            = variant.Orientation
	Road                   = variant.Road
	StringList             = variant.StringList
	TypedRef               = variant.TypedRef
	Variant                = variant.Variant
	Zone                   = variant.Zone

	AmbientPickupDistribution         = zone_layout.AmbientPickupDistribution
	ElevationMode                     = zone_layout.ElevationMode
	GuardedEncounterResourceFractions = zone_layout.GuardedEncounterResourceFractions
	ZoneLayoutDef                     = zone_layout.ZoneLayoutDef
)
