package mappers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
)

type MandatoryContentItemMapper struct {
	contentRuleService *content_rules.ContentRuleService
}

func NewMandatoryContentItemMapper() *MandatoryContentItemMapper {
	return &MandatoryContentItemMapper{contentRuleService: content_rules.NewContentRuleService()}
}

func (this *MandatoryContentItemMapper) FromRows(
	rows []models.ZoneContentRowSave,
) []entities.MandatoryContentItem {
	if len(rows) == 0 {
		return nil
	}

	var items []entities.MandatoryContentItem
	for _, rawRow := range rows {
		row := rawRow.Normalized()
		if row.Sid == "" {
			continue
		}
		for range row.Count {
			items = append(items, this.fromRow(row))
		}
	}
	return items
}

func (this *MandatoryContentItemMapper) fromRow(
	row models.ZoneContentRowSave,
) entities.MandatoryContentItem {
	item := entities.MandatoryContentItem{IsMine: row.IsMine}
	if row.IsGroup {
		item.IncludeLists = []string{row.Sid}
	} else {
		item.SID = row.Sid
	}

	rules := this.contentRuleService.RestoreRulesFromRow(row, models.SidMapping{Sid: row.Sid})
	this.contentRuleService.ApplyRulesToItem(&item, rules)
	return item
}
