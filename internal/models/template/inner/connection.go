package inner

import "github.com/Tariomka/hommoe_custom_templates/internal/models/template/inner/content"

// Connection links two zones together inside a Variant.
type Connection struct {
	Name string `json:"name"`
	From string `json:"from"`
	To   string `json:"to"`

	ConnectionType string `json:"connectionType"`

	SimTurnSquad bool  `json:"simTurnSquad,omitempty"`
	Road         *bool `json:"road,omitempty"`

	GuardZone   string `json:"guardZone,omitempty"`
	GuardEscape bool   `json:"guardEscape,omitempty"`

	GuardValue           int     `json:"guardValue"`
	GuardRandomization   float64 `json:"guardRandomization,omitempty"`
	GuardWeeklyIncrement float64 `json:"guardWeeklyIncrement"`

	GatePlacement string `json:"gatePlacement,omitempty"`

	// Optional "Proximity" pseudo-connection length factor (Arcade, Hallway, etc.).
	Length float64 `json:"length,omitempty"`

	// Identifier used to share a guard pool across multiple connections (Jebus Cross family).
	GuardMatchGroup string `json:"guardMatchGroup,omitempty"`

	// Portal placement rules - present on connections of type "Portal".
	PortalPlacementRulesFrom []content.PlacementRule `json:"portalPlacementRulesFrom,omitempty"`
	PortalPlacementRulesTo   []content.PlacementRule `json:"portalPlacementRulesTo,omitempty"`
}
