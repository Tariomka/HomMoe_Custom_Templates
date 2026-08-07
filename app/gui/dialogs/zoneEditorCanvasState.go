package dialogs

import (
	"image"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
)

type zoneEditorCanvasState struct {
	positions     map[string]image.Point
	previewZones  []preview.Zone
	radius        int
	side          int
	edges         []models.ZoneEditorEdge
	geometryDirty bool
	geometrySide  int

	canvasTag     int
	selected      *entities.Connection
	selectedZone  string
	addMode       bool
	addZoneMode   bool
	pendingFrom   string
	dragging      bool
	dragPos       image.Point
	zoneDragName  string
	zoneDragMoved bool
	pressPos      image.Point
	hint          string
}
