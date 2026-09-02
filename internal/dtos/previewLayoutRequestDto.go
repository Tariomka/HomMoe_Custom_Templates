package dtos

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

type PreviewLayoutRequestDto struct {
	// Template takes precedence; Zones and Connections provide an editor-only preview when it is nil.
	Template    *template_model.Template
	Zones       []entities.Zone
	Connections []entities.Connection
	Topology    config.MapTopology
	CanvasSide  float64
}
