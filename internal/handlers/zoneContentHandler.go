package handlers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/content_rule_model"
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
	rule, valid := this.zoneContentEditor.ComposeContentRule(request.ContentRuleComposition)
	return dtos.ContentRuleCompositionResultDto{Rule: rule, Valid: valid}
}

func (this *zoneContentHandler) UpsertContentRule(
	rules []editor_state_model.ContentRuleRow,
	rule editor_state_model.ContentRuleRow) []editor_state_model.ContentRuleRow {
	return this.zoneContentEditor.UpsertContentRule(rules, rule)
}

func (this *zoneContentHandler) GetDefaultContentRules(content models.SidMapping) []editor_state_model.ContentRuleRow {
	options := linq.FromSlice(this.GetContentRuleEditorOptions(content).Rules).
		Select(func(option dtos.ContentRuleOptionDto) content_rule_model.ContentRuleOption {
			return option.ContentRuleOption
		}).
		ToSlice()
	return this.zoneContentEditor.GetDefaultContentRules(options)
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
	rules []editor_state_model.ContentRuleRow) []content_rule_model.ContentRuleDescription {
	return linq.FromSlice(rules).
		Select(func(rule editor_state_model.ContentRuleRow) content_rule_model.ContentRuleDescription {
			return this.DescribeContentRule(content, rule).ContentRuleDescription
		}).
		ToSlice()
}
