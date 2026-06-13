package dtos

import "github.com/Tariomka/hommoe_custom_templates/internal/entities"

type TemplateDto struct {
	Template   *entities.RmgTemplate
	OutputPath string
}
