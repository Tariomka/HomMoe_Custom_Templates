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
	loaded, err := this.handler.LoadState(path, false)
	if err != nil {
		this.SetStatus(fmt.Sprintf("Load failed: %v.", err), true)
		return false
	}

	validation := this.handler.ValidateEditorState(loaded.State, true)
	outcome := this.innerState.OverrideStateFromLoad(loaded.State, validation.State)
	this.currentPath = path
	this.unsaved = false
	this.clearGeneratedState()

	this.pendingOutcome = outcome
	if outcome.TournamentCountCorrected {
		this.flagAsUnsaved()
	}

	this.SetStatus(this.describeLoadedState(path, validation.Warnings), false)
	return true
}

func (this *State) describeLoadedState(path string, warnings []string) string {
	status := "Loaded " + path
	if len(warnings) > 0 {
		status = fmt.Sprintf("Loaded %s (adjusted: %s)", path, strings.Join(warnings, "; "))
	}

	notice := stateTransitionNotice(this.pendingOutcome)
	if notice == "" {
		return status
	}

	return status + " " + notice
}

func (this *State) getWorkingDirectory() string {
	if this.currentPath != "" {
		return this.fileSystem.ResolveStartDirectory(this.currentPath)
	}

	return this.fileSystem.ResolveStartDirectory(".")
}
