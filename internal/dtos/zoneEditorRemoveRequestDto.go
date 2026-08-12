package dtos

import "github.com/Tariomka/hommoe_custom_templates/internal/entities"

type ZoneEditorRemoveRequestDto struct {
	Zones       []entities.Zone
	Connections []entities.Connection
	ZoneName    string
}
