package handlers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/bonuses"
)

type bonusHandler struct {
	bonusEntryService bonuses.IBonusEntryService
}

func NewBonusHandler(bonusEntryService bonuses.IBonusEntryService) handler_interfaces.IBonusHandler {
	return &bonusHandler{bonusEntryService: bonusEntryService}
}

func (this *bonusHandler) DescribeExistingBonuses(existing []config.BonusEntry) dtos.ExistingBonusesDto {
	return this.bonusEntryService.DescribeExistingBonuses(existing)
}

func (this *bonusHandler) BuildBonusEntries(
	request dtos.BonusCompositionRequestDto) dtos.BonusCompositionResultDto {
	return this.bonusEntryService.BuildBonusEntries(request)
}

func (this *bonusHandler) FilterNewBonusEntries(
	entries []config.BonusEntry,
	existingKeys map[string]bool) []config.BonusEntry {
	return this.bonusEntryService.FilterNewBonusEntries(entries, existingKeys)
}

func (this *bonusHandler) GetSpellCountLabel(count int) string {
	return this.bonusEntryService.GetSpellCountLabel(count)
}
