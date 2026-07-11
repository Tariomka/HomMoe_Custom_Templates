package dtos

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
)

type TemplateSaveDto struct {
	Template   *entities.RmgTemplate
	Topology   config.MapTopology
	OutputPath string
}
