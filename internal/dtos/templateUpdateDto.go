package dtos

import "github.com/Tariomka/hommoe_custom_templates/internal/entities"

type TemplateUpdateDto struct {
	Template    *entities.RmgTemplate
	Zones       []entities.Zone
	Connections []entities.Connection
}
