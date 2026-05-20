package components

import (
	"fmt"
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
	settings *models.SettingsFile

	// File state
	currentPath string
	dirty       bool

	// Output / status
	outputPath   widget.Editor
	lastTemplate *models.RmgTemplate
	statusMsg    string
	statusErr    bool
}

func NewUIState() *State {
	state := &State{settings: models.NewSettingsFile()}
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
	return this.settings
}

func (this *State) GetCurrentPath() string {
	return this.currentPath
}

func (this *State) IsDirty() bool {
	return this.dirty
}

func (this *State) GetLastTemplate() *models.RmgTemplate {
	return this.lastTemplate
}

func (this *State) GetOutputPath() string {
	return this.outputPath.Text()
}

func (this *State) Reset() {
	this.settings = models.NewSettingsFile()
	this.currentPath = ""
	this.dirty = false
	this.SetStatus("New settings file.", false)
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

func (this *State) Load() {
	path, err := utils.PickOpenFile("Open settings", "Settings (*.gen.json)|*.gen.json|All files|*.*", this.SuggestDirectory())
	if err != nil {
		this.SetStatus("Open dialog failed: "+err.Error(), true)
		return
	}

	if path == "" {
		return
	}

	loaded, err := services.LoadSettingsFile(path)
	if err != nil {
		this.SetStatus("Load failed: "+err.Error(), true)
		return
	}

	this.settings = loaded
	this.currentPath = path
	this.dirty = false
	this.SetStatus("Loaded "+path, false)
}

func (this *State) Save() {
	if this.currentPath == "" {
		this.SaveAs(this.settings.TemplateName)
		return
	}

	if err := services.SaveSettingsFile(this.currentPath, this.settings); err != nil {
		this.SetStatus("Save failed: "+err.Error(), true)
		return
	}

	this.dirty = false
	this.SetStatus("Saved "+this.currentPath, false)
}

func (this *State) SaveAs(templateName string) {
	defaultName := services.SanitizeFilename(strings.TrimSpace(templateName)) + ".gen.json"
	path, err := utils.PickSaveFile("Save settings as", "Settings (*.gen.json)|*.gen.json", this.SuggestDirectory(), defaultName)
	if err != nil {
		this.SetStatus("Save dialog failed: "+err.Error(), true)
		return
	}

	if path == "" {
		return
	}

	if err := services.SaveSettingsFile(path, this.settings); err != nil {
		this.SetStatus("Save failed: "+err.Error(), true)
		return
	}
	this.currentPath = path
	this.dirty = false
	this.SetStatus("Saved "+path, false)
}

func (this *State) Generate() {
	captured := this.GetSettingsFile()
	generatorSettings := services.SettingsToGenerator(captured)
	if generatorSettings.TemplateName == "" {
		this.SetStatus("Template name is required.", true)
		return
	}
	template, err := services.Generate(generatorSettings)
	if err != nil {
		this.SetStatus(fmt.Sprintf("Generation failed: %value", err), true)
		this.lastTemplate = nil
		return
	}
	this.lastTemplate = template
	zoneCount := 0
	connectionCount := 0
	if len(template.Variants) > 0 {
		zoneCount = len(template.Variants[0].Zones)
		connectionCount = len(template.Variants[0].Connections)
	}
	this.SetStatus(fmt.Sprintf("Generated '%s' — %d zones, %d connections.", template.Name, zoneCount, connectionCount), false)
}

// SaveTemplate writes the most recently generated template as .rmg.json.
func (this *State) SaveTemplate() {
	lastTemplate := this.GetLastTemplate()
	if lastTemplate == nil {
		this.SetStatus("Nothing to save — generate a template first.", true)
		return
	}
	dir := strings.TrimSpace(this.outputPath.Text())
	if dir == "" {
		this.SetStatus("Output directory is empty.", true)
		return
	}
	out, err := services.WriteTemplate(dir, lastTemplate)
	if err != nil {
		this.SetStatus(fmt.Sprintf("Save failed: %value", err), true)
		return
	}
	this.SetStatus("Saved template to "+out, false)
}

// PickOutputDir presents a folder picker for the template output directory.
func (this *State) PickOutputDir() {
	cur := strings.TrimSpace(this.outputPath.Text())
	dir, err := utils.PickFolder("Select output directory", cur)
	if err != nil {
		this.SetStatus("Folder dialog failed: "+err.Error(), true)
		return
	}
	if dir == "" {
		return
	}
	this.outputPath.SetText(dir)
}

func (this *State) UpdateState(updateFunc func(*models.SettingsFile)) {
	updateFunc(this.settings)
}
