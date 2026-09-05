package dtos

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

type ZoneEditorMutationDto struct {
	Zones       []template_model.Zone
	Connections []template_model.Connection
}
