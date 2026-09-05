package mappers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

type IMandatoryContentItemMapper interface {
	FromRows(rows []editor_state_model.ZoneContentRow) []template_model.MandatoryContentItem
}
