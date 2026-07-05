package models

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
)

type EditorState struct {
	current  *dtos.EditorStateDto
	previous *dtos.EditorStateDto
	next     *dtos.EditorStateDto
}

func NewEditorState() *EditorState {
	state := new(EditorState)
	state.ResetState()
	return state
}

func (this *EditorState) ResetState() { this.OverrideState(dtos.NewDefaultEditorStateDto()) }

func (this *EditorState) OverrideState(state dtos.EditorStateDto) {
	this.previous = nil
	this.current = &state
	this.next = nil
}

func (this *EditorState) GetCurrentState() dtos.EditorStateDto {
	return *this.current
}

func (this *EditorState) UpdateCurrentState(updateFunc func(state *dtos.EditorStateDto)) {
	// TODO: add validator for state updates, e.g. to prevent invalid map sizes or player counts
	// this.SnapshotCurrentState()

	updateFunc(this.current)
	if this.current.AdvancedMode {
		this.current.NeutralZoneCount = 0
	} else {
		this.current.NeutralLowNoCastleCount = 0
		this.current.NeutralLowCastleCount = 0
		this.current.NeutralMediumNoCastleCount = 0
		this.current.NeutralMediumCastleCount = 0
		this.current.NeutralHighNoCastleCount = 0
		this.current.NeutralHighCastleCount = 0
	}
}

func (this *EditorState) SnapshotCurrentState() {
	previousState := *this.current
	this.previous = &previousState
	// this.next = nil
}

func (this *EditorState) ResetPreviousState() { this.previous = nil }

func (this *EditorState) HasPreviousState() bool { return this.previous != nil }

func (this *EditorState) ResetNextState() { this.next = nil }

func (this *EditorState) ResetNextStateIfLayoutChanged() bool {
	if this.WasLayoutChanged() {
		this.ResetNextState()
		return true
	}

	return false
}

func (this *EditorState) ResetNextStateIfStateWasNotChanged() bool {
	if this.WasStateUnchanged() {
		this.ResetNextState()
		return true
	}

	return false
}

func (this *EditorState) SetNextFromCurrentIfStateIsBeingUpdated() bool {
	if !this.HasNextState() || this.HasPendingChanges() {
		this.SetNextState(*this.current)
		return true
	}

	return false
}

func (this *EditorState) SetNextState(state dtos.EditorStateDto) { this.next = &state }

func (this *EditorState) HasNextState() bool { return this.next != nil }

func (this *EditorState) WasStateChanged() bool {
	return this.HasPreviousState() && !this.previous.EqualsIgnoringManualEdits(this.current)
}

func (this *EditorState) WasStateUnchanged() bool {
	return this.HasPreviousState() && this.previous.EqualsIgnoringManualEdits(this.current)
}

func (this *EditorState) WasLayoutChanged() bool {
	return this.previous.LayoutDefiningOptionsChanged(this.current)
}

func (this *EditorState) WasLayoutUnchanged() bool {
	return this.HasPreviousState() && !this.previous.LayoutDefiningOptionsChanged(this.current)
}

func (this *EditorState) HasPendingChanges() bool {
	return this.HasNextState() && !this.next.EqualsIgnoringManualEdits(this.current)
}

// ── Manual zone editor edits ───────────────────────────────────────────
// The manual zones/connections snapshot lives inside the current
// EditorStateDto, so it is saved to and loaded from the .gen.json file with
// the rest of the editor state and needs no separate bookkeeping.

func (this *EditorState) HasManualEdits() bool { return this.current.HasManualEdits() }

// SetManualEdits stores the manual zone editor's result as the authoritative
// snapshot. Manual fields are ignored by the state-equality checks, so this
// never triggers an automatic regeneration by itself.
func (this *EditorState) SetManualEdits(zones []entities.Zone, connections []entities.Connection) {
	this.current.ManualZones = editor_state_dto.ToManualZoneSaves(zones)
	this.current.ManualConnections = editor_state_dto.ToManualConnectionSaves(connections)
}

// ClearManualEdits drops the manual snapshot, used when a layout-defining
// option change invalidates the hand-made layout.
func (this *EditorState) ClearManualEdits() {
	this.current.ManualZones = nil
	this.current.ManualConnections = nil
}

func (this *EditorState) GetManualZones() []entities.Zone {
	return editor_state_dto.FromManualZoneSaves(this.current.ManualZones)
}

func (this *EditorState) GetManualConnections() []entities.Connection {
	return editor_state_dto.FromManualConnectionSaves(this.current.ManualConnections)
}

// ShouldReapplyManualEdits reports whether the stored manual edits are still
// valid for the current state: they exist and no layout-defining option
// changed since the last generation. Right after a load there is no previous
// state; the stored edits are then trusted because they were saved against
// exactly this state.
func (this *EditorState) ShouldReapplyManualEdits() bool {
	if !this.HasManualEdits() {
		return false
	}
	return !this.HasPreviousState() || !this.WasLayoutChanged()
}

// CastleSettingsChangedSinceGeneration reports which castle-count options
// changed since the last generated state - the only option changes that are
// pushed into the manual snapshot. Zero-value when nothing was generated yet.
func (this *EditorState) CastleSettingsChangedSinceGeneration() editor_state_dto.CastleSettingChanges {
	if !this.HasPreviousState() {
		return editor_state_dto.CastleSettingChanges{}
	}
	return this.previous.DiffCastleSettings(this.current)
}
