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
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator"
)

type State struct {
	// Persistent stateModel file model. Updated continuously from widgets.
	stateModel *models.EditorStateModel

	// File state
	currentPath string
	unsaved     bool

	// Output / status
	outputPath   widget.Editor
	lastTemplate *template.RmgTemplateModel
	statusMsg    string
	statusErr    bool

	// connectionsModified records that the live template's connections were
	// hand-edited in the visual editor, so a later regeneration can warn that
	// those edits will be replaced.
	connectionsModified bool

	// dialogs renders modal dialogs (rule editors, pickers, the connection
	// editor) over the main UI.
	dialogs *DialogHost
}

func NewUIState() *State {
	state := &State{stateModel: models.NewEditorStateModel()}
	state.outputPath.SingleLine = true
	state.dialogs = &DialogHost{}

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

// Dialogs returns the modal host used to open and render dialogs.
func (this *State) Dialogs() *DialogHost {
	return this.dialogs
}

func (this *State) GetStateData() models.EditorStateModel {
	return *this.stateModel
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

// ApplyEditedConnections writes connections edited in the visual connection
// editor back into the live template and flags that manual edits now exist.
func (this *State) ApplyEditedConnections(connections []template.Connection) {
	if this.lastTemplate == nil || len(this.lastTemplate.Variants) == 0 {
		return
	}
	this.lastTemplate.Variants[0].Connections = connections
	this.connectionsModified = true
	if connection_editor.ComputeHasErrors(this.lastTemplate.Variants[0].Zones, connections) {
		this.setStatus(fmt.Sprintf("Applied %d connections — \u26a0 some reference a missing zone; fix before export.", len(connections)), true)
		return
	}
	this.setStatus(fmt.Sprintf("Applied %d connections from the editor.", len(connections)), false)
}

func (this *State) GetOutputPath() string {
	return this.outputPath.Text()
}

func (this *State) Reset() {
	this.stateModel = models.NewEditorStateModel()
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

	this.stateModel = loaded
	this.currentPath = path
	this.unsaved = false
	this.setStatus("Loaded "+path, false)
}

func (this *State) Save() {
	if this.currentPath == "" {
		this.SaveAs(this.stateModel.TemplateName)
		return
	}

	if err := services.SaveSettingsFile(this.currentPath, this.stateModel); err != nil {
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

	if err := services.SaveSettingsFile(path, this.stateModel); err != nil {
		this.setStatus("Save failed: "+err.Error(), true)
		return
	}
	this.currentPath = path
	this.unsaved = false
	this.setStatus("Saved "+path, false)
}

func (this *State) Generate() {
	generatorSettings := services.SettingsToGenerator(this.stateModel)
	if generatorSettings.TemplateName == "" {
		this.setStatus("Template name is required.", true)
		return
	}
	wasModified := this.connectionsModified
	this.connectionsModified = false
	template := template_generator.NewTemplateGenerator(generatorSettings).Generate()
	// template, err := services.Generate(generatorSettings)
	// if err != nil {
	// 	this.setStatus(fmt.Sprintf("Generation failed: %value", err), true)
	// 	this.lastTemplate = nil
	// 	return
	// }
	this.lastTemplate = template
	zoneCount := 0
	connectionCount := 0
	if len(template.Variants) > 0 {
		zoneCount = len(template.Variants[0].Zones)
		connectionCount = len(template.Variants[0].Connections)
	}
	status := fmt.Sprintf("Generated '%s' — %d zones, %d connections.", template.Name, zoneCount, connectionCount)
	if wasModified {
		status += " (Manual connection edits were replaced by regeneration.)"
	}
	this.setStatus(status, false)
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

func (this *State) UpdateState(updateFunc func(*models.EditorStateModel)) {
	// TODO: add validator for state updates, e.g. to prevent invalid map sizes or player counts
	updateFunc(this.stateModel)
}

func (this *State) setStatus(msg string, isErr bool) {
	this.statusMsg = msg
	this.statusErr = isErr
}
