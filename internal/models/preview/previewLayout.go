package preview

import "github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"

// Layout is the full geometry of a preview rendered into a square canvas of the requested side length.
// Coordinates and radius are unrounded: callers round once, at their draw call.
type Layout struct {
	Positions   map[string]data.Vec2[float64]
	Zones       []Zone
	Connections []Connection
	ZoneRadius  float64
}
