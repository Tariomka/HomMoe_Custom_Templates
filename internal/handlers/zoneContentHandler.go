package handlers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zone_content"
)

type zoneContentHandler struct {
	handler_interfaces.IContentRuleHandler

	zoneContentEditor zone_content.IZoneContentEditorService
}

func NewZoneContentHandler(
	contentRuleHandler handler_interfaces.IContentRuleHandler,
	zoneContentEditor zone_content.IZoneContentEditorService) handler_interfaces.IZoneContentHandler {
	return &zoneContentHandler{
		IContentRuleHandler: contentRuleHandler,
		zoneContentEditor:   zoneContentEditor,
	}
}

func (this *zoneContentHandler) ComposeContentRule(
	request dtos.ContentRuleCompositionRequestDto) dtos.ContentRuleCompositionResultDto {
	return this.zoneContentEditor.ComposeContentRule(request)
}

func (this *zoneContentHandler) UpsertContentRule(
	rules []editor_state_model.ContentRuleRow,
	rule editor_state_model.ContentRuleRow) []editor_state_model.ContentRuleRow {
	return this.zoneContentEditor.UpsertContentRule(rules, rule)
}

func (this *zoneContentHandler) GetDefaultContentRules(content models.SidMapping) []editor_state_model.ContentRuleRow {
	return this.zoneContentEditor.GetDefaultContentRules(this.GetContentRuleEditorOptions(content))
}

func (this *zoneContentHandler) GetContentRuleMarkers(
	content models.SidMapping,
	rules []editor_state_model.ContentRuleRow) string {
	return this.zoneContentEditor.GetContentRuleMarkers(this.describeContentRules(content, rules))
}

func (this *zoneContentHandler) GetContentRowDisplayName(
	content models.SidMapping,
	rules []editor_state_model.ContentRuleRow) string {
	return this.zoneContentEditor.GetContentRowDisplayName(content.Name, this.describeContentRules(content, rules))
}

func (this *zoneContentHandler) SortContentItemsByName(items []models.SidMapping) []models.SidMapping {
	return this.zoneContentEditor.SortContentItemsByName(items)
}

func (this *zoneContentHandler) ClampContentCount(count int, maxCount int) int {
	return this.zoneContentEditor.ClampContentCount(count, maxCount)
}

func (this *zoneContentHandler) describeContentRules(
	content models.SidMapping,
	rules []editor_state_model.ContentRuleRow) []dtos.ContentRuleDescriptionDto {
	return linq.FromSlice(rules).
		Select(func(rule editor_state_model.ContentRuleRow) dtos.ContentRuleDescriptionDto {
			return this.DescribeContentRule(content, rule)
		}).
		ToSlice()
}
