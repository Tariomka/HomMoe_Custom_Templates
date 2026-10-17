package editor_state_dto

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
)

type EditorStateValidationDto struct {
	State    editor_state_model.EditorState
	Warnings []string
	// LoadWarnings report what migrating an older file dropped or changed.
	LoadWarnings []string
	// Rejection is set when the state cannot be accepted at all; State is then not usable.
	Rejection error
}
