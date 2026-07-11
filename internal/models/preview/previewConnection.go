package preview

import "image"

// PreviewConnection is a drawn link between two zones on the preview canvas.
// Ctrl is the quadratic Bézier control point used to draw the edge: a lone
// edge keeps Ctrl on the midpoint (so it renders straight), while parallel
// edges between the same pair of zones spread their control points to either
// side so each connection stays individually visible.
type PreviewConnection struct {
	A, B   image.Point
	Ctrl   image.Point
	Portal bool
}
