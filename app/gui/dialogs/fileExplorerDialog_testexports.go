//go:build integration_test

package dialogs

// ConfirmSave ONLY FOR INTEGRATION TEST USE
func (this *FileExplorerDialog) ConfirmSave(path string) {
	if this.onSave != nil {
		this.onSave(path)
	}
}
