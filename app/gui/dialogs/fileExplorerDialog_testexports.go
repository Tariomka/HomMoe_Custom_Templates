//go:build integration_test

package dialogs

// ConfirmSave ONLY FOR INTEGRATION TEST USE
func (this *FileExplorerDialog) ConfirmSave(path string) {
	if this.onSave != nil {
		this.onSave(path)
	}
}

// ResolvedSaveName returns the read-only name the save will use.
// ONLY FOR INTEGRATION TEST USE
func (this *FileExplorerDialog) ResolvedSaveName() string { return this.filenameEd.Text() }

// SaveNameReadOnly reports whether the save-name field refuses edits. It only
// answers truthfully after a frame has been laid out, because the widget's
// read-only flag is applied by the textbox widget.
// ONLY FOR INTEGRATION TEST USE
func (this *FileExplorerDialog) SaveNameReadOnly() bool { return this.filenameEd.ReadOnly }

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

// ScrollPosition returns the listing's first visible row and its pixel offset.
// ONLY FOR INTEGRATION TEST USE
func (this *FileExplorerDialog) ScrollPosition() (int, int) {
	return this.list.Position.First, this.list.Position.Offset
}

// OverwriteActive ONLY FOR INTEGRATION TEST USE
func (this *FileExplorerDialog) OverwriteActive() bool { return this.overwriteActive }

// NewFolderActive reports whether the inline new-folder row is showing.
// ONLY FOR INTEGRATION TEST USE
func (this *FileExplorerDialog) NewFolderActive() bool { return this.newFolderActive }

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
