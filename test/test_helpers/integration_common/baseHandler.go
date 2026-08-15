//go:build integration_test

package integration_common

import (
	"gioui.org/f32"
)

// baseHandler drives the editor window through the AppRunner with a fluent API,
// standing in for the Playwright/Selenium-style driver that Gio has no
// equivalent of. Per-tab and per-dialog handlers, the state they need to track,
// and the scrolling support they depend on are designed in todo/backlog-opus5.md
// section 5.4 (items d-f); read that before adding to this file.
type baseHandler struct {
	runner *AppRunner

	isRandomTopology bool
}

func NewHandler(runner *AppRunner) *baseHandler {
	runner.tb.Helper()
	handler := baseHandler{runner: runner}
	handler.setRandomTopology()
	return &handler
}

func (this *baseHandler) WithSnapshots() *baseHandler {
	this.runner.tb.Helper()
	this.runner.EnableSnapshots()
	return this
}

func (this *baseHandler) ClickGeneralTab() *baseHandler {
	this.runner.tb.Helper()
	this.runner.ClickAt(f32.Pt(generalTabX, tabStripCenterY))
	this.runner.VerifySnapshot()
	return this
}

func (this *baseHandler) ClickLayoutAndZonesTab() *baseHandler {
	this.runner.tb.Helper()
	this.runner.ClickAt(f32.Pt(layoutAndZonesTabX, tabStripCenterY))
	this.runner.VerifySnapshot()
	return this
}

func (this *baseHandler) ClickBonusesAndBansTab() *baseHandler {
	this.runner.tb.Helper()
	this.runner.ClickAt(f32.Pt(bonusesAndBansTabX, tabStripCenterY))
	this.runner.VerifySnapshot()
	return this
}

// setRandomTopology records that the editor starts on the Random topology,
// whose preview is regenerated with fresh randomness on every run, and masks the
// three regions of the preview panel that cannot be compared against a golden.
func (this *baseHandler) setRandomTopology() *baseHandler {
	this.runner.tb.Helper()
	this.isRandomTopology = true
	this.runner.MaskRect(previewCanvasMask())
	this.runner.MaskRect(statusMessageMask())
	this.runner.MaskRect(outputDirectoryMask())
	return this
}
