//go:build integration_test

package editor

import (
	"image"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/dialogs"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

// LoadStateFromFile ONLY FOR INTEGRATION TEST USE
func (this *Window) LoadStateFromFile(path string) {
	this.state.LoadStateFromFile(path)
	this.load()
}

// SaveStateToFile ONLY FOR INTEGRATION TEST USE
func (this *Window) SaveStateToFile(path string) {
	this.save()
	this.state.SaveStateToFile(path)
}

// CurrentState ONLY FOR INTEGRATION TEST USE
func (this *Window) CurrentState() editor_state_dto.EditorStateDto {
	return this.state.GetStateData()
}

// TabCount ONLY FOR INTEGRATION TEST USE
func (this *Window) TabCount() int {
	return len(this.tabs)
}

// SelectedTabIndex ONLY FOR INTEGRATION TEST USE
func (this *Window) SelectedTabIndex() int {
	return this.selectedTab
}

// DialogsOpen ONLY FOR INTEGRATION TEST USE
func (this *Window) DialogsOpen() bool {
	return this.state.GetDialogHost().IsOpen()
}

// CloseTopDialog ONLY FOR INTEGRATION TEST USE
func (this *Window) CloseTopDialog() {
	this.state.GetDialogHost().Close()
}

// GetStateDriver ONLY FOR INTEGRATION TEST USE
func (this *Window) GetStateDriver() *drivers.State {
	return this.state
}

// SetTemplateName ONLY FOR INTEGRATION TEST USE
// Seeds the name the Save To dialog derives its filename from. The panels are
// resynced afterwards because every layout writes their widget values back over
// the editor state, which would otherwise undo this on the very next frame.
func (this *Window) SetTemplateName(name string) {
	this.state.UpdateState(func(state *editor_state_dto.EditorStateDto) { state.TemplateName = name })
	this.load()
}

// scrollablePanel is satisfied by the panels that expose their list position
// through a *_testexports.go file.
type scrollablePanel interface {
	ScrollPosition() (int, int)
}

// IFileExplorerDialog is the observation surface of dialogs.FileExplorerDialog,
// satisfied by the accessors on its own *_testexports.go file. It is declared
// here instead of in a *Interface.go file of its own because outside test/ only
// *_testexports.go may carry the integration_test tag (AGENTS.md 4.6.1), and a
// test-only contract must not reach production builds.
type IFileExplorerDialog interface {
	CurrentDir() string
	EntryNames() []string
	SelectedPath() string
	ScrollPosition() (int, int)
	ResolvedSaveName() string
	SaveNameReadOnly() bool
	ConfirmDisabled() bool
	NewFolderActive() bool
	OverwriteActive() bool
	SaveError() string
	NewFolderError() string
}

// TopFileExplorer ONLY FOR INTEGRATION TEST USE
// Returns the top-most dialog as a file explorer, and whether it is one.
func (this *Window) TopFileExplorer() (IFileExplorerDialog, bool) {
	dialog, ok := this.state.GetDialogHost().GetTopDialog().(IFileExplorerDialog)
	return dialog, ok
}

// IZoneEditorDialog is the observation surface of dialogs.ZoneEditorDialog,
// satisfied by the accessors on its own *_testexports.go file. Like
// IFileExplorerDialog it is declared here rather than in a *Interface.go file
// because outside test/ only *_testexports.go may carry the integration_test
// tag (AGENTS.md 4.6.1), and a test-only contract must not reach production
// builds. It is deliberately read-only: a test drives the dialog by clicking
// and dragging the real window, and only reads state back through this.
type IZoneEditorDialog interface {
	CanvasOrigin() image.Point
	CanvasSquareSide() int
	CanvasGridStep() float64
	CanvasZoneRadius() float64
	ZonePositions() map[string]models.Position
	EdgeGeometries() []dialogs.EdgeGeometry
	HitTestCanvasNode(pos models.Position) string
	HitTestCanvasEdge(pos models.Position) string
	SelectedZone() string
	SelectedConnection() string
	EditedZones() []entities.Zone
	EditedConnectionNames() []string
	AddConnectionModeActive() bool
	AddZoneModeActive() bool
	DraggingZone() string
	PendingConnectionSource() string
	SnapEnabled() bool
	SnapGuides() (x float64, xActive bool, y float64, yActive bool)
	StatusHint() string
}

// TopZoneEditor ONLY FOR INTEGRATION TEST USE
// Returns the top-most dialog as a zone editor, and whether it is one.
func (this *Window) TopZoneEditor() (IZoneEditorDialog, bool) {
	dialog, ok := this.state.GetDialogHost().GetTopDialog().(IZoneEditorDialog)
	return dialog, ok
}

// SelectedPanelScrollPosition ONLY FOR INTEGRATION TEST USE
// Returns the selected panel's first visible child and its pixel offset, plus
// whether that panel exposes a position at all.
func (this *Window) SelectedPanelScrollPosition() (int, int, bool) {
	panel, ok := this.tabs[this.selectedTab].GetPanel().(scrollablePanel)
	if !ok {
		return 0, 0, false
	}
	first, offset := panel.ScrollPosition()
	return first, offset, true
}
