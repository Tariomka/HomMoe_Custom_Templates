package template

// MandatoryContent is a named group of objects that must be placed somewhere on the map.
type MandatoryContent struct {
	Name    string                 `json:"name"`
	Content []MandatoryContentItem `json:"content"`
}

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

// WeightedContent is a sid + weight pair used inside MandatoryContentItem.Content rosters.
type WeightedContent struct {
	SID    string `json:"sid"`
	Weight int    `json:"weight"`
}

// PlacementRule biases placement of a MandatoryContentItem (or portal endpoint) towards a reference point.
// Observed `type` values include: "Crossroads", "Road", "MainObject", "Sid".
// `Args` is loosely-typed (most rules use strings; "Pyramid" / "Sand Clover" pass numeric thresholds).
// A rule may specify either an inclusive [targetMin, targetMax] range OR a single `target` value.
type PlacementRule struct {
	Type      string  `json:"type"`
	Args      []any   `json:"args"`
	TargetMin float64 `json:"targetMin,omitempty"`
	TargetMax float64 `json:"targetMax,omitempty"`
	Target    float64 `json:"target,omitempty"`
	Weight    int     `json:"weight"`
}

// ContentCountLimit is a named cap on how many of certain SIDs may appear.
type ContentCountLimit struct {
	Name   string         `json:"name"`
	Limits []ContentLimit `json:"limits"`
}

// ContentLimit caps a single object SID at MaxCount instances, optionally per variant index.
// A few templates (e.g. "Blitz") allow caps to apply to whole content lists via `includeLists`
// in place of a single `sid`.
type ContentLimit struct {
	SID          string            `json:"sid,omitempty"`
	IncludeLists []string          `json:"includeLists,omitempty"`
	Content      []WeightedContent `json:"content,omitempty"`
	Variant      *int              `json:"variant,omitempty"`
	MaxCount     int               `json:"maxCount"`
}

// ContentPool is a placeholder for the `contentPools` array entries.
// The bundled example templates always declare this as an empty array; the inner
// shape is preserved as a free-form JSON object so unknown future entries decode.
type ContentPool map[string]any

// ContentList is a placeholder for the `contentLists` array entries.
// Like ContentPool, every bundled example declares this empty; kept as a free-form
// JSON object for forward compatibility.
type ContentList map[string]any
