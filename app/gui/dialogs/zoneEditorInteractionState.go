package dialogs

import (
	"image"
	"math"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
)

// zoneDragDeadZonePx is how far the pointer must travel before a press on a
// zone becomes a move rather than a plain selection click.
const zoneDragDeadZonePx = 6.0

// zoneEditorInteractionState is the canvas view state: what is selected, which
// add mode is armed, and what the pointer is currently dragging. It is view
// state only - nothing here belongs to the domain or to the geometry service.
type zoneEditorInteractionState struct {
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

// selectConnection makes an edge the active selection.
func (this *zoneEditorInteractionState) selectConnection(connection *entities.Connection) {
	this.selected = connection
	this.selectedZone = ""
}

// selectZoneNamed makes a zone the active selection.
func (this *zoneEditorInteractionState) selectZoneNamed(name string) {
	this.selectedZone = name
	this.selected = nil
	this.hint = ""
}

// clearSelection drops both the edge and the zone selection.
func (this *zoneEditorInteractionState) clearSelection() {
	this.selected = nil
	this.selectedZone = ""
}

func (this *zoneEditorInteractionState) hasSelection() bool {
	return this.selected != nil || this.selectedZone != ""
}

// toggleAddConnectionMode arms drag-to-connect and disarms every other mode.
func (this *zoneEditorInteractionState) toggleAddConnectionMode() {
	armed := !this.addMode
	this.exitAddModes()
	this.hint = ""
	this.addMode = armed
}

// toggleAddZoneMode arms click-to-place and disarms every other mode.
func (this *zoneEditorInteractionState) toggleAddZoneMode() {
	armed := !this.addZoneMode
	this.exitAddModes()
	this.hint = ""
	this.addZoneMode = armed
}

// exitAddModes disarms both add modes and abandons a half-finished connection
// drag, leaving the status hint alone.
func (this *zoneEditorInteractionState) exitAddModes() {
	this.addMode = false
	this.addZoneMode = false
	this.pendingFrom = ""
	this.dragging = false
}

// beginConnectionDrag starts a rubber band from the given zone.
func (this *zoneEditorInteractionState) beginConnectionDrag(from string, pos image.Point) {
	this.pendingFrom = from
	this.dragging = true
	this.dragPos = pos
}

// finishConnectionDrag ends the rubber band and reports the zone it started
// from, or "" when no drag was in progress.
func (this *zoneEditorInteractionState) finishConnectionDrag() string {
	if !this.dragging {
		return ""
	}
	from := this.pendingFrom
	this.dragging = false
	this.pendingFrom = ""

	return from
}

// beginZoneDrag records a press on a zone; the move only starts once the
// pointer leaves the dead zone.
func (this *zoneEditorInteractionState) beginZoneDrag(name string, pos image.Point) {
	this.zoneDragName = name
	this.zoneDragMoved = false
	this.pressPos = pos
}

func (this *zoneEditorInteractionState) endZoneDrag() {
	this.zoneDragName = ""
	this.zoneDragMoved = false
}

// zoneDragLeftDeadZone reports whether the pointer travelled far enough from
// the press point to turn the press into a move.
func (this *zoneEditorInteractionState) zoneDragLeftDeadZone(pos image.Point) bool {
	deltaX := float64(pos.X - this.pressPos.X)
	deltaY := float64(pos.Y - this.pressPos.Y)

	return math.Hypot(deltaX, deltaY) >= zoneDragDeadZonePx
}

// reset returns the canvas to a neutral state: nothing selected, no mode
// armed, no drag in progress.
func (this *zoneEditorInteractionState) reset() {
	this.clearSelection()
	this.exitAddModes()
	this.endZoneDrag()
	this.hint = ""
}
