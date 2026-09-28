package preview

import "github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"

type ConnectionCurve struct {
	ConnectionIndex int // Index of the connection this curve represents.
	Start           data.Vec2[float64]
	End             data.Vec2[float64]
	Control         data.Vec2[float64]
}
