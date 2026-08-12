package dtos

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
)

type PreviewLayoutRequestDto struct {
	// Template takes precedence; Zones and Connections provide an editor-only preview when it is nil.
	Template    *entities.RmgTemplate
	Zones       []entities.Zone
	Connections []entities.Connection
	Topology    config.MapTopology
	CanvasSide  float64
}
