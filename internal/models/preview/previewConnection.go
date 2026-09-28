package preview

import "github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"

// Connection is a drawn link between two zones on the preview canvas.
// Ctrl is the quadratic Bézier control point.
type Connection struct {
	Start, Ctrl, End data.Vec2[float64]
	Type             ConnectionType
	HasRoad          bool
	ExplicitPortal   bool
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
