package handler_interfaces

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
)

type IBonusHandler interface {
	DescribeExistingBonuses(existing []config.BonusEntry) dtos.ExistingBonusesDto
	BuildBonusEntries(request dtos.BonusCompositionRequestDto) dtos.BonusCompositionResultDto
	FilterNewBonusEntries(entries []config.BonusEntry, existingKeys map[string]bool) []config.BonusEntry
	GetSpellCountLabel(count int) string
}
