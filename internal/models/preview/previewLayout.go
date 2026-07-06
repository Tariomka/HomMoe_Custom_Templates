package preview

import "image"

// PreviewLayout is the full geometry of a preview rendered into a square
// canvas of the requested side length.
type PreviewLayout struct {
	Positions   map[string]image.Point // TODO: use Vec2 instead of Point for sub-pixel precision
	Zones       []PreviewZone
	Connections []PreviewConnection
	ZoneRadius  int
}
