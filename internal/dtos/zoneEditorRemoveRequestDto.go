package dtos

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

type ZoneEditorRemoveRequestDto struct {
	Zones       []template_model.Zone
	Connections []entities.Connection
	ZoneName    string
}
