package template

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template/inner"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template/inner/content"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template/inner/gamerules"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template/inner/zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template/inner/zonelayout"
)

// Re-exported aliases so external callers can keep using template.Border, template.Variant, etc.
type (
	Border        = inner.Border
	Connection    = inner.Connection
	GlobalBans    = inner.GlobalBans
	MainObject    = inner.MainObject
	Noise         = inner.Noise
	Orientation   = inner.Orientation
	Road          = inner.Road
	StringList    = inner.StringList
	TypedRef      = inner.TypedRef
	ValueOverride = inner.ValueOverride

	ContentCountLimit    = content.ContentCountLimit
	ContentLimit         = content.ContentLimit
	ContentList          = content.ContentList
	ContentPool          = content.ContentPool
	MandatoryContent     = content.MandatoryContent
	MandatoryContentItem = content.MandatoryContentItem
	PlacementRule        = content.PlacementRule
	WeightedContent      = content.WeightedContent

	Bonus         = gamerules.Bonus
	BonusList     = gamerules.BonusList
	GameRules     = gamerules.GameRules
	WinConditions = gamerules.WinConditions

	EncounterHolesSettings = zone.EncounterHolesSettings
	Variant                = zone.Variant
	Zone                   = zone.Zone

	AmbientPickupDistribution         = zonelayout.AmbientPickupDistribution
	ElevationMode                     = zonelayout.ElevationMode
	GuardedEncounterResourceFractions = zonelayout.GuardedEncounterResourceFractions
	ZoneLayoutDef                     = zonelayout.ZoneLayoutDef
)
