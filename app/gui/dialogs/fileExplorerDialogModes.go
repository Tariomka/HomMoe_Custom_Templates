package dialogs

import "github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"

// fileDialogMode selects which task the explorer performs.
type fileDialogMode uint8

const (
	modeOpenFile fileDialogMode = iota
	modeSaveFile
	modePickFolder
	modeBrowse
)

// NewOpenFileDialog builds a single-file picker starting at initialDir. Only
// files whose name ends with one of filterSuffixes (case-insensitive) are shown;
// pass nil to list every file. onPick receives the chosen absolute file path.
func NewOpenFileDialog(
	fileSystem handler_interfaces.IFileSystemHandler,
	initialDir string,
	filterSuffixes []string,
	onPick func(path string)) *FileExplorerDialog {
	dialog := newFileExplorerDialog(fileSystem, modeOpenFile, "Open File")
	dialog.filterSuffixes = filterSuffixes
	dialog.onPick = onPick
	dialog.loadDir(fileSystem.ResolveStartDirectory(initialDir))
	return dialog
}

// NewSaveFileDialog builds a save-location picker starting at initialDir with
// the filename field prefilled to defaultName. onSave receives the chosen
// absolute path (with a guaranteed .gen.json suffix) after any overwrite
// confirmation.
func NewSaveFileDialog(
	fileSystem handler_interfaces.IFileSystemHandler,
	initialDir, defaultName string,
	onSave func(path string)) *FileExplorerDialog {
	dialog := newFileExplorerDialog(fileSystem, modeSaveFile, "Save File")
	dialog.onSave = onSave
	dialog.filenameEd.SetText(defaultName)
	dialog.loadDir(fileSystem.ResolveStartDirectory(initialDir))
	return dialog
}

// NewPickFolderDialog builds a single-folder picker starting at initialDir.
// onPick receives the directory the user navigated into and confirmed.
func NewPickFolderDialog(
	fileSystem handler_interfaces.IFileSystemHandler,
	initialDir string,
	onPick func(dir string)) *FileExplorerDialog {
	dialog := newFileExplorerDialog(fileSystem, modePickFolder, "Select Folder")
	dialog.onPick = onPick
	dialog.loadDir(fileSystem.ResolveStartDirectory(initialDir))
	return dialog
}

// NewBrowseDialog builds a read-only viewer starting at initialDir; it has no
// confirm action and is used to inspect a directory's contents in-app.
func NewBrowseDialog(
	fileSystem handler_interfaces.IFileSystemHandler,
	initialDir string) *FileExplorerDialog {
	dialog := newFileExplorerDialog(fileSystem, modeBrowse, "Browse")
	dialog.loadDir(fileSystem.ResolveStartDirectory(initialDir))
	return dialog
}
