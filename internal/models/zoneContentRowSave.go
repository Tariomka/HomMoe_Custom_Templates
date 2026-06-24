package models

// ZoneContentRowSave is the lightweight serialisation record for a single
// mandatory-content UI row. It preserves the row exactly as the user
// configured it, including the Count slider - so e.g. two separate sawmill
// rows stay as two rows after a round-trip.
type ZoneContentRowSave struct {
	// SID of the content item or include-list.
	Sid string `json:"sid"`
	// Spinner / Count value shown in the UI row.
	Count int `json:"count"`
	// True when the SID is an include-list group rather than a concrete item.
	IsGroup bool `json:"isGroup"`

	// True when this row lives in the Mines collection (affects IsMine on
	// the generated MandatoryContentItem).
	IsMine bool `json:"isMine,omitempty"`

	// Rules is the serialized list of content rules for the row. New settings
	// files use this in place of the deprecated flat fields above.
	Rules []ContentRuleRowSave `json:"rules,omitempty"`
}

// Normalized returns a copy with the default values applied
// (Count >= 1, RoadDistance == "Any" when empty and no rules are present).
func (this ZoneContentRowSave) Normalized() ZoneContentRowSave {
	out := this
	if out.Count < 1 {
		out.Count = 1
	}
	return out
}
