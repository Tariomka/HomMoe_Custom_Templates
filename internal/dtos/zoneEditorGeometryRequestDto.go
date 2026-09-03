package dtos

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

// ZoneEditorGeometryRequestDto asks for the manual zone editor's canvas
// geometry: the given zones and connections laid out into a square canvas of
// CanvasSide pixels per side.
type ZoneEditorGeometryRequestDto struct {
	Zones       []template_model.Zone
	Connections []entities.Connection
	Topology    config.MapTopology
	CanvasSide  int
}
