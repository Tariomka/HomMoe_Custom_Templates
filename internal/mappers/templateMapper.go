package mappers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

// TemplateMapper converts between the .rmg.json wire format and the template
// the service layer works with. The entity is required at exactly two places -
// reading and writing the file - so every other seam should hold the model and
// reach for this only at the boundary.
type TemplateMapper struct{}

func NewTemplateMapper() ITemplateMapper {
	return &TemplateMapper{}
}

// ToModel lifts a template read from disk. Nothing on the wire records a zone's
// tier, so every zone comes back with a nil Quality, meaning "infer it".
func (this *TemplateMapper) ToModel(entity template.RmgTemplate) template_model.Template {
	return template_model.ToTemplateModel(entity)
}

// ToEntity flattens the template back to the wire format, dropping the tier the
// schema has no field for.
func (this *TemplateMapper) ToEntity(model template_model.Template) template.RmgTemplate {
	return template_model.ToTemplateEntity(model)
}
