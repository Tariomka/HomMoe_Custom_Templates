package drivers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/dialogs"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
)

// Load opens the file picker and installs the picked .gen.json editor state.
// The picker is asynchronous, so callers must NOT resync their widgets at
// call time - onLoaded runs inside the pick handler, after the loaded state
// has been installed and only when loading succeeded.
func (this *State) Load(onLoaded func()) {
	dir, err := os.Getwd() // Editor state by default is loaded from the same directory as the executable
	if err != nil {
		dir = this.suggestDirectory()
	}
	this.dialogs.Open(dialogs.NewOpenFileDialog(dir, []string{configFileExtension}, func(path string) {
		if this.handleLoadState(path) && onLoaded != nil {
			onLoaded()
		}
	}))
}

func (this *State) Save() {
	if this.currentPath == "" {
		this.SaveAs(this.innerState.GetCurrentState().TemplateName)
		return
	}

	this.handleSaveState(this.currentPath)
}

func (this *State) SaveAs(templateName string) {
	dir, err := os.Getwd() // Editor state by default is saved in the same directory as the executable
	if err != nil {
		dir = this.suggestDirectory()
	}
	defaultName := helpers.SanitizeFilename(strings.TrimSpace(templateName)) + configFileExtension
	this.dialogs.Open(dialogs.NewSaveFileDialog(dir, defaultName, func(path string) {
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
	this.dialogs.Open(dialogs.NewPickFolderDialog(this.outputPath.Text(), func(path string) {
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
	this.dialogs.Open(dialogs.NewBrowseDialog(strings.TrimSpace(this.outputPath.Text())))
}

func (this *State) handleSaveState(path string) {
	savedPath, err := this.handler.SaveState(dtos.EditorStateSaveDto{
		State:      new(this.innerState.GetCurrentState()),
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
	dto, warnings, err := this.handler.LoadState(path, true)
	if err != nil {
		this.SetStatus(fmt.Sprintf("Load failed: %v.", err), true)
		return false
	}

	this.innerState.OverrideState(*dto)
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

func (this *State) suggestDirectory() string {
	if this.currentPath != "" {
		return filepath.Dir(this.currentPath)
	}

	if outputDir := strings.TrimSpace(this.outputPath.Text()); outputDir != "" {
		return outputDir
	}

	workingDir, _ := os.Getwd()
	return workingDir
}
