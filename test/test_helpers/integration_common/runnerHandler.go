//go:build integration_test

package integration_common

import (
	"image"

	"gioui.org/f32"
)

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
	this.runner.ClickAt(f32.Pt(672, 60))
	this.runner.VerifySnapshot()
	return this
}

func (this *baseHandler) ClickLayoutAndZonesTab() *baseHandler {
	this.runner.tb.Helper()
	this.runner.ClickAt(f32.Pt(789, 60))
	this.runner.VerifySnapshot()
	return this
}

func (this *baseHandler) ClickBonusesAndBansTab() *baseHandler {
	this.runner.tb.Helper()
	this.runner.ClickAt(f32.Pt(933, 60))
	this.runner.VerifySnapshot()
	return this
}

func (this *baseHandler) setRandomTopology() *baseHandler {
	this.runner.tb.Helper()
	this.isRandomTopology = true
	previewMask := image.Rect(
		WindowWidth-470, 0,
		WindowWidth, WindowHeight)
	this.runner.MaskRect(previewMask)
	return this
}
