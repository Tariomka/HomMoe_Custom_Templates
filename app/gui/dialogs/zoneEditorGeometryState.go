package dialogs

import (
	"image"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

// zoneEditorGeometryState caches the canvas layout the geometry service
// produced, together with the square canvas side it was built for.
type zoneEditorGeometryState struct {
	geometry models.ZoneEditorGeometry
	side     int
	// canvasOrigin is where the centred square starts inside the space the
	// canvas was given, which is the transform every square-local coordinate is
	// expressed against.
	canvasOrigin  image.Point
	geometryDirty bool
}
