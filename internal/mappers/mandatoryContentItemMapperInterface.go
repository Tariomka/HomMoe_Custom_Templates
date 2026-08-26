package mappers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
)

type IMandatoryContentItemMapper interface {
	FromRows(rows []editor_state_model.ZoneContentRow) []entities.MandatoryContentItem
}
