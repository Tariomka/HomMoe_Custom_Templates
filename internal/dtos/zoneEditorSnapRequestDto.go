package dtos

import "github.com/Tariomka/hommoe_custom_templates/internal/models"

// ZoneEditorSnapRequestDto asks for the snapped center of the zone named
// DraggedZone. That zone is excluded from the alignment guides it snaps to, so
// a zone never holds onto itself.
type ZoneEditorSnapRequestDto struct {
	Position    models.Position
	Positions   map[string]models.Position
	ZoneRadius  float64
	DraggedZone string
}
