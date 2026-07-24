package dtos

import "github.com/Tariomka/hommoe_custom_templates/internal/entities"

type ZoneEditorMutationDto struct {
	Zones       []entities.Zone
	Connections []entities.Connection
}
