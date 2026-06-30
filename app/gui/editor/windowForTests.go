package editor

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
)

type WindowForTests struct {
	Window
}

func NewWindowForTests() *WindowForTests {
	return &WindowForTests{Window: *NewWindow()}
}

// LoadStateFromFile loads a saved editor state from path and resynchronises the
// panels, mirroring what happens when the user picks a file in the Load dialog
// (but without the interactive picker). Exposed for programmatic loads and
// integration tests.
func (this *WindowForTests) LoadStateFromFile(path string) error {
	return this.state.LoadStateFromFile(path, this.load)
}

// SaveStateToFile writes the current editor state, including manual zone edits,
// to path. Exposed for programmatic saves and integration tests.
func (this *WindowForTests) SaveStateToFile(path string) error {
	this.save()
	return this.state.SaveStateToFile(path)
}

// CurrentState returns a copy of the live editor state. Exposed for tests.
func (this *WindowForTests) CurrentState() dtos.EditorStateDto {
	return this.state.GetStateData()
}

// TabCount returns the number of editor tabs. Exposed for tests and benchmarks.
func (this *WindowForTests) TabCount() int {
	return len(this.tabs)
}

// SelectedTabIndex returns the index of the currently selected tab. Exposed for
// tests and benchmarks that need to verify tab navigation.
func (this *WindowForTests) SelectedTabIndex() int {
	return this.selectedTab
}

// DialogsOpen reports whether a modal dialog is currently open. Exposed so tests
// and benchmarks can detect (and dismiss) dialogs they may trigger while
// programmatically clicking around the UI.
func (this *WindowForTests) DialogsOpen() bool {
	return this.state.Dialogs().IsOpen()
}

// CloseTopDialog dismisses the top-most modal dialog, if any. Exposed for tests
// and benchmarks.
func (this *WindowForTests) CloseTopDialog() {
	this.state.Dialogs().Close()
}
