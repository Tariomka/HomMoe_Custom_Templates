package components

import (
	"os"
	"path/filepath"
	"strings"

	"gioui.org/widget"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/utils"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/services"
)

type State struct {
	// Persistent settings file model. Updated continuously from widgets.
	settingsFile *models.SettingsFile

	// File state
	currentPath string
	dirty       bool

	// Output / status
	outputPath   widget.Editor
	lastTemplate *models.RmgTemplate
	statusMsg    string
	statusErr    bool
}

func NewUiState() *State {
	state := &State{
		settingsFile: models.NewSettingsFile(),
	}
	state.outputPath.SingleLine = true
	if workingDir, err := os.Getwd(); err == nil {
		state.outputPath.SetText(workingDir)
	}
	return state
}

func (this *State) GetStatus() (msg string, isErr bool) {
	return this.statusMsg, this.statusErr
}

func (this *State) GetSettingsFile() *models.SettingsFile {
	return this.settingsFile
}

func (this *State) GetCurrentPath() string {
	return this.currentPath
}

func (this *State) GetOutputPath() string {
	return this.outputPath.Text()
}

func (this *State) Reset() {
	this.settingsFile = models.NewSettingsFile()
	this.currentPath = ""
	this.dirty = false
}

func (this *State) SetStatus(msg string, isErr bool) {
	this.statusMsg = msg
	this.statusErr = isErr
}

func (this *State) SuggestDirectory() string {
	if this.currentPath != "" {
		return filepath.Dir(this.currentPath)
	}
	if outputDir := strings.TrimSpace(this.outputPath.Text()); outputDir != "" {
		return outputDir
	}
	workingDir, _ := os.Getwd()
	return workingDir
}

func (this *State) Save(filename string) {
	defaultName := services.SanitizeFilename(strings.TrimSpace(filename)) + ".gen.json"
	path, err := utils.PickSaveFile("Save settings as", "Settings (*.gen.json)|*.gen.json", this.SuggestDirectory(), defaultName)
	if err != nil {
		this.SetStatus("Save dialog failed: "+err.Error(), true)
		return
	}
	if path == "" {
		return
	}
	if err := services.SaveSettingsFile(path, this.captureToSettingsFile()); err != nil {
		this.SetStatus("Save failed: "+err.Error(), true)
		return
	}
	this.currentPath = path
	this.dirty = false
	this.SetStatus("Saved "+path, false)
}
