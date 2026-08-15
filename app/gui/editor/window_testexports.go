//go:build integration_test

package editor

import (
	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
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
func (this *Window) CurrentState() dtos.EditorStateDto {
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

// scrollablePanel is satisfied by the panels that expose their list position
// through a *_testexports.go file.
type scrollablePanel interface {
	ScrollPosition() (int, int)
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
