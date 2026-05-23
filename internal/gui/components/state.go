package components

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gioui.org/widget"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/utils"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/services"
)

type State struct {
	// Persistent settings file model. Updated continuously from widgets.
	settings *models.SettingsFile

	// File state
	currentPath string
	unsaved     bool

	// Output / status
	outputPath   widget.Editor
	lastTemplate *template.RmgTemplateModel
	statusMsg    string
	statusErr    bool
}

func NewUIState() *State {
	state := &State{settings: models.NewSettingsFile()}
	state.outputPath.SingleLine = true

	templateDir := helpers.FindOldenEraTemplatesDir(false)
	if templateDir == "" {
		if workingDir, err := os.Getwd(); err == nil {
			templateDir = workingDir
		}
	}
	state.outputPath.SetText(templateDir)
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

func (this *State) IsUnsaved() bool {
	return this.unsaved
}

func (this *State) GetLastTemplate() *template.RmgTemplateModel {
	return this.lastTemplate
}

func (this *State) GetOutputPath() string {
	return this.outputPath.Text()
}

func (this *State) Reset() {
	this.settings = models.NewSettingsFile()
	this.currentPath = ""
	this.unsaved = false
	this.setStatus("New settings file.", false)
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
		this.setStatus("Open dialog failed: "+err.Error(), true)
		return
	}

	if path == "" {
		return
	}

	loaded, err := services.LoadSettingsFile(path)
	if err != nil {
		this.setStatus("Load failed: "+err.Error(), true)
		return
	}

	this.settings = loaded
	this.currentPath = path
	this.unsaved = false
	this.setStatus("Loaded "+path, false)
}

func (this *State) Save() {
	if this.currentPath == "" {
		this.SaveAs(this.settings.TemplateName)
		return
	}

	if err := services.SaveSettingsFile(this.currentPath, this.settings); err != nil {
		this.setStatus("Save failed: "+err.Error(), true)
		return
	}

	this.unsaved = false
	this.setStatus("Saved "+this.currentPath, false)
}

func (this *State) SaveAs(templateName string) {
	defaultName := services.SanitizeFilename(strings.TrimSpace(templateName)) + ".gen.json"
	path, err := utils.PickSaveFile("Save settings as", "Settings (*.gen.json)|*.gen.json", this.SuggestDirectory(), defaultName)
	if err != nil {
		this.setStatus("Save dialog failed: "+err.Error(), true)
		return
	}

	if path == "" {
		return
	}

	if err := services.SaveSettingsFile(path, this.settings); err != nil {
		this.setStatus("Save failed: "+err.Error(), true)
		return
	}
	this.currentPath = path
	this.unsaved = false
	this.setStatus("Saved "+path, false)
}

func (this *State) Generate() {
	generatorSettings := services.SettingsToGenerator(this.GetSettingsFile())
	if generatorSettings.TemplateName == "" {
		this.setStatus("Template name is required.", true)
		return
	}
	template, err := services.Generate(generatorSettings)
	if err != nil {
		this.setStatus(fmt.Sprintf("Generation failed: %value", err), true)
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
	this.setStatus(fmt.Sprintf("Generated '%s' — %d zones, %d connections.", template.Name, zoneCount, connectionCount), false)
}

// SaveTemplate writes the most recently generated template as .rmg.json.
func (this *State) SaveTemplate() {
	lastTemplate := this.GetLastTemplate()
	if lastTemplate == nil {
		this.setStatus("Nothing to save — generate a template first.", true)
		return
	}
	dir := strings.TrimSpace(this.outputPath.Text())
	if dir == "" {
		this.setStatus("Output directory is empty.", true)
		return
	}
	out, err := services.WriteTemplate(dir, lastTemplate)
	if err != nil {
		this.setStatus(fmt.Sprintf("Save failed: %value", err), true)
		return
	}
	this.setStatus("Saved template to "+out, false)
}

// PickOutputDir presents a folder picker for the template output directory.
func (this *State) PickOutputDir() {
	cur := strings.TrimSpace(this.outputPath.Text())
	dir, err := utils.PickFolder("Select output directory", cur)
	if err != nil {
		this.setStatus("Folder dialog failed: "+err.Error(), true)
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

func (this *State) setStatus(msg string, isErr bool) {
	this.statusMsg = msg
	this.statusErr = isErr
}
