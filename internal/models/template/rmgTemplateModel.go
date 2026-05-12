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

// ValueOverride overrides the default guard value of a specific object SID, optionally per variant index.
type ValueOverride struct {
	SID        string `json:"sid"`
	Variant    int    `json:"variant"`
	GuardValue int    `json:"guardValue"`
}

// GlobalBans declares globally banned content (items, magics, heroes) at the template level.
type GlobalBans struct {
	Items  []string `json:"items,omitempty"`
	Magics []string `json:"magics,omitempty"`
	Heroes []string `json:"heroes,omitempty"`
}
