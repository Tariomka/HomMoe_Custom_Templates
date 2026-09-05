//go:build integration_test

package integration_common

import (
	"image"

	"gioui.org/f32"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/editor"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/zone_helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

// The labels the zone editor's buttons are addressed by. They are looked up
// inside the dialog panel, so a label the tab behind the scrim also uses still
// resolves to the dialog's own button.
const (
	addConnectionButtonLabel    = "Add connection"
	addingConnectionButtonLabel = "Adding... (click empty to stop)"
	addZoneButtonLabel          = "Add zone"
	placingZoneButtonLabel      = "Placing... (click a zone to stop)"
	deleteSelectedButtonLabel   = "Delete selected"
	undoButtonLabel             = "Undo"
	revertToBaseButtonLabel     = "Revert to Base"
	applyChangesButtonLabel     = "Apply changes"
	zoneEditorCancelButtonLabel = "Cancel"
	deleteThisZoneButtonLabel   = "Delete this zone"
	deleteConnectionButtonLabel = "Delete this connection"
)

// ZoneEditorHandler drives the manual zone editor dialog. It holds the
// originating BaseHandler but deliberately does not embed it: an open dialog's
// scrim absorbs every background click, so promoting the tab clicks would be a
// lie.
//
// It snapshots every action, but only once the tab has been taken off the Random
// topology - the dialog draws the generated layout, and a freshly randomised one
// is not something a golden can hold.
type ZoneEditorHandler struct {
	base *BaseHandler
}

// newZoneEditorHandler wraps the dialog the caller just opened and snapshots the
// opened state.
func newZoneEditorHandler(base *BaseHandler) *ZoneEditorHandler {
	base.runner.tb.Helper()
	handler := &ZoneEditorHandler{base: base}
	handler.verifySnapshot()
	return handler
}

func (this *ZoneEditorHandler) IsOpen() bool {
	this.base.runner.tb.Helper()
	return this.base.runner.DialogsOpen()
}

// Dialog exposes what the open dialog currently holds: its canvas geometry, its
// selection, the drag in flight and the zones and connections it has edited.
func (this *ZoneEditorHandler) Dialog() editor.IZoneEditorDialog {
	this.base.runner.tb.Helper()
	dialog, ok := this.base.runner.TopZoneEditor()
	if !ok {
		this.base.runner.tb.Fatal("the top-most dialog is not a zone editor")
	}

	return dialog
}

// CanvasPoint maps a canvas-local position to the window point that presses it,
// so a test can drive and assert against the same coordinates the dialog
// reports. The canvas centres its square inside the space it was given and
// reports that offset, so only the origin of that space is a measured constant.
func (this *ZoneEditorHandler) CanvasPoint(position models.Position) f32.Point {
	this.base.runner.tb.Helper()
	origin := this.Dialog().CanvasOrigin()
	return f32.Pt(
		float32(zoneEditorCanvasBoxLeft+origin.X)+float32(position.X),
		float32(zoneEditorCanvasBoxTop+origin.Y)+float32(position.Y))
}

// ButtonLabels reports every button the open dialog offers, for the assertions
// that are about what the toolbar says rather than about pressing it.
func (this *ZoneEditorHandler) ButtonLabels() []string {
	this.base.runner.tb.Helper()
	return this.base.runner.ButtonLabelsIn(zoneEditorRect())
}

// ZonePosition reads where the canvas currently draws the named zone.
func (this *ZoneEditorHandler) ZonePosition(name string) models.Position {
	this.base.runner.tb.Helper()
	position, ok := this.Dialog().ZonePositions()[name]
	if !ok {
		this.base.runner.tb.Fatalf("the canvas is not showing a zone called %q", name)
	}

	return position
}

// ClickCanvasAt taps the canvas at a canvas-local position. Whether that selects
// a zone, selects a connection, places a zone or clears the selection is the
// canvas' decision, which is what the pointer tests are there to pin down.
func (this *ZoneEditorHandler) ClickCanvasAt(position models.Position) *ZoneEditorHandler {
	this.base.runner.tb.Helper()
	this.base.runner.ClickAt(this.CanvasPoint(position))
	this.base.runner.NextFrame()
	this.verifySnapshot()
	return this
}

// ClickZone taps the centre of the named zone.
func (this *ZoneEditorHandler) ClickZone(name string) *ZoneEditorHandler {
	this.base.runner.tb.Helper()
	return this.ClickCanvasAt(this.ZonePosition(name))
}

// ClickConnection taps the midpoint of the named connection's curve.
func (this *ZoneEditorHandler) ClickConnection(name string) *ZoneEditorHandler {
	this.base.runner.tb.Helper()
	return this.ClickCanvasAt(this.connectionMidPoint(name))
}

// RightClickEdge presses the secondary mouse button on the named connection's
// curve, which is how the canvas deletes a connection.
func (this *ZoneEditorHandler) RightClickEdge(name string) *ZoneEditorHandler {
	this.base.runner.tb.Helper()
	this.base.runner.RightClickAt(this.CanvasPoint(this.connectionMidPoint(name)))
	this.base.runner.NextFrame()
	this.verifySnapshot()
	return this
}

// DragZone presses the named zone and releases at a canvas-local position. The
// canvas ignores a drag shorter than its dead zone, so a caller after a move
// rather than a selection has to travel further than that.
func (this *ZoneEditorHandler) DragZone(name string, to models.Position) *ZoneEditorHandler {
	this.base.runner.tb.Helper()
	return this.dragBetween(this.ZonePosition(name), to)
}

// DragFromZoneTo presses the named zone and releases on the named target zone,
// which is how a connection is drawn while the add-connection mode is armed.
func (this *ZoneEditorHandler) DragFromZoneTo(from string, to string) *ZoneEditorHandler {
	this.base.runner.tb.Helper()
	return this.dragBetween(this.ZonePosition(from), this.ZonePosition(to))
}

// ClickAddConnection arms the mode that turns the next drag between two zones
// into a connection, or disarms it. Both buttons relabel themselves while armed,
// so the label to click depends on the mode that is currently running.
func (this *ZoneEditorHandler) ClickAddConnection() *ZoneEditorHandler {
	if this.Dialog().AddConnectionModeActive() {
		return this.clickDialogButton(addingConnectionButtonLabel)
	}

	return this.clickDialogButton(addConnectionButtonLabel)
}

// ClickAddZone arms the mode that turns the next click on empty canvas into a
// new zone, or disarms it.
func (this *ZoneEditorHandler) ClickAddZone() *ZoneEditorHandler {
	if this.Dialog().AddZoneModeActive() {
		return this.clickDialogButton(placingZoneButtonLabel)
	}

	return this.clickDialogButton(addZoneButtonLabel)
}

// ClickDeleteSelected removes whatever the canvas currently has selected.
func (this *ZoneEditorHandler) ClickDeleteSelected() *ZoneEditorHandler {
	return this.clickDialogButton(deleteSelectedButtonLabel)
}

// ClickUndo steps back one edit.
func (this *ZoneEditorHandler) ClickUndo() *ZoneEditorHandler {
	return this.clickDialogButton(undoButtonLabel)
}

// ClickRevertToBase throws the manual layout away and regenerates it.
func (this *ZoneEditorHandler) ClickRevertToBase() *ZoneEditorHandler {
	return this.clickDialogButton(revertToBaseButtonLabel)
}

// ToggleSnap flips the toolbar's Snap checkbox, which emits a semantic.CheckBox
// rather than a labelled button and so is the one toolbar control addressed by
// coordinate.
func (this *ZoneEditorHandler) ToggleSnap() *ZoneEditorHandler {
	this.base.runner.tb.Helper()
	this.base.runner.ClickAt(f32.Pt(zoneEditorSnapCheckboxX, zoneEditorSnapCheckboxY))
	this.base.runner.NextFrame()
	this.verifySnapshot()
	return this
}

// ClickApply commits the edited layout to the editor state and closes the
// dialog.
func (this *ZoneEditorHandler) ClickApply() *BaseHandler {
	this.base.runner.tb.Helper()
	this.base.runner.ClickButtonIn(zoneEditorRect(), applyChangesButtonLabel)
	this.base.commit()
	this.verifySnapshot()
	return this.base
}

// ClickCancel discards the edited layout and closes the dialog.
func (this *ZoneEditorHandler) ClickCancel() *BaseHandler {
	this.base.runner.tb.Helper()
	this.base.runner.ClickButtonIn(zoneEditorRect(), zoneEditorCancelButtonLabel)
	this.base.commit()
	this.verifySnapshot()
	return this.base
}

// Close dismisses the dialog without going through its footer, for the tests
// that only need the editor back.
func (this *ZoneEditorHandler) Close() *BaseHandler {
	this.base.runner.tb.Helper()
	this.base.runner.CloseTopDialog()
	return this.base
}

// TypeZoneSize types into the selected zone's size field. Gio inserts at the
// caret, so this adds to whatever the field already shows.
func (this *ZoneEditorHandler) TypeZoneSize(text string) *ZoneEditorHandler {
	return this.typeInSidePanelField(this.zoneRowY(zoneEditorZoneSizeY), text)
}

// TypeZoneGuard types into the selected zone's guard multiplier field.
func (this *ZoneEditorHandler) TypeZoneGuard(text string) *ZoneEditorHandler {
	return this.typeInSidePanelField(this.zoneRowY(zoneEditorZoneGuardY), text)
}

// TypeZoneWeekly types into the selected zone's weekly increment field.
func (this *ZoneEditorHandler) TypeZoneWeekly(text string) *ZoneEditorHandler {
	return this.typeInSidePanelField(this.zoneRowY(zoneEditorZoneWeeklyY), text)
}

// SelectZoneQuality picks a quality preset for the selected neutral zone.
func (this *ZoneEditorHandler) SelectZoneQuality(label string) *ZoneEditorHandler {
	return this.selectInSidePanelDropdown(this.zoneRowY(zoneEditorZoneQualityY), label)
}

// SelectZoneCastles picks how many castles the selected neutral zone holds.
func (this *ZoneEditorHandler) SelectZoneCastles(label string) *ZoneEditorHandler {
	return this.selectInSidePanelDropdown(this.zoneRowY(zoneEditorZoneCastlesY), label)
}

// ClickDeleteZone removes the selected zone. The button is disabled on a player
// spawn, whose content the generator owns.
func (this *ZoneEditorHandler) ClickDeleteZone() *ZoneEditorHandler {
	return this.clickDialogButton(deleteThisZoneButtonLabel)
}

// SelectConnectionType picks the selected connection's type.
func (this *ZoneEditorHandler) SelectConnectionType(label string) *ZoneEditorHandler {
	return this.selectInSidePanelDropdown(zoneEditorConnectionTypeY, label)
}

// SelectConnectionGuardZone picks which of the connected zones the guard is
// drawn from.
func (this *ZoneEditorHandler) SelectConnectionGuardZone(label string) *ZoneEditorHandler {
	return this.selectInSidePanelDropdown(zoneEditorConnectionGuardZoneY, label)
}

// SelectConnectionGuardPreset picks a guard strength preset, which rewrites the
// guard value field below it.
func (this *ZoneEditorHandler) SelectConnectionGuardPreset(label string) *ZoneEditorHandler {
	return this.selectInSidePanelDropdown(zoneEditorConnectionGuardPresetY, label)
}

// SelectConnectionWeekly picks the connection's weekly increment mode.
func (this *ZoneEditorHandler) SelectConnectionWeekly(label string) *ZoneEditorHandler {
	return this.selectInSidePanelDropdown(zoneEditorConnectionWeeklyY, label)
}

// TypeConnectionGuardValue types into the connection's guard value field.
func (this *ZoneEditorHandler) TypeConnectionGuardValue(text string) *ZoneEditorHandler {
	return this.typeInSidePanelField(zoneEditorConnectionGuardValueY, text)
}

// TypeConnectionIncrement types into the connection's increment field.
func (this *ZoneEditorHandler) TypeConnectionIncrement(text string) *ZoneEditorHandler {
	return this.typeInSidePanelField(zoneEditorConnectionIncrementY, text)
}

// TypeConnectionMatchGroup types into the connection's match group field, which
// only exists while the advanced options are shown.
func (this *ZoneEditorHandler) TypeConnectionMatchGroup(text string) *ZoneEditorHandler {
	return this.typeInSidePanelField(zoneEditorConnectionMatchGroupY, text)
}

// ToggleAdvancedOptions shows or hides the connection's advanced rows.
func (this *ZoneEditorHandler) ToggleAdvancedOptions() *ZoneEditorHandler {
	return this.clickSidePanelCheckbox(zoneEditorConnectionAdvancedY)
}

// ToggleGuardEscape flips the connection's guard escape flag.
func (this *ZoneEditorHandler) ToggleGuardEscape() *ZoneEditorHandler {
	return this.clickSidePanelCheckbox(zoneEditorConnectionGuardEscapeY)
}

// ToggleSimTurnSquad flips the connection's simultaneous turn squad flag.
func (this *ZoneEditorHandler) ToggleSimTurnSquad() *ZoneEditorHandler {
	return this.clickSidePanelCheckbox(zoneEditorConnectionSimTurnSquadY)
}

// ClickDeleteConnection removes the selected connection.
func (this *ZoneEditorHandler) ClickDeleteConnection() *ZoneEditorHandler {
	return this.clickDialogButton(deleteConnectionButtonLabel)
}

// dragBetween presses one canvas-local position and releases at another,
// travelling through interpolated moves so the canvas sees a drag rather than a
// jump.
func (this *ZoneEditorHandler) dragBetween(
	from models.Position, to models.Position) *ZoneEditorHandler {
	this.base.runner.tb.Helper()
	this.base.runner.DragTo(toWindowPixel(this.CanvasPoint(from)), toWindowPixel(this.CanvasPoint(to)))
	this.base.runner.NextFrame()
	this.verifySnapshot()
	return this
}

// clickDialogButton presses a labelled button anywhere inside the dialog panel.
func (this *ZoneEditorHandler) clickDialogButton(label string) *ZoneEditorHandler {
	this.base.runner.tb.Helper()
	this.base.runner.ClickButtonIn(zoneEditorRect(), label)
	this.base.runner.NextFrame()
	this.verifySnapshot()
	return this
}

// typeInSidePanelField focuses the field on a side panel row and types into it.
func (this *ZoneEditorHandler) typeInSidePanelField(rowY int, text string) *ZoneEditorHandler {
	this.base.runner.tb.Helper()
	this.base.runner.ClickAt(f32.Pt(zoneEditorSidePanelFieldX, float32(rowY)))
	this.base.runner.InputText(text)
	this.base.runner.NextFrame()
	this.verifySnapshot()
	return this
}

// selectInSidePanelDropdown opens the dropdown on a side panel row and clicks
// the option labelled label. An open dropdown pushes every row below it down, so
// an option is always addressed by label rather than by coordinate.
func (this *ZoneEditorHandler) selectInSidePanelDropdown(rowY int, label string) *ZoneEditorHandler {
	this.base.runner.tb.Helper()
	this.base.runner.ClickAt(f32.Pt(zoneEditorSidePanelFieldX, float32(rowY)))
	this.base.runner.ClickButtonIn(zoneEditorSidePanelRect(), label)
	this.base.runner.NextFrame()
	this.verifySnapshot()
	return this
}

// clickSidePanelCheckbox flips the checkbox on a side panel row. A checkbox
// emits no labelled button, so it is addressed by coordinate.
func (this *ZoneEditorHandler) clickSidePanelCheckbox(rowY int) *ZoneEditorHandler {
	this.base.runner.tb.Helper()
	this.base.runner.ClickAt(f32.Pt(zoneEditorSidePanelButtonX, float32(rowY)))
	this.base.runner.NextFrame()
	this.verifySnapshot()
	return this
}

// zoneRowY shifts a row measured on a zone that carries a note row up for a
// neutral zone, which carries none.
func (this *ZoneEditorHandler) zoneRowY(noteRowY int) int {
	this.base.runner.tb.Helper()
	if zone_helpers.IsZoneNameNeutral(this.Dialog().SelectedZone()) {
		return noteRowY - zoneEditorSidePanelNoteDrop
	}

	return noteRowY
}

// connectionMidPoint reads where the canvas currently draws the middle of the
// named connection's curve.
func (this *ZoneEditorHandler) connectionMidPoint(name string) models.Position {
	this.base.runner.tb.Helper()
	for _, edge := range this.Dialog().EdgeGeometries() {
		if edge.Name == name {
			return edge.MidPoint
		}
	}

	this.base.runner.tb.Fatalf("the canvas is not showing a connection called %q", name)
	return models.Position{}
}

// verifySnapshot compares the golden for the action that just ran, refusing to
// snapshot a dialog whose layout the Random topology re-randomises every run.
func (this *ZoneEditorHandler) verifySnapshot() {
	this.base.runner.tb.Helper()
	if !this.base.runner.SnapshotsEnabled() {
		return
	}
	if this.base.isRandomTopology {
		this.base.runner.tb.Fatal(
			"zone editor snapshots need LayoutAndZonesTabHandler.SelectTopology off Random")
	}

	this.base.runner.VerifySnapshot()
}

// toWindowPixel rounds a window point to the integer pixel the drag injector
// travels through.
func toWindowPixel(point f32.Point) image.Point {
	return image.Pt(int(point.X+0.5), int(point.Y+0.5))
}
