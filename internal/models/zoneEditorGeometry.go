package models

import (
	"image"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
)

// ZoneEditorGeometry is the complete canvas geometry of the manual zone editor
// for one square canvas side: where every zone node sits, how the nodes should
// be drawn, how big they are, and the curve of every connection between them.
type ZoneEditorGeometry struct {
	Positions  map[string]image.Point
	Zones      []preview.Zone
	ZoneRadius int
	Edges      []ZoneEditorEdge
}
