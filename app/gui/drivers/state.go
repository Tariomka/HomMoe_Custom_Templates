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

	templateDir, err := helpers.FindOldenEraTemplatesDir(false)
	if templateDir == "" {
		if err != nil {
			state.SetStatus(fmt.Sprintf("Failed to find game template directory: %v", err), true)
		} else {
			state.SetStatus("Failed to find game template directory, using fallback directory.", false)
		}
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
	dto, err := this.handler.UpdateTemplate(dtos.TemplateUpdateDto{
		Template:    this.lastTemplate,
		Zones:       zones,
		Connections: connections,
	})
	this.lastTemplate = dto.Template
	this.connectionsModified = true
	if err != nil {
		status := fmt.Sprintf(
			"Applied %d zones and %d connections. \u26a0 Error: %v; fix before export.",
			len(zones), len(connections), err)
		this.SetStatus(status, true)
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
		this.SetStatus(fmt.Sprintf("Open dialog failed: %v.", err), true)
		return
	}

	if path == "" {
		return
	}

	dto, err := this.handler.LoadState(path)
	if err != nil {
		this.SetStatus(fmt.Sprintf("Load failed: %v.", err), true)
		return
	}

	this.stateDto = dto
	this.currentPath = path
	this.unsaved = false
	this.SetStatus("Loaded "+path, false)
}

func (this *State) Save() {
	if this.currentPath == "" {
		this.SaveAs(this.stateDto.TemplateName)
		return
	}

	if _, err := this.handler.SaveState(dtos.EditorStateSaveDto{
		State:      this.stateDto,
		OutputPath: this.currentPath,
	}); err != nil {
		this.SetStatus(fmt.Sprintf("Save failed: %v.", err), true)
		return
	}

	this.unsaved = false
	this.SetStatus("Saved "+this.currentPath, false)
}

func (this *State) SaveAs(templateName string) {
	defaultName := helpers.SanitizeFilename(strings.TrimSpace(templateName)) + ".gen.json"
	path, err := utils.PickSaveFile("Save settings as", "Settings (*.gen.json)|*.gen.json", this.SuggestDirectory(), defaultName)
	if err != nil {
		this.SetStatus(fmt.Sprintf("Save dialog failed: %v.", err), true)
		return
	}

	if path == "" {
		return
	}

	if _, err := this.handler.SaveState(dtos.EditorStateSaveDto{
		State:      this.stateDto,
		OutputPath: path,
	}); err != nil {
		this.SetStatus(fmt.Sprintf("Save failed: %v.", err), true)
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
	savedPath, err := this.handler.SaveTemplate(dtos.TemplateSaveDto{
		Template:   this.GetLastTemplate(),
		Topology:   this.stateDto.Topology,
		OutputPath: strings.TrimSpace(this.outputPath.Text()),
	})
	if err != nil && savedPath == "" {
		this.SetStatus(fmt.Sprintf("Save failed: %v.", err), true)
		return
	} else if err != nil {
		this.SetStatus(
			fmt.Sprintf("Saved template to %s, but failed to write preview PNG with error: %v.", savedPath, err),
			true)
		return
	}

	this.SetStatus("Saved template to "+savedPath, false)
}

// PickOutputDir presents a folder picker for the template output directory.
func (this *State) PickOutputDir() {
	cur := strings.TrimSpace(this.outputPath.Text())
	dir, err := utils.PickFolder("Select output directory", cur)
	if err != nil {
		this.SetStatus(fmt.Sprintf("Folder dialog failed: %v.", err), true)
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
