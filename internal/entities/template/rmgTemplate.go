package template

// RmgTemplate represents the top-level template structure for .rmg.json files.
// Mirrors the on-disk schema observed across all bundled `data/ExampleTemplates/*.rmg.json` files.
type RmgTemplate struct {
	Name string `json:"name"`

	GameMode            string `json:"gameMode"`
	Description         string `json:"description"`
	DisplayWinCondition string `json:"displayWinCondition"`

	SizeX int `json:"sizeX"`
	SizeZ int `json:"sizeZ"`

	ValueOverrides []ValueOverride `json:"valueOverrides,omitempty"`

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
