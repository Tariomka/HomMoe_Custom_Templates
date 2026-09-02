package dtos

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

type TemplateLoadDto struct {
	Template *template_model.Template
	Warnings []string
}
