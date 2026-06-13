package dtos

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config/config_inner"
)

type TemplateSaveDto struct {
	Template   *entities.RmgTemplate
	Topology   config_inner.MapTopology
	OutputPath string
}
