package provider_interfaces

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
)

type IContentLimitProvider interface {
	CreateContentCountLimits(settings config.GeneratorConfig) []entities.ContentCountLimit
}
