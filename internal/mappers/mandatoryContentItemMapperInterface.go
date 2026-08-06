package mappers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

type IMandatoryContentItemMapper interface {
	FromRows(rows []models.ZoneContentRowSave) []entities.MandatoryContentItem
}
