package provider_interfaces

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

type IContentLimitProvider interface {
	CreateContentCountLimits(settings config.GeneratorConfig) []template_model.ContentCountLimit
}
