package drivers

import (
	"fmt"
	"strings"
	"time"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
)

func (this *State) Generate() { this.handleGenerateTemplate(false) }

// SaveTemplate writes the most recently generated template as .rmg.json.
func (this *State) SaveTemplate() { this.handleSaveTemplate() }

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

func (this *State) handleGenerateTemplate(createStateSnapshotOnFailure bool) {
	dto, err := this.handler.GenerateTemplate(this.innerState.GetCurrentState())
	if err != nil {
		this.SetStatus(fmt.Sprintf("Generation failed: %v.", err), true)
		if createStateSnapshotOnFailure {
			this.innerState.SnapshotCurrentState()
		}
		return
	}

	// The reapply decision and the castle-option diff both compare against the
	// state of the LAST generation, so they must be taken before
	// applyGeneratedTemplate snapshots the current state.
	reapplyManual := this.innerState.ShouldReapplyManualEdits()
	castleChanges := this.innerState.CastleSettingsChangedSinceGeneration()
	this.applyGeneratedTemplate(dto.Template)
	if reapplyManual && this.hasTemplateVariants() {
		this.reapplyManualEdits(castleChanges)
	} else if !reapplyManual {
		this.innerState.ClearManualEdits()
	}

	zoneCount, connectionCount := this.lastTemplateZoneAndConnectionCount()
	status := fmt.Sprintf(
		"Template '%s' generated with latest changes - %d zones, %d connections.",
		this.lastTemplate.Name, zoneCount, connectionCount)
	if len(dto.Warnings) > 0 {
		status += fmt.Sprintf(" Note that %d warning(s) were found and fixed.", len(dto.Warnings))
	}
	if reapplyManual {
		status += " (Manual zone edits reapplied.)"
	}
	status += fmt.Sprintf("\n%s", time.Now().Format("15:04:05"))
	this.SetStatus(status, false)
}

// applyGeneratedTemplate stores a freshly generated template as the live one
// and records the editor state that produced it.
func (this *State) applyGeneratedTemplate(template *entities.RmgTemplate) {
	this.lastTemplate = template
	this.innerState.SnapshotCurrentState()
}

// clearGeneratedState forgets the last generated template, used when a
// brand-new or loaded settings file replaces the current one. Manual edits
// need no separate handling: they live inside the editor state itself, which
// the caller is replacing.
func (this *State) clearGeneratedState() {
	this.lastTemplate = nil
}

func (this *State) lastTemplateZoneAndConnectionCount() (zoneCount, connectionCount int) {
	if this.hasTemplateVariants() {
		zoneCount = len(this.lastTemplate.Variants[0].Zones)
		connectionCount = len(this.lastTemplate.Variants[0].Connections)
	}
	return zoneCount, connectionCount
}
