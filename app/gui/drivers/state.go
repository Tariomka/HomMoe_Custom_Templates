package drivers

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
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

	// lastGeneratedState is a copy of the editor state used for the most recent
	// successful generation. Auto-regeneration compares the live state against
	// it to detect option changes and to decide whether manual zone edits can
	// be reapplied.
	lastGeneratedState *dtos.EditorStateDto

	// manualZones / manualConnections hold the latest zones and connections
	// applied from the manual zone editor. When a regeneration leaves the
	// layout-defining options unchanged, these are reapplied so hand edits
	// survive option tweaks (e.g. changing a castle count).
	manualZones       []entities.Zone
	manualConnections []entities.Connection
	hasManualEdits    bool

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
	this.rememberManualEdits(zones, connections)
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
	this.clearGeneratedState()
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
	this.clearGeneratedState()
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
	this.applyGeneratedTemplate(dto.Template)
	this.discardManualEdits()
	zoneCount, connectionCount := this.lastTemplateZoneAndConnectionCount()
	status := fmt.Sprintf(
		"Generated '%s' — %d zones, %d connections.",
		this.lastTemplate.Name, zoneCount, connectionCount)
	if wasModified {
		status += " (Manual connection edits were replaced by regeneration.)"
	}
	this.SetStatus(status, false)
}

// AutoRegenerate regenerates the template whenever the live editor state has
// changed since the last generation. Manual zone edits are reapplied when the
// layout-defining options (player/zone counts, topology and connection
// settings) are unchanged, so tweaks such as castle counts keep hand edits.
// It is a no-op when nothing changed, keeping it cheap to call every frame.
func (this *State) AutoRegenerate() {
	if this.lastGeneratedState != nil && reflect.DeepEqual(*this.lastGeneratedState, *this.stateDto) {
		return
	}

	reapplyManual := this.lastGeneratedState != nil &&
		this.hasManualEdits &&
		!layoutDefiningOptionsChanged(this.lastGeneratedState, this.stateDto)

	dto, err := this.handler.GenerateTemplate(*this.stateDto)
	if err != nil {
		// Keep the previous template; only refresh the snapshot so we do not
		// spin trying (and failing) to regenerate every frame for the same
		// invalid state.
		this.snapshotGeneratedState()
		this.SetStatus(fmt.Sprintf("Auto-generation failed: %v.", err), true)
		return
	}

	this.applyGeneratedTemplate(dto.Template)
	if reapplyManual {
		this.reapplyManualEdits()
	} else {
		this.discardManualEdits()
	}

	zoneCount, connectionCount := this.lastTemplateZoneAndConnectionCount()
	notice := fmt.Sprintf(
		"Template regenerated with latest changes — %d zones, %d connections.",
		zoneCount, connectionCount)
	if reapplyManual {
		notice += " (Manual zone edits reapplied.)"
	}
	this.SetStatus(notice, false)
}

// applyGeneratedTemplate stores a freshly generated template as the live one
// and records the editor state that produced it.
func (this *State) applyGeneratedTemplate(template *entities.RmgTemplate) {
	this.connectionsModified = false
	this.lastTemplate = template
	this.snapshotGeneratedState()
}

// snapshotGeneratedState copies the current editor state so future changes can
// be detected by AutoRegenerate.
func (this *State) snapshotGeneratedState() {
	snapshot := *this.stateDto
	this.lastGeneratedState = &snapshot
}

// clearGeneratedState forgets the last generation and any manual edits, used
// when a brand-new or loaded settings file replaces the current one.
func (this *State) clearGeneratedState() {
	this.lastGeneratedState = nil
	this.lastTemplate = nil
	this.connectionsModified = false
	this.discardManualEdits()
}

// rememberManualEdits stores copies of the latest manually edited zones and
// connections so they can be reapplied after a non-structural regeneration.
func (this *State) rememberManualEdits(zones []entities.Zone, connections []entities.Connection) {
	this.manualZones = append([]entities.Zone(nil), zones...)
	this.manualConnections = append([]entities.Connection(nil), connections...)
	this.hasManualEdits = true
}

// discardManualEdits drops any stored manual edits.
func (this *State) discardManualEdits() {
	this.manualZones = nil
	this.manualConnections = nil
	this.hasManualEdits = false
}

// reapplyManualEdits pushes the stored manual zones and connections back onto
// the freshly generated template.
func (this *State) reapplyManualEdits() {
	if !this.hasManualEdits || this.lastTemplate == nil || len(this.lastTemplate.Variants) == 0 {
		return
	}
	dto, _ := this.handler.UpdateTemplate(dtos.TemplateUpdateDto{
		Template:    this.lastTemplate,
		Zones:       append([]entities.Zone(nil), this.manualZones...),
		Connections: append([]entities.Connection(nil), this.manualConnections...),
	})
	this.lastTemplate = dto.Template
	this.connectionsModified = true
}

func (this *State) lastTemplateZoneAndConnectionCount() (zoneCount, connectionCount int) {
	if this.lastTemplate != nil && len(this.lastTemplate.Variants) > 0 {
		zoneCount = len(this.lastTemplate.Variants[0].Zones)
		connectionCount = len(this.lastTemplate.Variants[0].Connections)
	}
	return zoneCount, connectionCount
}

// layoutDefiningOptionsChanged reports whether any option that changes the set
// of zones or the connection graph differs between two editor states. When
// these are unchanged, manual zone edits remain valid and can be reapplied.
func layoutDefiningOptionsChanged(previous, current *dtos.EditorStateDto) bool {
	return previous.PlayerCount != current.PlayerCount ||
		previous.Topology != current.Topology ||
		previous.GenerateRoads != current.GenerateRoads ||
		previous.RandomPortals != current.RandomPortals ||
		previous.NoDirectPlayerConn != current.NoDirectPlayerConn ||
		previous.MaxPortalConnections != current.MaxPortalConnections ||
		previous.MinNeutralZonesBetweenPlayers != current.MinNeutralZonesBetweenPlayers ||
		previous.SpawnRemoteFootholds != current.SpawnRemoteFootholds ||
		zoneCountOptionsChanged(previous, current)
}

// zoneCountOptionsChanged reports whether the number of neutral zones differs
// between two editor states.
func zoneCountOptionsChanged(previous, current *dtos.EditorStateDto) bool {
	return previous.AdvancedMode != current.AdvancedMode ||
		previous.NeutralZoneCount != current.NeutralZoneCount ||
		previous.NeutralLowNoCastleCount != current.NeutralLowNoCastleCount ||
		previous.NeutralLowCastleCount != current.NeutralLowCastleCount ||
		previous.NeutralMediumNoCastleCount != current.NeutralMediumNoCastleCount ||
		previous.NeutralMediumCastleCount != current.NeutralMediumCastleCount ||
		previous.NeutralHighNoCastleCount != current.NeutralHighNoCastleCount ||
		previous.NeutralHighCastleCount != current.NeutralHighCastleCount
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
