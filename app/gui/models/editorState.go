package models

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
)

type EditorState struct {
	validationHandler handler_interfaces.IStateValidationHandler
	current           *editor_state_dto.EditorStateDto
	previous          *editor_state_dto.EditorStateDto
	next              *editor_state_dto.EditorStateDto
}

func NewEditorState(validationHandler handler_interfaces.IStateValidationHandler) *EditorState {
	state := &EditorState{validationHandler: validationHandler}
	state.ResetState()
	return state
}

func (this *EditorState) ResetState() {
	this.OverrideState(editor_state_dto.NewDefaultEditorStateDto())
}

func (this *EditorState) OverrideState(state editor_state_dto.EditorStateDto) {
	storedState := state.Clone()
	this.previous = nil
	this.current = &storedState
	this.next = nil
}

func (this *EditorState) GetCurrentState() editor_state_dto.EditorStateDto {
	return this.current.Clone()
}

func (this *EditorState) GetTemplateName() string { return this.current.TemplateName }

func (this *EditorState) GetMapSize() int { return this.current.MapSize }

func (this *EditorState) GetTopology() config.MapTopology { return this.current.Topology }

func (this *EditorState) GetExperimentalMapSizes() bool { return this.current.ExperimentalMapSizes }

func (this *EditorState) UpdateCurrentState(updateFunc func(state *editor_state_dto.EditorStateDto)) {
	updatedState := this.current.Clone()
	updateFunc(&updatedState)
	validation := this.validationHandler.ValidateEditorState(updatedState, true)
	this.current = &validation.State
}

func (this *EditorState) SnapshotCurrentState() {
	previousState := this.current.Clone()
	this.previous = &previousState
	this.next = nil
}

func (this *EditorState) HasPreviousState() bool { return this.previous != nil }

func (this *EditorState) GetPreviousState() *editor_state_dto.EditorStateDto {
	if this.previous == nil {
		return nil
	}

	previousState := this.previous.Clone()
	return &previousState
}

func (this *EditorState) GetNextState() *editor_state_dto.EditorStateDto {
	if this.next == nil {
		return nil
	}

	nextState := this.next.Clone()
	return &nextState
}

func (this *EditorState) ResetNextState() { this.next = nil }

func (this *EditorState) SetNextState(state editor_state_dto.EditorStateDto) {
	storedState := state.Clone()
	this.next = &storedState
}

func (this *EditorState) HasNextState() bool { return this.next != nil }

func (this *EditorState) WasStateChanged() bool {
	return this.HasPreviousState() && !this.previous.EqualsIgnoringManualEdits(this.current)
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
