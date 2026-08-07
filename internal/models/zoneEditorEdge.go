package models

import (
	"image"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
)

// ZoneEditorEdge is one laid-out connection curve of the manual zone editor: a
// quadratic Bezier running from StartPoint to EndPoint, bent by ControlPoint,
// with MidPoint marking where the edge's label belongs. ConnectionIndex points
// back into the connection slice the geometry was built from.
type ZoneEditorEdge struct {
	ConnectionIndex int
	StartPoint      data.Vec2[float64]
	EndPoint        data.Vec2[float64]
	ControlPoint    data.Vec2[float64]
	MidPoint        image.Point
}
