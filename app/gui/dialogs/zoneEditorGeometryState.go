package dialogs

import "github.com/Tariomka/hommoe_custom_templates/internal/models"

// zoneEditorGeometryState caches the canvas layout the geometry service
// produced, together with the square canvas side it was built for.
type zoneEditorGeometryState struct {
	geometry      models.ZoneEditorGeometry
	side          int
	geometryDirty bool
}
