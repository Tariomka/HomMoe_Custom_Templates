// Package entities re-exports the .rmg.json schema types under their historical
// short names. The canonical home is internal/entities/template/types.go; these
// are aliases of aliases, so both spellings name the identical type. This needs
// to be removed and internal/entities/template should be used instead.
package entities

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template"
)

type (
	RmgTemplate = template.RmgTemplate

	ValueOverride = template.ValueOverride

	PlacementRule = template.PlacementRule

	ContentCountLimit    = template.ContentCountLimit
	ContentLimit         = template.ContentLimit
	ContentList          = template.ContentList
	ContentPool          = template.ContentPool
	MandatoryContent     = template.MandatoryContent
	MandatoryContentItem = template.MandatoryContentItem
	WeightedContent      = template.WeightedContent

	Bonus         = template.Bonus
	BonusList     = template.BonusList
	GameRules     = template.GameRules
	GlobalBans    = template.GlobalBans
	WinConditions = template.WinConditions

	Border                 = template.Border
	Connection             = template.Connection
	EncounterHolesSettings = template.EncounterHolesSettings
	MainObject             = template.MainObject
	Noise                  = template.Noise
	Orientation            = template.Orientation
	Road                   = template.Road
	StringList             = template.StringList
	TypedRef               = template.TypedRef
	Variant                = template.Variant
	Zone                   = template.Zone

	AmbientPickupDistribution         = template.AmbientPickupDistribution
	ElevationMode                     = template.ElevationMode
	GuardedEncounterResourceFractions = template.GuardedEncounterResourceFractions
	ZoneLayoutDef                     = template.ZoneLayoutDef
)
