package models

// ZoneContentRowSave is the lightweight serialisation record for a single
// mandatory-content UI row. It preserves the row exactly as the user
// configured it, including the Count slider — so e.g. two separate sawmill
// rows stay as two rows after a round-trip.
//
// Mirror of the C# Models/Generator/ZoneContentRowSave.cs record.
type ZoneContentRowSave struct {
	// SID of the content item or include-list.
	Sid string `json:"sid"`
	// Spinner / Count value shown in the UI row.
	Count int `json:"count"`
	// True when the SID is an include-list group rather than a concrete item.
	IsGroup bool `json:"isGroup"`
	// Whether the content is guarded.
	IsGuarded bool `json:"isGuarded"`
	// Whether the Near Castle placement rule is active.
	NearCastle bool `json:"nearCastle"`
	// Road-distance label: "Any" | "Next To" | "Near" | "Medium" | "Far" | "Very Far".
	RoadDistance string `json:"roadDistance"`
	// True when this row lives in the Mines collection (affects IsMine on
	// the generated MandatoryContentItem).
	IsMine bool `json:"isMine"`
}

// Normalised returns a copy with the default values applied
// (Count >= 1, RoadDistance == "Any" when empty).
func (this ZoneContentRowSave) Normalised() ZoneContentRowSave {
	out := this
	if out.Count < 1 {
		out.Count = 1
	}
	if out.RoadDistance == "" {
		out.RoadDistance = "Any"
	}
	return out
}
