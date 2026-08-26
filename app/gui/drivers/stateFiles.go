package drivers

import (
	"fmt"
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/dialogs"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
)

// Load loads in an editor state from a file. The file picker is asynchronous,
// so onLoaded is used to resync the UI only when loading succeeded.
func (this *State) Load(onLoaded func()) {
	dir := this.getWorkingDirectory()
	this.dialogs.Open(dialogs.NewOpenFileDialog(this.fileSystem, dir, []string{configFileExtension}, func(path string) {
		if this.handleLoadState(path) && onLoaded != nil {
			onLoaded()
		}
	}))
}

func (this *State) Save() {
	if this.currentPath == "" {
		this.SaveTo(this.innerState.GetTemplateName())
		return
	}

	this.handleSaveState(this.currentPath)
}

// SaveTo saves editor state to the file designated by templateName. The directory is picked via dialog.
func (this *State) SaveTo(templateName string) {
	dir := this.getWorkingDirectory()
	resolvedName := strings.TrimSpace(templateName)
	if resolvedName != "" {
		resolvedName = helpers.SanitizeFilename(resolvedName) + configFileExtension
	}

	this.dialogs.Open(dialogs.NewSaveFileDialog(this.fileSystem, dir, resolvedName, func(path string) {
		this.handleSaveState(path)
	}))
}

// SetOnExit installs the callback Exit uses to close the application window.
func (this *State) SetOnExit(onExit func()) { this.onExit = onExit }

func (this *State) Exit() {
	if this.unsaved && !this.confirmExit {
		this.SetStatus("Unsaved changes exist - save first or press Exit again.", true)
		this.confirmExit = true
		return
	}

	if this.onExit != nil {
		this.onExit()
	}
}

// PickOutputDir presents a folder picker for the template output directory.
func (this *State) PickOutputDir() {
	this.dialogs.Open(dialogs.NewPickFolderDialog(this.fileSystem, this.outputPath.Text(), func(path string) {
		if path == "" {
			this.SetStatus("No output directory selected.", true)
			return
		}

		path = strings.TrimSpace(path)
		this.outputPath.SetText(path)
	}))
}

// RevealOutputDir opens the in-app explorer in read-only Browse mode at the
// configured output directory.
func (this *State) RevealOutputDir() {
	this.dialogs.Open(dialogs.NewBrowseDialog(this.fileSystem, strings.TrimSpace(this.outputPath.Text())))
}

func (this *State) handleSaveState(path string) {
	savedPath, err := this.handler.SaveState(editor_state_dto.EditorStateSaveDto{
		State:      new(this.GetStateDto()),
		OutputPath: path,
	})
	if err != nil {
		this.SetStatus(fmt.Sprintf("Save failed: %v.", err), true)
		return
	}

	this.currentPath = savedPath
	this.unsaved = false
	this.confirmExit = false
	this.SetStatus("Saved "+savedPath, false)
}

func (this *State) handleLoadState(path string) bool {
	loaded, warnings, err := this.handler.LoadState(path, true)
	if err != nil {
		this.SetStatus(fmt.Sprintf("Load failed: %v.", err), true)
		return false
	}

	this.innerState.OverrideState(loaded.EditorState)
	this.currentPath = path
	this.unsaved = false
	this.clearGeneratedState()
	if len(warnings) > 0 {
		this.SetStatus(fmt.Sprintf("Loaded %s (adjusted: %s)", path, strings.Join(warnings, "; ")), false)
		return true
	}

	this.SetStatus("Loaded "+path, false)
	return true
}

// getWorkingDirectory returns the directory a file dialog should open at: the
// one holding the file currently being edited, so a second Load or Save To
// lands where the user last worked, falling back to the process working
// directory before anything has been opened or saved. Resolution goes through
// the filesystem handler so the driver performs no path arithmetic of its own -
// ResolveStartDirectory climbs from a file path to its containing directory and
// always yields an existing one, so unlike the raw [os.Getwd] it never needs a
// caller-side fallback.
func (this *State) getWorkingDirectory() string {
	if this.currentPath != "" {
		return this.fileSystem.ResolveStartDirectory(this.currentPath)
	}

	return this.fileSystem.ResolveStartDirectory(".")
}
