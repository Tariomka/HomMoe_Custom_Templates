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
