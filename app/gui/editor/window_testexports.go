//go:build integration_test

package editor

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
)

// LoadStateFromFile ONLY FOR INTEGRATION TEST USE
func (this *Window) LoadStateFromFile(path string) error {
	return this.state.LoadStateFromFile(path, this.load)
}

// SaveStateToFile  ONLY FOR INTEGRATION TEST USE
func (this *Window) SaveStateToFile(path string) error {
	this.save()
	return this.state.SaveStateToFile(path)
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
	return this.state.Dialogs().IsOpen()
}

// CloseTopDialog ONLY FOR INTEGRATION TEST USE
func (this *Window) CloseTopDialog() {
	this.state.Dialogs().Close()
}
