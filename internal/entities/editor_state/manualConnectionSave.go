package editor_state

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
)

type ManualConnectionSave struct {
	Connection  template_entity.Connection `json:"connection"`
	IsUserAdded bool                       `json:"isUserAdded,omitempty"`
}
