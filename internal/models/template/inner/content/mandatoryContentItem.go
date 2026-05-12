package content

// MandatoryContentItem describes a single mandatory object and how it should be placed.
// In some templates the item references a content list (via `includeLists`) and/or a
// nested weighted content roster (via `content`) instead of a single `sid`; all forms
// are supported.
type MandatoryContentItem struct {
	SID       string          `json:"sid,omitempty"`
	Name      string          `json:"name,omitempty"`
	IsMine    bool            `json:"isMine,omitempty"`
	IsGuarded bool            `json:"isGuarded,omitempty"`
	Rules     []PlacementRule `json:"rules,omitempty"`

	// Variant index this item applies to (-1 for all variants).
	Variant *int `json:"variant,omitempty"`

	// Optional starting owner of the placed object (e.g. "Player1").
	Owner string `json:"owner,omitempty"`

	// Override of the placed object's guard value.
	GuardValue int `json:"guardValue,omitempty"`

	// Alternative form: pull mandatory content from named content lists (Blitz, Arcade).
	IncludeLists []string `json:"includeLists,omitempty"`

	// Nested weighted content entries (Harmony) - typically a roster of random_item / random_hire SIDs.
	Content []WeightedContent `json:"content,omitempty"`

	// Whether the spawned object counts as a designated encounter (Anarchy, Maze).
	DesignatedEncounter *bool `json:"designatedEncounter,omitempty"`

	// Forces the item to be the only encounter in its host slot (Blitz, Crossroads).
	SoloEncounter bool `json:"soloEncounter,omitempty"`

	// Whether a road must connect to (or skip) the placed object (Christmas Tree, Hallway).
	Road *bool `json:"road,omitempty"`
}
