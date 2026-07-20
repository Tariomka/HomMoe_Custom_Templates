package dtos

import "github.com/Tariomka/hommoe_custom_templates/internal/entities"

type TemplateLoadDto struct {
	Template *entities.RmgTemplate
	Warnings []string
}
