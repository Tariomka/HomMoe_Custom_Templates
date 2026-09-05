package models

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
)

// ZoneEditorGeometry is the complete canvas geometry of the manual zone editor
// for one square canvas side: where every zone node sits, how the nodes should
// be drawn, how big they are, and the curve of every connection between them.
// Coordinates and radius are unrounded: the canvas rounds once, when it draws.
type ZoneEditorGeometry struct {
	Positions  map[string]Position
	Zones      []preview.Zone
	ZoneRadius float64
	Edges      []ZoneEditorEdge
}
