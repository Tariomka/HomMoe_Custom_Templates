//go:build integration_test

package drivers

// SaveStateToFile ONLY FOR INTEGRATION TEST USE
func (this *State) SaveStateToFile(path string) {
	this.handleSaveState(path)
}

// LoadStateFromFile ONLY FOR INTEGRATION TEST USE
func (this *State) LoadStateFromFile(path string) {
	this.handleLoadState(path)
}

// SetCurrentPath seeds the file the editor is working on, which is also what
// the Load and Save To dialogs open at. ONLY FOR INTEGRATION TEST USE
func (this *State) SetCurrentPath(path string) {
	this.currentPath = path
}
