package bonuses

import (
	"strconv"
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/config_helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
)

const (
	spellFreeParam    = "1"
	spellNotFreeParam = "0"
)

type BonusEntryService struct{}

func NewBonusEntryService() IBonusEntryService {
	return &BonusEntryService{}
}

func (this *BonusEntryService) DescribeExistingBonuses(existing []config.BonusEntry) dtos.ExistingBonusesDto {
	summary := dtos.ExistingBonusesDto{Keys: make(map[string]bool, len(existing))}
	for _, entry := range existing {
		summary.Keys[config_helpers.GetHash(entry)] = true
		if entry.PresetType == config.BonusSpell && entry.Param != "" {
			summary.SpellIDs = append(summary.SpellIDs, entry.Param)
		}
	}

	return summary
}

func (this *BonusEntryService) BuildBonusEntries(
	request dtos.BonusCompositionRequestDto) dtos.BonusCompositionResultDto {
	switch request.PresetType {
	case config.BonusTownPortalFree:
		return singleEntry(request, "")
	case config.BonusSpell:
		return this.buildSpellEntries(request)
	case config.BonusUnitMultiplier:
		return numericEntry(request, request.MultiplierText, "Enter a numeric multiplier.")
	case config.BonusMovementBonus:
		return numericEntry(request, request.MovementText, "Enter a numeric movement value.")
	case config.BonusStartingItem:
		id := strings.TrimSpace(request.ItemText)
		if id == "" {
			return dtos.BonusCompositionResultDto{Error: "Pick or enter an item."}
		}

		return singleEntry(request, id)
	default:
		return numericEntry(request, request.ResourceText, "Enter a numeric amount.")
	}
}

func (this *BonusEntryService) FilterNewBonusEntries(
	entries []config.BonusEntry,
	existingKeys map[string]bool) []config.BonusEntry {
	fresh := make([]config.BonusEntry, 0, len(entries))
	for _, entry := range entries {
		if !existingKeys[config_helpers.GetHash(entry)] {
			fresh = append(fresh, entry)
		}
	}

	return fresh
}

func (this *BonusEntryService) GetSpellCountLabel(count int) string {
	switch count {
	case 0:
		return "Spells"
	case 1:
		return "1 spell picked"
	default:
		return strconv.Itoa(count) + " spells picked"
	}
}

func (this *BonusEntryService) buildSpellEntries(
	request dtos.BonusCompositionRequestDto) dtos.BonusCompositionResultDto {
	if len(request.SelectedSpells) == 0 {
		return dtos.BonusCompositionResultDto{Error: "Pick at least one spell."}
	}

	free := spellNotFreeParam
	if request.MakeSpellsFree {
		free = spellFreeParam
	}

	entries := make([]config.BonusEntry, 0, len(request.SelectedSpells))
	for _, id := range request.SelectedSpells {
		entries = append(entries, config.BonusEntry{
			PresetType:     request.PresetType,
			ReceiverFilter: request.ReceiverFilter,
			Param:          id,
			Param2:         free,
		})
	}

	return dtos.BonusCompositionResultDto{Entries: entries}
}

func numericEntry(
	request dtos.BonusCompositionRequestDto,
	text string,
	errorText string) dtos.BonusCompositionResultDto {
	value := strings.TrimSpace(text)
	if !isNumeric(value) {
		return dtos.BonusCompositionResultDto{Error: errorText}
	}

	return singleEntry(request, value)
}

func singleEntry(request dtos.BonusCompositionRequestDto, param string) dtos.BonusCompositionResultDto {
	return dtos.BonusCompositionResultDto{
		Entries: []config.BonusEntry{{
			PresetType:     request.PresetType,
			ReceiverFilter: request.ReceiverFilter,
			Param:          param,
		}},
	}
}

func isNumeric(value string) bool {
	if value == "" {
		return false
	}

	_, err := strconv.ParseFloat(value, 64)

	return err == nil
}
