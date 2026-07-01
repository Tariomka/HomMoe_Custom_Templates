//go:build integration_test

// This file exposes Window internals for out-of-package tests (the test tree
// lives in a different directory, so the standard export_test.go mechanism does
// not reach it). It is guarded by the "integration_test" build tag and is
// therefore compiled ONLY for test builds run with `-tags integration_test`;
// production builds (`go build ./...`) exclude it entirely, keeping the public
// API clean.
package editor

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
)

// LoadStateFromFile loads a saved editor state from path and resynchronises the
// panels, mirroring what happens when the user picks a file in the Load dialog
// (but without the interactive picker). Exposed for programmatic loads and
// integration tests.
func (this *Window) LoadStateFromFile(path string) error {
	return this.state.LoadStateFromFile(path, this.load)
}

// SaveStateToFile writes the current editor state, including manual zone edits,
// to path. Exposed for programmatic saves and integration tests.
func (this *Window) SaveStateToFile(path string) error {
	this.save()
	return this.state.SaveStateToFile(path)
}

// CurrentState returns a copy of the live editor state. Exposed for tests.
func (this *Window) CurrentState() dtos.EditorStateDto {
	return this.state.GetStateData()
}

// TabCount returns the number of editor tabs. Exposed for tests and benchmarks.
func (this *Window) TabCount() int {
	return len(this.tabs)
}

// SelectedTabIndex returns the index of the currently selected tab. Exposed for
// tests and benchmarks that need to verify tab navigation.
func (this *Window) SelectedTabIndex() int {
	return this.selectedTab
}

// DialogsOpen reports whether a modal dialog is currently open. Exposed so tests
// and benchmarks can detect (and dismiss) dialogs they may trigger while
// programmatically clicking around the UI.
func (this *Window) DialogsOpen() bool {
	return this.state.Dialogs().IsOpen()
}

// CloseTopDialog dismisses the top-most modal dialog, if any. Exposed for tests
// and benchmarks.
func (this *Window) CloseTopDialog() {
	this.state.Dialogs().Close()
}
