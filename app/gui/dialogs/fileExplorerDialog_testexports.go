//go:build integration_test

package dialogs

// ConfirmSave ONLY FOR INTEGRATION TEST USE
func (this *FileExplorerDialog) ConfirmSave(path string) {
	if this.onSave != nil {
		this.onSave(path)
	}
}

// ClickEntry queues a click on the listed entry called name and reports whether
// such an entry exists. ONLY FOR INTEGRATION TEST USE
func (this *FileExplorerDialog) ClickEntry(name string) bool {
	for _, entry := range this.entries {
		if entry.Name == name {
			this.clickFor(entry.Path).Click()
			return true
		}
	}

	return false
}

// ClickUp ONLY FOR INTEGRATION TEST USE
func (this *FileExplorerDialog) ClickUp() { this.upBtn.Click() }

// ClickConfirm ONLY FOR INTEGRATION TEST USE
func (this *FileExplorerDialog) ClickConfirm() { this.confirmBtn.Click() }

// ClickOverwriteConfirm ONLY FOR INTEGRATION TEST USE
func (this *FileExplorerDialog) ClickOverwriteConfirm() { this.overwriteConfirmBtn.Click() }

// ClickOverwriteCancel ONLY FOR INTEGRATION TEST USE
func (this *FileExplorerDialog) ClickOverwriteCancel() { this.overwriteCancelBtn.Click() }

// ClickNewFolder ONLY FOR INTEGRATION TEST USE
func (this *FileExplorerDialog) ClickNewFolder() { this.newFolderBtn.Click() }

// ClickCreateFolder ONLY FOR INTEGRATION TEST USE
func (this *FileExplorerDialog) ClickCreateFolder() { this.createFolderBtn.Click() }

// SetFilename ONLY FOR INTEGRATION TEST USE
func (this *FileExplorerDialog) SetFilename(name string) { this.filenameEd.SetText(name) }

// SetNewFolderName ONLY FOR INTEGRATION TEST USE
func (this *FileExplorerDialog) SetNewFolderName(name string) { this.newFolderEd.SetText(name) }

// CurrentDir ONLY FOR INTEGRATION TEST USE
func (this *FileExplorerDialog) CurrentDir() string { return this.currentDir }

// EntryNames returns the names of the cached listing, in display order.
// ONLY FOR INTEGRATION TEST USE
func (this *FileExplorerDialog) EntryNames() []string {
	names := make([]string, 0, len(this.entries))
	for _, entry := range this.entries {
		names = append(names, entry.Name)
	}

	return names
}

// SelectedPath ONLY FOR INTEGRATION TEST USE
func (this *FileExplorerDialog) SelectedPath() string { return this.selectedPath }

// OverwriteActive ONLY FOR INTEGRATION TEST USE
func (this *FileExplorerDialog) OverwriteActive() bool { return this.overwriteActive }

// SaveError ONLY FOR INTEGRATION TEST USE
func (this *FileExplorerDialog) SaveError() string { return this.saveErr }

// NewFolderError ONLY FOR INTEGRATION TEST USE
func (this *FileExplorerDialog) NewFolderError() string { return this.newFolderErr }

// ConfirmDisabled reports whether the primary button is currently greyed out.
// ONLY FOR INTEGRATION TEST USE
func (this *FileExplorerDialog) ConfirmDisabled() bool {
	_, _, disabled := this.confirmButtonState()
	return disabled
}
