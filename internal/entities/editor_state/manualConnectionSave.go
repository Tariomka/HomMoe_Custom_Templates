package editor_state

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
)

// ManualConnectionSave persists a connection edited in the manual zone editor,
// capturing the runtime-only IsUserAdded flag that entities.Connection omits
// from JSON (json:"-").
type ManualConnectionSave struct {
	Connection  entities.Connection `json:"connection"`
	IsUserAdded bool                `json:"isUserAdded,omitempty"`
}
