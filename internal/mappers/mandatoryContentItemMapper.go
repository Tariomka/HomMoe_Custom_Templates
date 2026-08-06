package mappers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/mandatory_content"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
)

type MandatoryContentItemMapper struct {
	contentRuleService content_rules.IContentRuleService
}

func NewMandatoryContentItemMapper(contentRuleService content_rules.IContentRuleService) IMandatoryContentItemMapper {
	return &MandatoryContentItemMapper{contentRuleService: contentRuleService}
}

func (this *MandatoryContentItemMapper) FromRows(rows []models.ZoneContentRowSave) []entities.MandatoryContentItem {
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

func (this *MandatoryContentItemMapper) fromRow(row models.ZoneContentRowSave) entities.MandatoryContentItem {
	sid := row.Sid
	if row.IsGroup {
		sid = ""
	}
	builder := mandatory_content.NewContentItemBuilder(sid)
	if row.IsGroup {
		builder.WithIncludeList(row.Sid)
	}
	if row.IsMine {
		builder.WithMine()
	}
	item := builder.Build()

	rules := this.contentRuleService.RestoreRulesFromRow(row, models.SidMapping{Sid: row.Sid})
	this.contentRuleService.ApplyRulesToItem(&item, rules)
	return item
}
