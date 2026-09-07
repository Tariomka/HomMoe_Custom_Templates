package dtos

import "github.com/Tariomka/hommoe_custom_templates/internal/models"

// ZoneEditorHitTestRequestDto asks which zone node of a laid-out canvas covers
// Position, given the node centers and their shared radius.
type ZoneEditorHitTestRequestDto struct {
	Position   models.Position
	Positions  map[string]models.Position
	ZoneRadius float64
}
