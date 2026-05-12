package template

// RmgTemplateModel represents the top-level template structure for .rmg.json files.
// Mirrors the on-disk schema observed across all bundled `data/ExampleTemplates/*.rmg.json` files.
type RmgTemplateModel struct {
	Name string `json:"name"`

	GameMode            string `json:"gameMode"`
	Description         string `json:"description"`
	DisplayWinCondition string `json:"displayWinCondition"`

	SizeX int `json:"sizeX"`
	SizeZ int `json:"sizeZ"`

	ValueOverrides []ValueOverride `json:"valueOverrides,omitempty"`

	// A handful of templates ("OctoJebus") declare a stray top-level `orientation` /
	// `border` block alongside `variants`. These hold the same data as Variant's own
	// orientation/border and are preserved here for round-trip fidelity.
	Orientation *Orientation `json:"orientation,omitempty"`
	Border      *Border      `json:"border,omitempty"`

	GameRules  GameRules   `json:"gameRules"`
	GlobalBans *GlobalBans `json:"globalBans,omitempty"`

	Variants []Variant `json:"variants"`

	ZoneLayouts        []ZoneLayoutDef     `json:"zoneLayouts,omitempty"`
	MandatoryContent   []MandatoryContent  `json:"mandatoryContent,omitempty"`
	ContentCountLimits []ContentCountLimit `json:"contentCountLimits,omitempty"`
	ContentPools       []ContentPool       `json:"contentPools"`
	ContentLists       []ContentList       `json:"contentLists"`
}
