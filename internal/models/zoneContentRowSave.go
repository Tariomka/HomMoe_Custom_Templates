package models

import "slices"

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

	// Rules is the serialized list of content rules for the row.
	Rules []ContentRuleRowSave `json:"rules,omitempty"`
}

// Normalized returns a copy with the default values applied.
func (this ZoneContentRowSave) Normalized() ZoneContentRowSave {
	out := this
	if out.Count < 1 {
		out.Count = 1
	}
	return out
}

// Clone returns a copy that shares no backing array or pointer with the
// receiver. A nil Rules slice stays nil, because the change detection in
// EditorStateModel distinguishes a nil slice from an empty one.
func (this ZoneContentRowSave) Clone() ZoneContentRowSave {
	clone := this
	clone.Rules = slices.Clone(this.Rules)
	for ruleIndex := range clone.Rules {
		clone.Rules[ruleIndex] = clone.Rules[ruleIndex].Clone()
	}
	return clone
}
