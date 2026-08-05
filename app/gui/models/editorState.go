package models

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
)

type EditorState struct {
	validationHandler handler_interfaces.IStateValidationHandler
	current           *dtos.EditorStateDto
	previous          *dtos.EditorStateDto
	next              *dtos.EditorStateDto
}

func NewEditorState(validationHandler handler_interfaces.IStateValidationHandler) *EditorState {
	state := &EditorState{validationHandler: validationHandler}
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
	updateFunc(this.current)
	validation := this.validationHandler.ValidateEditorState(*this.current, true)
	*this.current = validation.State
}

func (this *EditorState) SnapshotCurrentState() {
	previousState := *this.current
	this.previous = &previousState
	this.next = nil
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
	return this.HasPreviousState() && this.previous.LayoutDefiningOptionsChanged(this.current)
}

func (this *EditorState) WasLayoutUnchanged() bool {
	return this.HasPreviousState() && !this.previous.LayoutDefiningOptionsChanged(this.current)
}

func (this *EditorState) HasPendingChanges() bool {
	return this.HasNextState() && !this.next.EqualsIgnoringManualEdits(this.current)
}

func (this *EditorState) HasManualEdits() bool { return this.current.HasManualEdits() }

func (this *EditorState) SetManualEdits(zones []entities.Zone, connections []entities.Connection) {
	this.current.ManualZones = editor_state_dto.ToManualZoneSaves(zones)
	this.current.ManualConnections = editor_state_dto.ToManualConnectionSaves(connections)
}

// ClearManualEdits drops the manual snapshot, used when a layout-defining
// option change invalidates the hand-made layout.
// This does not clear manual edits applied to the template,
// effectively making it useless if you want to reset the entire layout.
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

func (this *EditorState) ShouldReapplyManualEdits() bool {
	if !this.HasManualEdits() {
		return false
	}

	return !this.WasLayoutChanged()
}

func (this *EditorState) CastleSettingsChangedSinceGeneration() editor_state_dto.CastleSettingChanges {
	if !this.HasPreviousState() {
		return editor_state_dto.CastleSettingChanges{}
	}

	return this.previous.DiffCastleSettings(this.current)
}
