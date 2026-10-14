package dtos

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

type ManualEditExpansionDto struct {
	Zones       []template_model.Zone
	Connections []template_model.Connection
	ManualEdits editor_state_model.ManualEditSettings
	Changed     bool
	Warnings    []string
}
