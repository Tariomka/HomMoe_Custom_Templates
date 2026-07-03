package drivers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/dialogs"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/models"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/widgets"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
)

const (
	autoRegenDebounce   = 300 * time.Millisecond
	configFileExtension = ".gen.json"
)

type State struct {
	handler *handlers.GUIHandler
	mapper  *mappers.GeneratorConfigMapper

	innerState *models.EditorState

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

	// manualZones / manualConnections hold the latest zones and connections
	// applied from the manual zone editor. When a regeneration leaves the
	// layout-defining options unchanged, these are reapplied so hand edits
	// survive option tweaks (e.g. changing a castle count).
	manualZones       []entities.Zone
	manualConnections []entities.Connection
	hasManualEdits    bool

	confirmExit bool

	// pendingManualReapply requests that the stored manual edits be reapplied on
	// the next regeneration even though there is no prior generated state to
	// compare against. It is set when a saved state with manual edits is loaded
	// so the hand-made layout is restored once the loaded template is
	// regenerated.
	pendingManualReapply bool

	applyNextStateAt time.Time

	// dialogs renders modal dialogs (rule editors, pickers, the connection
	// editor) over the main UI.
	dialogs *DialogHost
}

func NewUIState() *State {
	state := &State{
		handler:    handlers.NewGuiHandler(),
		mapper:     mappers.NewConfigMapper(),
		innerState: models.NewEditorState(),
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

func (this *State) GetStatus() (msg string, isErr bool) { return this.statusMsg, this.statusErr }

// Dialogs returns the modal host used to open and render dialogs.
func (this *State) Dialogs() *DialogHost { return this.dialogs }

func (this *State) GetStateData() dtos.EditorStateDto { return this.innerState.GetCurrentState() }

func (this *State) GetGeneratorConfig() *config.GeneratorConfig {
	return this.mapper.FromEditorState(this.innerState.GetCurrentState())
}

func (this *State) GetCurrentPath() string { return this.currentPath }

func (this *State) IsUnsaved() bool { return this.unsaved }

func (this *State) GetLastTemplate() *entities.RmgTemplate { return this.lastTemplate }

// ApplyEditedZones writes zones and connections edited in the manual zone
// editor back into the live template and flags that manual edits now exist.
func (this *State) ApplyEditedZones(zones []entities.Zone, connections []entities.Connection) {
	if !this.hasTemplateVariants() {
		return
	}

	this.handleUpdateTemplate(zones, connections)
	this.rememberManualEdits(zones, connections)
}

func (this *State) GetOutputPath() string { return this.outputPath.Text() }

func (this *State) GetOutputPathWidget(theme *material.Theme) layout.Widget {
	return widgets.NewTextboxWidget(theme, &this.outputPath, "Choose folder", true)
}

func (this *State) Reset() {
	this.innerState.ResetState()
	this.currentPath = ""
	this.unsaved = false
	this.clearGeneratedState()
	this.SetStatus("New settings file.", false)
}

func (this *State) Load() {
	dir, err := os.Getwd()
	if err != nil {
		dir = this.suggestDirectory()
	}
	this.dialogs.Open(dialogs.NewOpenFileDialog(dir, []string{".gen.json"}, func(path string) {
		this.handleLoadState(path)
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
	dir, err := os.Getwd()
	if err != nil {
		dir = this.suggestDirectory()
	}
	defaultName := helpers.SanitizeFilename(strings.TrimSpace(templateName)) + ".gen.json"
	this.dialogs.Open(dialogs.NewSaveFileDialog(dir, defaultName, func(path string) {
		this.handleSaveState(path)
		this.currentPath = path
	}))
}

func (this *State) Exit() {
	// Give State an onExit func() callback (injected from program.go).
	// In program.go, pass func() { window.Perform(system.ActionClose) } so the normal app.DestroyEvent path runs.
	if this.unsaved && !this.confirmExit {
		this.SetStatus("Unsaved changes exist - save first or press Exit again.", true)
		this.confirmExit = true // reset if not selected right after
		return
	}
	// this.onExit()

	os.Exit(0)
}

func (this *State) Generate() { this.handleGenerateTemplate(false) }

// SaveTemplate writes the most recently generated template as .rmg.json.
func (this *State) SaveTemplate() { this.handleSaveTemplate() }

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

func (this *State) UpdateState(updateFunc func(*dtos.EditorStateDto)) {
	this.innerState.UpdateCurrentState(updateFunc)
	if this.innerState.WasStateChanged() {
		this.unsaved = true
	}
}

func (this *State) SetStatus(msg string, isErr bool) {
	this.statusMsg = msg
	this.statusErr = isErr
}

// AutoRegenerate regenerates the template when the live editor state has
// changed since the last generation.
//
// Preview-affecting changes (player/zone counts, topology and connection
// settings) regenerate immediately so the live preview tracks the control.
// All other changes (non-structural sliders and the template name) are
// debounced and only regenerate once editing has paused for autoRegenDebounce,
// avoiding a regeneration on every frame while a slider is dragged.
//
// now is the current frame time. It returns the time at which the caller
// should request another frame and whether such a redraw must be scheduled,
// used to wake the UI back up once the debounce window elapses without further
// input.
func (this *State) AutoRegenerate(now time.Time) (redrawAt time.Time, scheduleRedraw bool) {
	// Nothing changed since the last generation → cancel any pending debounce.
	if this.innerState.ResetNextStateIfStateWasNotChanged() {
		return time.Time{}, false
	}

	// First generation: populate the preview immediately on startup.
	if !this.innerState.HasPreviousState() {
		this.handleGenerateTemplate(true)
		return time.Time{}, false
	}

	// Preview-affecting changes regenerate immediately so the preview follows
	// the control live.
	if this.innerState.ResetNextStateIfLayoutChanged() {
		this.handleGenerateTemplate(true)
		return time.Time{}, false
	}

	// Non-preview change: (re)arm the debounce timer whenever the state is
	// still moving, and ask to be woken up when the timer is due.
	if this.innerState.SetNextFromCurrentIfStateIsBeingUpdated() {
		this.applyNextStateAt = now.Add(autoRegenDebounce)
		return this.applyNextStateAt, true
	}

	// State has been stable since the last frame; keep waiting until due.
	if now.Before(this.applyNextStateAt) {
		return this.applyNextStateAt, true
	}

	// Editing paused long enough -> regenerate now.
	this.innerState.ResetNextState()
	this.handleGenerateTemplate(true)
	return time.Time{}, false
}

func (this *State) handleSaveState(path string) {
	this.syncManualEditsToDto()
	if _, err := this.handler.SaveState(dtos.EditorStateSaveDto{
		State:      new(this.innerState.GetCurrentState()),
		OutputPath: path,
	}); err != nil {
		this.SetStatus(fmt.Sprintf("Save failed: %v.", err), true)
		return
	}

	this.unsaved = false
	this.SetStatus("Saved "+path, false)
}

func (this *State) handleLoadState(path string) {
	dto, err := this.handler.LoadState(path)
	if err != nil {
		this.SetStatus(fmt.Sprintf("Load failed: %v.", err), true)
		return
	}

	this.innerState.OverrideState(*dto)
	this.currentPath = path
	this.unsaved = false
	this.clearGeneratedState()
	this.restoreManualEdits(dto)
	this.SetStatus("Loaded "+path, false)
}

func (this *State) handleSaveTemplate() {
	savedPath, err := this.handler.SaveTemplate(dtos.TemplateSaveDto{
		Template:   this.GetLastTemplate(),
		Topology:   this.innerState.GetCurrentState().Topology,
		OutputPath: strings.TrimSpace(this.outputPath.Text()),
	})

	if err == nil {
		this.SetStatus("Saved template to "+savedPath, false)
		return
	}

	if savedPath == "" {
		this.SetStatus(fmt.Sprintf("Save failed: %v.", err), true)
		return
	}

	this.SetStatus(
		fmt.Sprintf("Saved template to %s, but failed to write preview PNG with error: %v.", savedPath, err),
		true)
}

func (this *State) handleUpdateTemplate(zones []entities.Zone, connections []entities.Connection) {
	dto, err := this.handler.UpdateTemplate(dtos.TemplateUpdateDto{
		Template:    this.lastTemplate,
		Zones:       zones,
		Connections: connections,
		Config:      this.GetGeneratorConfig(),
	})
	this.lastTemplate = dto.Template
	this.connectionsModified = true

	if err != nil {
		this.SetStatus(
			fmt.Sprintf(
				"Applied %d zones and %d connections. ⚠ Error: %v; fix before export.",
				len(zones), len(connections), err),
			true)
		return
	}

	this.SetStatus(
		fmt.Sprintf("Applied %d zones and %d connections from the editor.", len(zones), len(connections)),
		false)
}

func (this *State) handleGenerateTemplate(createStateSnapshotOnFailure bool) {
	dto, err := this.handler.GenerateTemplate(this.innerState.GetCurrentState())
	if err != nil {
		this.SetStatus(fmt.Sprintf("Generation failed: %v.", err), true)
		if createStateSnapshotOnFailure {
			this.innerState.SnapshotCurrentState()
		}
		return
	}

	reapplyManual := this.shouldReapplyManualEdits()
	this.applyGeneratedTemplate(dto.Template)
	if reapplyManual && this.hasTemplateVariants() {
		this.handleUpdateTemplate(
			append([]entities.Zone(nil), this.manualZones...),
			append([]entities.Connection(nil), this.manualConnections...))
	} else if !reapplyManual {
		this.discardManualEdits()
	}
	this.pendingManualReapply = false

	zoneCount, connectionCount := this.lastTemplateZoneAndConnectionCount()
	status := fmt.Sprintf(
		"Template '%s' generated with latest changes - %d zones, %d connections.",
		this.lastTemplate.Name, zoneCount, connectionCount)
	if reapplyManual {
		status += " (Manual zone edits reapplied.)"
	}
	status += fmt.Sprintf("\n%s", time.Now().Format("15:04:05"))
	this.SetStatus(status, false)
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

// applyGeneratedTemplate stores a freshly generated template as the live one
// and records the editor state that produced it.
func (this *State) applyGeneratedTemplate(template *entities.RmgTemplate) {
	this.connectionsModified = false
	this.lastTemplate = template
	this.innerState.SnapshotCurrentState()
}

// clearGeneratedState forgets the last generation and any manual edits, used
// when a brand-new or loaded settings file replaces the current one.
func (this *State) clearGeneratedState() {
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
	this.pendingManualReapply = false
}

// syncManualEditsToDto mirrors the in-memory manual edits onto the persistent
// editor state so they are written to the .gen.json file. When no manual edits
// exist the persisted fields are cleared.
func (this *State) syncManualEditsToDto() {
	manualZones := []dtos.ManualZoneSave(nil)
	manualConnections := []dtos.ManualConnectionSave(nil)
	if this.hasManualEdits {
		manualZones = dtos.ToManualZoneSaves(this.manualZones)
		manualConnections = dtos.ToManualConnectionSaves(this.manualConnections)
	}
	this.innerState.UpdateCurrentState(func(state *dtos.EditorStateDto) {
		state.ManualZones = manualZones
		state.ManualConnections = manualConnections
	})
}

// restoreManualEdits loads any persisted manual edits from a freshly loaded
// editor state back into memory and arms a reapply so the hand-made layout is
// restored once the loaded template is regenerated.
func (this *State) restoreManualEdits(dto *dtos.EditorStateDto) {
	if dto == nil || !dto.HasManualEdits() {
		return
	}

	this.manualZones = dtos.FromManualZoneSaves(dto.ManualZones)
	this.manualConnections = dtos.FromManualConnectionSaves(dto.ManualConnections)
	this.hasManualEdits = true
	this.pendingManualReapply = true
}

func (this *State) shouldReapplyManualEdits() bool {
	return this.hasManualEdits && (this.pendingManualReapply || this.innerState.WasLayoutUnchanged())
}

func (this *State) hasTemplateVariants() bool {
	return this.lastTemplate != nil && len(this.lastTemplate.Variants) > 0
}

func (this *State) lastTemplateZoneAndConnectionCount() (zoneCount, connectionCount int) {
	if this.hasTemplateVariants() {
		zoneCount = len(this.lastTemplate.Variants[0].Zones)
		connectionCount = len(this.lastTemplate.Variants[0].Connections)
	}
	return
}
