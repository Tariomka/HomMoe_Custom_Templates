package mappers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

type ITemplateMapper interface {
	ToModel(entity template.RmgTemplate) template_model.Template
	ToEntity(model template_model.Template) template.RmgTemplate
}
