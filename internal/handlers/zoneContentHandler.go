package handlers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
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
	rules []models.ContentRuleRowSave,
	rule models.ContentRuleRowSave) []models.ContentRuleRowSave {
	return this.zoneContentEditor.UpsertContentRule(rules, rule)
}

func (this *zoneContentHandler) GetDefaultContentRules(
	content models.SidMapping) []models.ContentRuleRowSave {
	return this.zoneContentEditor.GetDefaultContentRules(this.GetContentRuleEditorOptions(content))
}

func (this *zoneContentHandler) GetContentRuleMarkers(
	content models.SidMapping,
	rules []models.ContentRuleRowSave) string {
	return this.zoneContentEditor.GetContentRuleMarkers(this.describeContentRules(content, rules))
}

func (this *zoneContentHandler) GetContentRowDisplayName(
	content models.SidMapping,
	rules []models.ContentRuleRowSave) string {
	return this.zoneContentEditor.GetContentRowDisplayName(
		content.Name,
		this.describeContentRules(content, rules))
}

func (this *zoneContentHandler) SortContentItemsByName(items []models.SidMapping) []models.SidMapping {
	return this.zoneContentEditor.SortContentItemsByName(items)
}

func (this *zoneContentHandler) ClampContentCount(count int, maxCount int) int {
	return this.zoneContentEditor.ClampContentCount(count, maxCount)
}

func (this *zoneContentHandler) describeContentRules(
	content models.SidMapping,
	rules []models.ContentRuleRowSave) []dtos.ContentRuleDescriptionDto {
	descriptions := make([]dtos.ContentRuleDescriptionDto, 0, len(rules))
	for _, savedRule := range rules {
		descriptions = append(descriptions, this.DescribeContentRule(content, savedRule))
	}

	return descriptions
}
