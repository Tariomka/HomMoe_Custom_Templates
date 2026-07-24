package dtos

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
)

type PreviewLayoutRequestDto struct {
	Template   *entities.RmgTemplate
	Topology   config.MapTopology
	CanvasSide float64
}
