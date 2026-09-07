package dtos

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

type TemplateSaveDto struct {
	Template   *template_model.Template
	Topology   config.MapTopology
	OutputPath string
}
