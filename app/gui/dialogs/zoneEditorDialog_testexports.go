//go:build integration_test

package dialogs

import (
	"image"

	"gioui.org/f32"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
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
	MidPoint     image.Point
}

// RecomputeGeometry rebuilds node positions and edge curves for a square canvas
// of the given side, exactly as a laid-out frame would. ONLY FOR INTEGRATION TEST USE
func (this *ZoneEditorDialog) RecomputeGeometry(side int) { this.recomputeGeometry(side) }

// ZonePositions ONLY FOR INTEGRATION TEST USE
func (this *ZoneEditorDialog) ZonePositions() map[string]image.Point { return this.positions }

// CanvasZoneRadius ONLY FOR INTEGRATION TEST USE
func (this *ZoneEditorDialog) CanvasZoneRadius() int { return this.radius }

// EdgeGeometries returns the laid-out connection curves in draw order.
// ONLY FOR INTEGRATION TEST USE
func (this *ZoneEditorDialog) EdgeGeometries() []EdgeGeometry {
	edges := make([]EdgeGeometry, 0, len(this.edges))
	for _, edge := range this.edges {
		connection := this.edgeConnection(edge)
		if connection == nil {
			continue
		}
		edges = append(edges, EdgeGeometry{
			Name:         connection.Name,
			From:         connection.From,
			To:           connection.To,
			StartPoint:   toCanvasPoint(edge.StartPoint),
			EndPoint:     toCanvasPoint(edge.EndPoint),
			ControlPoint: toCanvasPoint(edge.ControlPoint),
			MidPoint:     edge.MidPoint,
		})
	}

	return edges
}

// HitTestCanvasNode returns the zone whose node covers pos, or "" when none
// does. ONLY FOR INTEGRATION TEST USE
func (this *ZoneEditorDialog) HitTestCanvasNode(pos image.Point) string {
	return this.hitTestNode(pos)
}

// HitTestCanvasEdge returns the name of the connection whose curve passes
// closest to pos, or "" when no curve is within reach.
// ONLY FOR INTEGRATION TEST USE
func (this *ZoneEditorDialog) HitTestCanvasEdge(pos image.Point) string {
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
func (this *ZoneEditorDialog) SnapDraggedPosition(pos image.Point) image.Point {
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
			this.selected = connection
			this.selectedZone = ""
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

// ClickReset ONLY FOR INTEGRATION TEST USE
func (this *ZoneEditorDialog) ClickReset() { this.resetBtn.Click() }

// ClickApply ONLY FOR INTEGRATION TEST USE
func (this *ZoneEditorDialog) ClickApply() { this.applyBtn.Click() }

// ClickAddConnection ONLY FOR INTEGRATION TEST USE
func (this *ZoneEditorDialog) ClickAddConnection() { this.addBtn.Click() }

// ClickAddZone ONLY FOR INTEGRATION TEST USE
func (this *ZoneEditorDialog) ClickAddZone() { this.addZoneBtn.Click() }

// ClickDeleteSelected ONLY FOR INTEGRATION TEST USE
func (this *ZoneEditorDialog) ClickDeleteSelected() { this.deleteBtn.Click() }

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
