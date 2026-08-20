//go:build integration_test

package dialogs

import (
	"image"

	"gioui.org/f32"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/utils"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

// EdgeGeometry is a test-visible copy of one laid-out connection curve.
// ONLY FOR INTEGRATION TEST USE
type EdgeGeometry struct {
	Name         string
	From         string
	To           string
	StartPoint   f32.Point
	EndPoint     f32.Point
	ControlPoint f32.Point
	MidPoint     models.Position
}

// CanvasOrigin returns where the centred canvas square starts inside the space
// the canvas was laid out in, so a window-coordinate caller can convert to the
// square-local coordinates every other accessor here speaks.
// ONLY FOR INTEGRATION TEST USE
func (this *ZoneEditorDialog) CanvasOrigin() image.Point { return this.canvasOrigin }

// CanvasSquareSide ONLY FOR INTEGRATION TEST USE
func (this *ZoneEditorDialog) CanvasSquareSide() int { return this.side }

// RecomputeGeometry rebuilds node positions and edge curves for a square canvas
// of the given side, exactly as a laid-out frame would. ONLY FOR INTEGRATION TEST USE
func (this *ZoneEditorDialog) RecomputeGeometry(side int) { this.recomputeGeometry(side) }

// ZonePositions ONLY FOR INTEGRATION TEST USE
func (this *ZoneEditorDialog) ZonePositions() map[string]models.Position {
	return this.geometry.Positions
}

// CanvasZoneRadius ONLY FOR INTEGRATION TEST USE
func (this *ZoneEditorDialog) CanvasZoneRadius() float64 { return this.geometry.ZoneRadius }

// EdgeGeometries returns the laid-out connection curves in draw order.
// ONLY FOR INTEGRATION TEST USE
func (this *ZoneEditorDialog) EdgeGeometries() []EdgeGeometry {
	edges := make([]EdgeGeometry, 0, len(this.geometry.Edges))
	for _, edge := range this.geometry.Edges {
		connection := this.edgeConnection(edge)
		if connection == nil {
			continue
		}

		edges = append(edges, EdgeGeometry{
			Name:         connection.Name,
			From:         connection.From,
			To:           connection.To,
			StartPoint:   utils.ToF32Point(edge.StartPoint),
			EndPoint:     utils.ToF32Point(edge.EndPoint),
			ControlPoint: utils.ToF32Point(edge.ControlPoint),
			MidPoint:     edge.MidPoint,
		})
	}

	return edges
}

// HitTestCanvasNode returns the zone whose node covers pos, or "" when none
// does. ONLY FOR INTEGRATION TEST USE
func (this *ZoneEditorDialog) HitTestCanvasNode(pos models.Position) string {
	return this.hitTestNode(pos)
}

// HitTestCanvasEdge returns the name of the connection whose curve passes
// closest to pos, or "" when no curve is within reach.
// ONLY FOR INTEGRATION TEST USE
func (this *ZoneEditorDialog) HitTestCanvasEdge(pos models.Position) string {
	if edge := this.hitTestEdge(pos); edge != nil {
		return edge.Name
	}

	return ""
}

// CanvasGridStep ONLY FOR INTEGRATION TEST USE
func (this *ZoneEditorDialog) CanvasGridStep() float64 { return this.gridStep() }

// SetSnapEnabled ONLY FOR INTEGRATION TEST USE
func (this *ZoneEditorDialog) SetSnapEnabled(value bool) { this.snapBool.Value = value }

// BeginZoneDrag marks a zone as the one currently being dragged, which excludes
// it from the snap guides and enables the guide overlay.
// ONLY FOR INTEGRATION TEST USE
func (this *ZoneEditorDialog) BeginZoneDrag(name string) {
	this.zoneDragName = name
	this.zoneDragMoved = true
}

// SnapDraggedPosition ONLY FOR INTEGRATION TEST USE
func (this *ZoneEditorDialog) SnapDraggedPosition(pos models.Position) models.Position {
	return this.snapDraggedPosition(pos)
}

// SnapGuides reports the alignment guides the last snap produced.
// ONLY FOR INTEGRATION TEST USE
func (this *ZoneEditorDialog) SnapGuides() (x float64, xActive bool, y float64, yActive bool) {
	return this.snapGuideX, this.snapGuideXActive, this.snapGuideY, this.snapGuideYActive
}

// SelectZone ONLY FOR INTEGRATION TEST USE
func (this *ZoneEditorDialog) SelectZone(name string) { this.selectZone(name) }

// SelectedZone ONLY FOR INTEGRATION TEST USE
func (this *ZoneEditorDialog) SelectedZone() string { return this.selectedZone }

// SelectConnection selects the connection with the given name and reports
// whether it exists. ONLY FOR INTEGRATION TEST USE
func (this *ZoneEditorDialog) SelectConnection(name string) bool {
	for _, connection := range this.working {
		if connection.Name == name {
			this.selectConnection(connection)
			this.syncedFor = nil
			return true
		}
	}

	return false
}

// SelectedConnection ONLY FOR INTEGRATION TEST USE
func (this *ZoneEditorDialog) SelectedConnection() string {
	if this.selected == nil {
		return ""
	}

	return this.selected.Name
}

// ClickUndo ONLY FOR INTEGRATION TEST USE
func (this *ZoneEditorDialog) ClickUndo() { this.undoBtn.Click() }

// ClickRevertToBase ONLY FOR INTEGRATION TEST USE
func (this *ZoneEditorDialog) ClickRevertToBase() { this.revertBaseBtn.Click() }

// ClickApply ONLY FOR INTEGRATION TEST USE
func (this *ZoneEditorDialog) ClickApply() { this.applyBtn.Click() }

// ClickAddConnection ONLY FOR INTEGRATION TEST USE
func (this *ZoneEditorDialog) ClickAddConnection() { this.addBtn.Click() }

// ClickAddZone ONLY FOR INTEGRATION TEST USE
func (this *ZoneEditorDialog) ClickAddZone() { this.addZoneBtn.Click() }

// ClickDeleteSelected ONLY FOR INTEGRATION TEST USE
func (this *ZoneEditorDialog) ClickDeleteSelected() { this.deleteBtn.Click() }

// AddConnectionModeActive ONLY FOR INTEGRATION TEST USE
func (this *ZoneEditorDialog) AddConnectionModeActive() bool { return this.addMode }

// AddZoneModeActive ONLY FOR INTEGRATION TEST USE
func (this *ZoneEditorDialog) AddZoneModeActive() bool { return this.addZoneMode }

// DraggingZone returns the zone the pointer is currently moving, or "" when no
// zone drag is in progress. ONLY FOR INTEGRATION TEST USE
func (this *ZoneEditorDialog) DraggingZone() string { return this.zoneDragName }

// PendingConnectionSource returns the zone an in-progress connection drag
// started from, or "" when no rubber band is being drawn.
// ONLY FOR INTEGRATION TEST USE
func (this *ZoneEditorDialog) PendingConnectionSource() string { return this.pendingFrom }

// SnapEnabled ONLY FOR INTEGRATION TEST USE
func (this *ZoneEditorDialog) SnapEnabled() bool { return this.snapBool.Value }

// EditedZones ONLY FOR INTEGRATION TEST USE
func (this *ZoneEditorDialog) EditedZones() []entities.Zone { return this.zones }

// EditedConnectionNames ONLY FOR INTEGRATION TEST USE
func (this *ZoneEditorDialog) EditedConnectionNames() []string {
	names := make([]string, 0, len(this.working))
	for _, connection := range this.working {
		names = append(names, connection.Name)
	}

	return names
}

// StatusHint ONLY FOR INTEGRATION TEST USE
func (this *ZoneEditorDialog) StatusHint() string { return this.hint }
