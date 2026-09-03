package models

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

type EditorState struct {
	validationHandler handler_interfaces.IStateValidationHandler
	current           *editor_state_model.EditorState
	previous          *editor_state_model.EditorState
	next              *editor_state_model.EditorState
}

func NewEditorState(validationHandler handler_interfaces.IStateValidationHandler) *EditorState {
	state := &EditorState{validationHandler: validationHandler}
	state.ResetState()
	return state
}

func (this *EditorState) ResetState() {
	this.OverrideState(editor_state_model.NewDefaultEditorStateModel())
}

func (this *EditorState) OverrideState(state editor_state_model.EditorState) {
	storedState := state.Clone()
	this.previous = nil
	this.current = &storedState
	this.next = nil
}

func (this *EditorState) GetCurrentState() editor_state_model.EditorState {
	return this.current.Clone()
}

func (this *EditorState) GetTemplateName() string { return this.current.TemplateName }

func (this *EditorState) GetMapSize() int { return this.current.MapSize }

func (this *EditorState) GetTopology() config.MapTopology { return this.current.Topology }

func (this *EditorState) GetExperimentalMapSizes() bool { return this.current.ExperimentalMapSizes }

func (this *EditorState) UpdateCurrentState(updateFunc func(state *editor_state_model.EditorState)) {
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

func (this *EditorState) GetPreviousState() *editor_state_model.EditorState {
	if this.previous == nil {
		return nil
	}

	previousState := this.previous.Clone()
	return &previousState
}

func (this *EditorState) GetNextState() *editor_state_model.EditorState {
	if this.next == nil {
		return nil
	}

	nextState := this.next.Clone()
	return &nextState
}

func (this *EditorState) ResetNextState() { this.next = nil }

func (this *EditorState) SetNextState(state editor_state_model.EditorState) {
	storedState := state.Clone()
	this.next = &storedState
}

func (this *EditorState) HasNextState() bool { return this.next != nil }

func (this *EditorState) WasStateChanged() bool {
	return this.HasPreviousState() && !this.previous.EqualsIgnoringManualEdits(this.current)
}

func (this *EditorState) HasManualEdits() bool { return this.current.HasManualEdits() }

func (this *EditorState) SetManualEdits(zones []template_model.Zone, connections []entities.Connection) {
	this.current.ManualZones = editor_state_model.ToManualZoneSaves(zones)
	this.current.ManualConnections = editor_state_model.ToManualConnectionSaves(connections)
}

// ClearManualEdits drops the manual snapshot, used when a layout-defining
// option change invalidates the hand-made layout.
// This does not clear manual edits applied to the template,
// effectively making it useless if you want to reset the entire layout.
func (this *EditorState) ClearManualEdits() {
	this.current.ManualZones = nil
	this.current.ManualConnections = nil
}

func (this *EditorState) GetManualZones() []template_model.Zone {
	return editor_state_model.FromManualZoneSaves(this.current.ManualZones)
}

func (this *EditorState) GetManualConnections() []entities.Connection {
	return editor_state_model.FromManualConnectionSaves(this.current.ManualConnections)
}
