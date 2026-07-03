package models

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
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

// func (this *EditorState) Undo() {
// 	if this.previous != nil {
// 		nextState := *this.current
// 		this.next = &nextState
// 		this.current = this.previous
// 		this.previous = nil
// 	}
// }

// func (this *EditorState) Redo() {
// 	if this.next != nil {
// 		previousState := *this.current
// 		this.previous = &previousState
// 		this.current = this.next
// 		this.next = nil
// 	}
// }
