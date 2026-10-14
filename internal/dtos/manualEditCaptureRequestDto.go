package dtos

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

// ManualEditCaptureRequestDto carries the settings to generate the derivation
// from and the edited layout to record against it.
type ManualEditCaptureRequestDto struct {
	State       editor_state_dto.EditorStateDto
	Zones       []template_model.Zone
	Connections []template_model.Connection
}
