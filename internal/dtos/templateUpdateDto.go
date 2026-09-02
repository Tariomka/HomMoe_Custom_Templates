package dtos

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

type TemplateUpdateDto struct {
	Template *template_model.Template
	// Zones and Connections are still entities: the zone editor has not moved
	// onto template_model yet, which is phase 4 of this batch.
	Zones       []entities.Zone
	Connections []entities.Connection

	// EditorState, when supplied, lets UpdateTemplate rebuild the template's
	// mandatory content from the final (manually edited) zones so a zone whose
	// quality was changed in the editor gets the content of its new tier.
	EditorState *editor_state_dto.EditorStateDto
}
