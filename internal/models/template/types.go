package template

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template/inner"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template/inner/common"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template/inner/content"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template/inner/game_rules"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template/inner/zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template/inner/zone_layout"
)

// Re-exported aliases so external callers can keep using template.Border, template.Variant, etc.
type (
	ValueOverride = inner.ValueOverride

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

	Border                 = zone.Border
	Connection             = zone.Connection
	EncounterHolesSettings = zone.EncounterHolesSettings
	MainObject             = zone.MainObject
	Noise                  = zone.Noise
	Orientation            = zone.Orientation
	Road                   = zone.Road
	StringList             = zone.StringList
	TypedRef               = zone.TypedRef
	Variant                = zone.Variant
	Zone                   = zone.Zone

	AmbientPickupDistribution         = zone_layout.AmbientPickupDistribution
	ElevationMode                     = zone_layout.ElevationMode
	GuardedEncounterResourceFractions = zone_layout.GuardedEncounterResourceFractions
	ZoneLayoutDef                     = zone_layout.ZoneLayoutDef
)
