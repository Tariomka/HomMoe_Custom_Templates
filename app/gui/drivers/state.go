package drivers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gioui.org/widget"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/utils"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/services"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
)

type State struct {
	handler *handlers.GUIHandler

	// Persistent stateDto file model. Updated continuously from widgets.
	stateDto *dtos.EditorStateDto

	// File state
	currentPath string
	unsaved     bool

	// Output / status
	outputPath   widget.Editor
	lastTemplate *entities.RmgTemplate
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
	stateDto := dtos.NewDefaultEditorStateDto()
	state := &State{
		handler:  handlers.NewGuiHandler(),
		stateDto: &stateDto,
	}
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

func (this *State) GetStateData() dtos.EditorStateDto {
	return *this.stateDto
}

func (this *State) GetCurrentPath() string {
	return this.currentPath
}

func (this *State) IsUnsaved() bool {
	return this.unsaved
}

func (this *State) GetLastTemplate() *entities.RmgTemplate {
	return this.lastTemplate
}

// ApplyEditedZones writes zones and connections edited in the manual zone
// editor back into the live template and flags that manual edits now exist.
func (this *State) ApplyEditedZones(zones []entities.Zone, connections []entities.Connection) {
	if this.lastTemplate == nil || len(this.lastTemplate.Variants) == 0 {
		return
	}
	this.lastTemplate.Variants[0].Zones = zones
	this.lastTemplate.Variants[0].Connections = connections
	this.connectionsModified = true
	if connection_editor.ComputeHasErrors(zones, connections) {
		this.SetStatus(fmt.Sprintf("Applied %d zones and %d connections — \u26a0 some connections reference a missing zone; fix before export.", len(zones), len(connections)), true)
		return
	}
	this.SetStatus(fmt.Sprintf("Applied %d zones and %d connections from the editor.", len(zones), len(connections)), false)
}

func (this *State) GetOutputPath() string {
	return this.outputPath.Text()
}

func (this *State) GetOutputPathEditor() *widget.Editor {
	return &this.outputPath
}

func (this *State) Reset() {
	stateDto := dtos.NewDefaultEditorStateDto()
	this.stateDto = &stateDto
	this.currentPath = ""
	this.unsaved = false
	this.SetStatus("New settings file.", false)
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

	this.stateDto = loaded
	this.currentPath = path
	this.unsaved = false
	this.SetStatus("Loaded "+path, false)
}

func (this *State) Save() {
	if this.currentPath == "" {
		this.SaveAs(this.stateDto.TemplateName)
		return
	}

	if err := services.SaveSettingsFile(this.currentPath, this.stateDto); err != nil {
		this.SetStatus("Save failed: "+err.Error(), true)
		return
	}

	this.unsaved = false
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

	if err := services.SaveSettingsFile(path, this.stateDto); err != nil {
		this.SetStatus("Save failed: "+err.Error(), true)
		return
	}
	this.currentPath = path
	this.unsaved = false
	this.SetStatus("Saved "+path, false)
}

func (this *State) Generate() {
	dto, err := this.handler.GenerateTemplate(*this.stateDto)
	if err != nil {
		this.SetStatus(fmt.Sprintf("Generation failed: %v.", err), true)
		return
	}

	wasModified := this.connectionsModified
	this.connectionsModified = false
	this.lastTemplate = dto.Template
	zoneCount := 0
	connectionCount := 0
	if len(this.lastTemplate.Variants) > 0 {
		zoneCount = len(this.lastTemplate.Variants[0].Zones)
		connectionCount = len(this.lastTemplate.Variants[0].Connections)
	}
	status := fmt.Sprintf(
		"Generated '%s' — %d zones, %d connections.",
		this.lastTemplate.Name, zoneCount, connectionCount)
	if wasModified {
		status += " (Manual connection edits were replaced by regeneration.)"
	}
	this.SetStatus(status, false)
}

// SaveTemplate writes the most recently generated template as .rmg.json.
func (this *State) SaveTemplate() {
	savedPath, err := this.handler.SaveTemplate(dtos.TemplateDto{
		Template:   this.GetLastTemplate(),
		OutputPath: strings.TrimSpace(this.outputPath.Text()),
	})
	if err != nil {
		this.SetStatus(fmt.Sprintf("Save failed: %v", err), true)
		return
	}

	this.SetStatus("Saved template to "+savedPath, false)
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

func (this *State) UpdateState(updateFunc func(*dtos.EditorStateDto)) {
	// TODO: add validator for state updates, e.g. to prevent invalid map sizes or player counts
	updateFunc(this.stateDto)
}

func (this *State) SetStatus(msg string, isErr bool) {
	this.statusMsg = msg
	this.statusErr = isErr
}
