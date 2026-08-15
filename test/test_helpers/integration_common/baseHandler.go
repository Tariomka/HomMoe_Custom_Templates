//go:build integration_test

package integration_common

import (
	"gioui.org/f32"
)

// BaseHandler drives the editor window through the AppRunner with a fluent API,
// standing in for the Playwright/Selenium-style driver that Gio has no
// equivalent of. Per-tab and per-dialog handlers, the state they need to track,
// and the scrolling support they depend on are designed in
// plans/gui-handler-framework.md; read that before adding to this file.
type BaseHandler struct {
	runner *AppRunner

	isRandomTopology bool

	// isExperimentalMapSizes tracks the checkbox that grows the map size
	// dropdown from 11 rows to 28, because the row a coordinate addresses is
	// only valid against the option count the checkbox implies.
	isExperimentalMapSizes bool
}

func NewHandler(runner *AppRunner) *BaseHandler {
	runner.tb.Helper()
	handler := BaseHandler{runner: runner}
	handler.setRandomTopology()
	return &handler
}

func (this *BaseHandler) WithSnapshots() *BaseHandler {
	this.runner.tb.Helper()
	this.runner.EnableSnapshots()
	return this
}

func (this *BaseHandler) ClickGeneralTab() *GeneralTabHandler {
	this.runner.tb.Helper()
	this.runner.ClickAt(f32.Pt(generalTabX, tabStripCenterY))
	this.runner.VerifySnapshot()
	return &GeneralTabHandler{BaseHandler: this}
}

func (this *BaseHandler) ClickLayoutAndZonesTab() *LayoutAndZonesTabHandler {
	this.runner.tb.Helper()
	this.runner.ClickAt(f32.Pt(layoutAndZonesTabX, tabStripCenterY))
	this.runner.VerifySnapshot()
	return &LayoutAndZonesTabHandler{BaseHandler: this}
}

func (this *BaseHandler) ClickBonusesAndBansTab() *BonusesAndBansTabHandler {
	this.runner.tb.Helper()
	this.runner.ClickAt(f32.Pt(bonusesAndBansTabX, tabStripCenterY))
	this.runner.VerifySnapshot()
	return &BonusesAndBansTabHandler{BaseHandler: this}
}

// ClickNew discards the editor state and starts a fresh template. The toolbar
// methods take no snapshot: Load and Save As open a dialog listing the
// per-machine templates directory (AGENTS.md 2.7), which no golden can hold.
func (this *BaseHandler) ClickNew() *BaseHandler {
	this.runner.tb.Helper()
	this.runner.ClickAt(f32.Pt(newButtonX, toolbarCenterY))
	return this
}

func (this *BaseHandler) ClickLoad() *FileExplorerHandler {
	this.runner.tb.Helper()
	this.runner.ClickAt(f32.Pt(loadButtonX, toolbarCenterY))
	return &FileExplorerHandler{base: this}
}

func (this *BaseHandler) ClickSaveAs() *FileExplorerHandler {
	this.runner.tb.Helper()
	this.runner.ClickAt(f32.Pt(saveAsButtonX, toolbarCenterY))
	return &FileExplorerHandler{base: this}
}

// ScrollPanel turns the mouse wheel over the settings panel by delta pixels;
// positive scrolls the content up. Gio clamps to the panel's scrollable range,
// so an oversized delta scrolls to the end.
func (this *BaseHandler) ScrollPanel(delta float32) *BaseHandler {
	this.runner.tb.Helper()
	this.runner.Scroll(f32.Pt(panelScrollX, panelScrollY), f32.Pt(0, delta))
	this.runner.VerifySnapshot()
	return this
}

// commit renders one more frame so the panel's SaveToState runs. A click is
// processed on the frame it is queued against, but the panel only writes the
// resulting widget values back to the editor state on the layout after that.
func (this *BaseHandler) commit() {
	this.runner.tb.Helper()
	this.runner.NextFrame()
}

// setRandomTopology records that the editor starts on the Random topology,
// whose preview is regenerated with fresh randomness on every run, and masks the
// three regions of the preview panel that cannot be compared against a golden.
func (this *BaseHandler) setRandomTopology() *BaseHandler {
	this.runner.tb.Helper()
	this.isRandomTopology = true
	this.runner.MaskRect(previewCanvasMask())
	this.runner.MaskRect(statusMessageMask())
	this.runner.MaskRect(outputDirectoryMask())
	return this
}
