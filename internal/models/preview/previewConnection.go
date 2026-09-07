package preview

import "github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"

// Connection is a drawn link between two zones on the preview canvas.
// Ctrl is the quadratic Bézier control point used to draw the edge: a lone
// edge keeps Ctrl on the midpoint (so it renders straight), while parallel
// edges between the same pair of zones spread their control points to either
// side so each connection stays individually visible.
type Connection struct {
	Start, Ctrl, End data.Vec2[float64]
	Type             ConnectionType
}

func (this Connection) IsPortal() bool {
	return this.Type == ConnectionTypePortal
}

func (this Connection) IsGladiatorArena() bool {
	return this.Type == ConnectionTypeGladiatorArena
}

type ConnectionType uint8

const (
	ConnectionTypeDirect ConnectionType = iota
	ConnectionTypePortal
	ConnectionTypeGladiatorArena
	ConnectionTypeProximity
)
