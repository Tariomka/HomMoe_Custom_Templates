package template

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template/inner"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template/inner/content"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template/inner/gamerules"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template/inner/zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template/inner/zonelayout"
)

// RmgTemplateModel represents the top-level template structure for .rmg.json files.
// Mirrors the on-disk schema observed across all bundled `data/ExampleTemplates/*.rmg.json` files.
type RmgTemplateModel struct {
	Name string `json:"name"`

	GameMode            string `json:"gameMode"`
	Description         string `json:"description"`
	DisplayWinCondition string `json:"displayWinCondition"`

	SizeX int `json:"sizeX"`
	SizeZ int `json:"sizeZ"`

	ValueOverrides []inner.ValueOverride `json:"valueOverrides,omitempty"`

	// A handful of templates ("OctoJebus") declare a stray top-level `orientation` /
	// `border` block alongside `variants`. These hold the same data as Variant's own
	// orientation/border and are preserved here for round-trip fidelity.
	Orientation *inner.Orientation `json:"orientation,omitempty"`
	Border      *inner.Border      `json:"border,omitempty"`

	GameRules  gamerules.GameRules `json:"gameRules"`
	GlobalBans *inner.GlobalBans   `json:"globalBans,omitempty"`

	Variants []zone.Variant `json:"variants"`

	ZoneLayouts        []zonelayout.ZoneLayoutDef  `json:"zoneLayouts,omitempty"`
	MandatoryContent   []content.MandatoryContent  `json:"mandatoryContent,omitempty"`
	ContentCountLimits []content.ContentCountLimit `json:"contentCountLimits,omitempty"`
	ContentPools       []content.ContentPool       `json:"contentPools"`
	ContentLists       []content.ContentList       `json:"contentLists"`
}
