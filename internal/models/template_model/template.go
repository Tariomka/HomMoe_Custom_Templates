package template_model

import (
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
)

type Template struct {
	Name    string
	MapSize int

	GameMode            string
	Description         string
	DisplayWinCondition string

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
