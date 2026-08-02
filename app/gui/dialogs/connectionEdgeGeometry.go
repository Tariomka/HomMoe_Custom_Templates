package dialogs

import (
	"image"

	"gioui.org/f32"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
)

type connectionEdgeGeometry struct {
	connection   *entities.Connection
	startPoint   f32.Point
	endPoint     f32.Point
	controlPoint f32.Point
	midPoint     image.Point
}
