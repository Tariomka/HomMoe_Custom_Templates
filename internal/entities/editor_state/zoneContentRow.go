package editor_state

type ZoneContentRow struct {
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
	Rules []ContentRuleRow `json:"rules,omitempty"`
}
