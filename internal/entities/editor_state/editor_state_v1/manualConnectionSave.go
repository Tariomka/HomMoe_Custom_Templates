package editor_state_v1

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
)

// ManualConnectionSave is the frozen v1 snapshot - see editorState.go. Never edit.
type ManualConnectionSave struct {
	Connection  entities.Connection `json:"connection"`
	IsUserAdded bool                `json:"isUserAdded,omitempty"`
}
