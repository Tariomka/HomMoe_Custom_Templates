package dtos

import "github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"

// ZoneEditorConnectionTypeRequestDto carries a pending connection-type change
// from the manual editor, together with the current between-zone road setting.
type ZoneEditorConnectionTypeRequestDto struct {
	Connection     template_model.Connection
	ConnectionType string
	GenerateRoads  bool
}
