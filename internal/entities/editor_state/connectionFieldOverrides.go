package editor_state

import "github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"

type ConnectionFieldOverrides struct {
	SimTurnSquad *bool `json:"simTurnSquad,omitzero"`
	Road         *bool `json:"road,omitzero"`

	GuardZone   *string `json:"guardZone,omitzero"`
	GuardEscape *bool   `json:"guardEscape,omitzero"`

	GuardValue           *int     `json:"guardValue,omitzero"`
	GuardRandomization   *float64 `json:"guardRandomization,omitzero"`
	GuardWeeklyIncrement *float64 `json:"guardWeeklyIncrement,omitzero"`

	GatePlacement *string `json:"gatePlacement,omitzero"`

	Length *float64 `json:"length,omitzero"`

	GuardMatchGroup *string `json:"guardMatchGroup,omitzero"`

	PortalPlacementRulesFrom []template_entity.PlacementRule `json:"portalPlacementRulesFrom,omitzero"`
	PortalPlacementRulesTo   []template_entity.PlacementRule `json:"portalPlacementRulesTo,omitzero"`
}
