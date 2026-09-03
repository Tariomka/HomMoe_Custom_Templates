package dtos

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

// ZoneEditorZonesDto carries a manual zone editor's zone and connection lists
// in either direction: out of the editor when the user applies, and back into
// an open editor when the template is reverted to a freshly generated base.
type ZoneEditorZonesDto struct {
	Zones       []template_model.Zone
	Connections []entities.Connection
	// RevertToBase reports that the editing session reverted to a freshly
	// generated layout. It is only meaningful on the apply direction.
	RevertToBase bool
}
