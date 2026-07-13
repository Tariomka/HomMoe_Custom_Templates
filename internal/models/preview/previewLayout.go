package preview

import "image"

// Layout is the full geometry of a preview rendered into a square canvas of the requested side length.
type Layout struct {
	Positions   map[string]image.Point
	Zones       []Zone
	Connections []Connection
	ZoneRadius  int
}
