package editor_state

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
)

type ManualZoneSave struct {
	Zone template_entity.Zone `json:"zone"`

	GeneratorPosition *data.Vec2[float64] `json:"generatorPosition,omitempty"`
	GeneratorRing     *int                `json:"generatorRing,omitempty"`

	ManualPosition *data.Vec2[float64] `json:"manualPosition,omitempty"`

	Quality *int8 `json:"quality,omitempty"`
}
