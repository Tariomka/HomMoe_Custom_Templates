package models

// ZoneContentRowSave is the lightweight serialisation record for a single
// mandatory-content UI row. It preserves the row exactly as the user
// configured it, including the Count slider — so e.g. two separate sawmill
// rows stay as two rows after a round-trip.
type ZoneContentRowSave struct {
	// SID of the content item or include-list.
	Sid string `json:"sid"`
	// Spinner / Count value shown in the UI row.
	Count int `json:"count"`
	// True when the SID is an include-list group rather than a concrete item.
	IsGroup bool `json:"isGroup"`

	// Whether the content is guarded.
	//
	// Deprecated: left in for backward-compatibility with old saves
	// (0.7.1 and earlier). New settings files use the Rules list instead.
	IsGuarded bool `json:"isGuarded,omitempty"`
	// Whether the Near Castle placement rule is active.
	//
	// Deprecated: left in for backward-compatibility with old saves
	// (0.7.1 and earlier). New settings files use the Rules list instead.
	NearCastle bool `json:"nearCastle,omitempty"`
	// Road-distance label: "Any" | "Next To" | "Near" | "Medium" | "Far" | "Very Far".
	//
	// Deprecated: left in for backward-compatibility with old saves
	// (0.7.1 and earlier). New settings files use the Rules list instead.
	RoadDistance string `json:"roadDistance,omitempty"`

	// True when this row lives in the Mines collection (affects IsMine on
	// the generated MandatoryContentItem).
	IsMine bool `json:"isMine,omitempty"`

	// Rules is the serialized list of content rules for the row. New settings
	// files use this in place of the deprecated flat fields above.
	Rules []ContentRuleRowSave `json:"rules,omitempty"`
}

// Normalised returns a copy with the default values applied
// (Count >= 1, RoadDistance == "Any" when empty and no rules are present).
func (this ZoneContentRowSave) Normalised() ZoneContentRowSave {
	out := this
	if out.Count < 1 {
		out.Count = 1
	}
	// Only seed the legacy default when there is no new-style rule list, so
	// that migrated/new rows do not gain a phantom "Any" road-distance label.
	if out.RoadDistance == "" && len(out.Rules) == 0 {
		out.RoadDistance = "Any"
	}
	return out
}

// HasLegacyRuleData reports whether the row carries pre-rules flat fields that
// still need migrating into the Rules list.
func (this ZoneContentRowSave) HasLegacyRuleData() bool {
	return this.IsGuarded || this.NearCastle ||
		(this.RoadDistance != "" && this.RoadDistance != "Any")
}
