package template

// RmgTemplate represents the top-level template structure for .rmg.json files
type RmgTemplate struct {
	Name string `json:"name"`

	GameMode            string `json:"gameMode"`
	Description         string `json:"description"`
	DisplayWinCondition string `json:"displayWinCondition"`

	SizeX int `json:"sizeX"`
	SizeY int `json:"sizeY"`

	GameRules          GameRules           `json:"gameRules"`
	ValueOverrides     []ValueOverride     `json:"valueOverrides"`
	Variants           []Variant           `json:"variants"`
	ZoneLayouts        []ZoneLayout        `json:"zoneLayouts"`
	MandatoryContent   []ContentGroup      `json:"mandatoryContent"`
	ContentCountLimits []ContentCountLimit `json:"contentCountLimits"`
	ContentPools       []ContentPool       `json:"contentPools"`
	ContentLists       []ContentList       `json:"contentLists"`
}
