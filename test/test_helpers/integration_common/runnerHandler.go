//go:build integration_test

package integration_common

import "gioui.org/f32"

type baseHandler struct {
	runner *AppRunner
}

func NewHandler(runner *AppRunner) *baseHandler {
	return &baseHandler{
		runner: runner,
	}
}

func (this *baseHandler) ClickGeneralTab() *baseHandler {
	this.runner.ClickAt(f32.Pt(672, 60))
	return this
}

func (this *baseHandler) ClickLayoutAndZonesTab() *baseHandler {
	this.runner.ClickAt(f32.Pt(789, 60))
	return this
}

func (this *baseHandler) ClickBonusesAndBansTab() *baseHandler {
	this.runner.ClickAt(f32.Pt(933, 60))
	return this
}
