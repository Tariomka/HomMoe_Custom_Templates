package mappers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

type ITemplateMapper interface {
	ToModel(entity template_entity.RmgTemplate) template_model.Template
	ToEntity(model template_model.Template) template_entity.RmgTemplate
}
