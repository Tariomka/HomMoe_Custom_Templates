package editor_state_v1

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
)

type ManualZoneSave struct {
	Zone           template_entity.Zone `json:"zone"`
	ManualPosition *[2]float64          `json:"manualPosition,omitempty"`
	Quality        *int8                `json:"quality,omitempty"`
}
