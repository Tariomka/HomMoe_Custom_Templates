package preview

import "image"

// Connection is a drawn link between two zones on the preview canvas.
// Ctrl is the quadratic Bézier control point used to draw the edge: a lone
// edge keeps Ctrl on the midpoint (so it renders straight), while parallel
// edges between the same pair of zones spread their control points to either
// side so each connection stays individually visible.
type Connection struct {
	Start, Ctrl, End image.Point
	Portal           bool
}
