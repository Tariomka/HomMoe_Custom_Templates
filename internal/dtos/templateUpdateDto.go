package dtos

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

type TemplateUpdateDto struct {
	Template    *template_model.Template
	Zones       []template_model.Zone
	Connections []template_model.Connection

	// EditorState, when supplied, lets UpdateTemplate rebuild the template's
	// mandatory content from the final (manually edited) zones so a zone whose
	// quality was changed in the editor gets the content of its new tier.
	EditorState *editor_state_dto.EditorStateDto
}
