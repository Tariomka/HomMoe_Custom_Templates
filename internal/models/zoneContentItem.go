package models

// ZoneContentItem is the editable zone-content row used by the GUI
type ZoneContentItem struct {
	Sid          string `json:"sid,omitempty"`
	Name         string `json:"name,omitempty"`
	Count        int    `json:"count"`
	IsGuarded    bool   `json:"isGuarded"`
	NearCastle   bool   `json:"nearCastle,omitempty"`
	RoadDistance string `json:"roadDistance,omitempty"`
	IsGroup      bool   `json:"isGroup,omitempty"`
}
