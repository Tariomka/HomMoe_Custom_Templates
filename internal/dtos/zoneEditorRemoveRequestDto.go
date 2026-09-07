package dtos

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

type ZoneEditorRemoveRequestDto struct {
	Zones       []template_model.Zone
	Connections []template_model.Connection
	ZoneName    string
}
