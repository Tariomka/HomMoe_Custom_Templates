package editor_state

import "github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"

type ManualZoneEdit struct {
	Name    string `json:"name"`
	Quality *int8  `json:"quality,omitempty"`

	GeneratorPosition *data.Vec2[float64] `json:"generatorPosition,omitempty"`
	GeneratorRing     *int                `json:"generatorRing,omitempty"`
	ManualPosition    *data.Vec2[float64] `json:"manualPosition,omitempty"`

	Origin  ZoneOrigin `json:"origin,omitempty"`
	Castles *int       `json:"castles,omitempty"`

	FieldOverrides ZoneFieldOverrides `json:"fieldOverrides,omitzero"`
	Unset          []string           `json:"unset,omitempty"`
}
