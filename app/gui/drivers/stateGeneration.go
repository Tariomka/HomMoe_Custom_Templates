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
// The decision itself belongs to the regeneration handler; this method only
// feeds it the current snapshots and frame time, applies the pending-snapshot
// mutation it asks for, and performs the regeneration.
//
// now is the current frame time. It returns the time at which the caller
// should request another frame and whether such a redraw must be scheduled,
// used to wake the UI back up once the debounce window elapses without further
// input.
func (this *State) AutoRegenerate(now time.Time) (redrawAt time.Time, scheduleRedraw bool) {
	decision := this.regeneration.DecideRegeneration(dtos.RegenerationDecisionRequestDto{
		Previous:      this.editorStateMapper.ToDtoPointer(this.innerState.GetPreviousState()),
		Current:       new(this.GetStateDto()),
		Next:          this.editorStateMapper.ToDtoPointer(this.innerState.GetNextState()),
		Now:           now,
		DebounceDueAt: this.applyNextStateAt,
	})

	switch decision.NextStateAction {
	case dtos.NextStateClear:
		this.innerState.ResetNextState()
	case dtos.NextStateSetFromCurrent:
		this.innerState.SetNextState(this.innerState.GetCurrentState())
	case dtos.NextStateLeave:
	}

	if decision.ScheduleRedraw {
		this.applyNextStateAt = decision.RedrawAt
	}

	if decision.Regenerate {
		this.handleGenerateTemplate(true)
	}

	return decision.RedrawAt, decision.ScheduleRedraw
}

func (this *State) handleSaveTemplate() {
	savedPath, err := this.handler.SaveTemplate(dtos.TemplateSaveDto{
		Template:   this.GetLastTemplate(),
		Topology:   this.innerState.GetTopology(),
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

// handleGenerateTemplate regenerates the template; on failure the previous
// template is left in place.
func (this *State) handleGenerateTemplate(createStateSnapshotOnFailure bool) {
	currentState := this.GetStateDto()
	dto, err := this.handler.GenerateTemplate(currentState)
	if err != nil {
		this.SetStatus(fmt.Sprintf("Generation failed: %v.", err), true)
		if createStateSnapshotOnFailure {
			this.innerState.SnapshotCurrentState()
		}
		return
	}

	// The decision compares against the state of the LAST generation, so it
	// must be taken before applyGeneratedTemplate snapshots the current state.
	manualEdits := this.regeneration.DecideManualEditReapplication(
		this.editorStateMapper.ToDtoPointer(this.innerState.GetPreviousState()), &currentState)
	this.applyGeneratedTemplate(dto.Template)
	if manualEdits.ReapplyWithCastleChanges != nil && this.hasTemplateVariants() {
		this.reapplyManualEdits(
			this.editorStateMapper.ToCastleSettingChangesModel(*manualEdits.ReapplyWithCastleChanges))
	} else if manualEdits.ReapplyWithCastleChanges == nil {
		this.innerState.ClearManualEdits()
	}

	zoneCount, connectionCount := this.lastTemplateZoneAndConnectionCount()
	status := fmt.Sprintf(
		"Template '%s' generated with latest changes - %d zones, %d connections.",
		this.lastTemplate.Name, zoneCount, connectionCount)
	if len(dto.Warnings) > 0 {
		status += fmt.Sprintf(" Note that %d warning(s) were found and fixed.", len(dto.Warnings))
	}
	if manualEdits.ReapplyWithCastleChanges != nil {
		status += " (Manual zone edits reapplied.)"
	}
	status += fmt.Sprintf("\n%s", time.Now().Format("15:04:05"))
	this.SetStatus(status, false)
}

// applyGeneratedTemplate stores a freshly generated template as the live one
// and records the editor state that produced it.
func (this *State) applyGeneratedTemplate(template *entities.RmgTemplate) {
	this.setLastTemplate(template)
	this.innerState.SnapshotCurrentState()
}

// clearGeneratedState forgets the last generated template, used when a
// brand-new or loaded settings file replaces the current one. Manual edits
// need no separate handling: they live inside the editor state itself, which
// the caller is replacing.
func (this *State) clearGeneratedState() {
	this.setLastTemplate(nil)
}

func (this *State) lastTemplateZoneAndConnectionCount() (zoneCount, connectionCount int) {
	if this.hasTemplateVariants() {
		zoneCount = len(this.lastTemplate.Variants[0].Zones)
		connectionCount = len(this.lastTemplate.Variants[0].Connections)
	}
	return zoneCount, connectionCount
}
