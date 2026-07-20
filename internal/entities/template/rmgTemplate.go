package template

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template/template_content"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template/template_layout"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template/template_override"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template/template_rule"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template/template_variant"
)

// RmgTemplate represents the top-level template structure for .rmg.json files.
// Mirrors the on-disk schema observed across all bundled `data/ExampleTemplates/*.rmg.json` files.
type RmgTemplate struct {
	Name string `json:"name"`

	GameMode            string `json:"gameMode"`
	Description         string `json:"description"`
	DisplayWinCondition string `json:"displayWinCondition"`

	SizeX int `json:"sizeX"`
	SizeZ int `json:"sizeZ"`

	ValueOverrides []template_override.ValueOverride `json:"valueOverrides,omitempty"`

	Orientation *template_variant.Orientation `json:"orientation,omitempty"`
	Border      *template_variant.Border      `json:"border,omitempty"`

	GameRules  template_rule.GameRules   `json:"gameRules"`
	GlobalBans *template_rule.GlobalBans `json:"globalBans,omitempty"`

	Variants []template_variant.Variant `json:"variants"`

	ZoneLayouts        []template_layout.ZoneLayoutDef      `json:"zoneLayouts,omitempty"`
	MandatoryContent   []template_content.MandatoryContent  `json:"mandatoryContent,omitempty"`
	ContentCountLimits []template_content.ContentCountLimit `json:"contentCountLimits,omitempty"`
	ContentPools       []template_content.ContentPool       `json:"contentPools"`
	ContentLists       []template_content.ContentList       `json:"contentLists"`
}
