package handler_interfaces

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

// IZoneContentHandler backs the zone-content editor: it composes and merges the
// content rules attached to a row and shapes how rows and the content catalogue
// are presented. It embeds IContentRuleHandler because every zone-content view
// also needs the rule catalogue and rule descriptions.
type IZoneContentHandler interface {
	IContentRuleHandler

	ComposeContentRule(request dtos.ContentRuleCompositionRequestDto) dtos.ContentRuleCompositionResultDto
	UpsertContentRule(rules []models.ContentRuleRow, rule models.ContentRuleRow) []models.ContentRuleRow
	GetDefaultContentRules(content models.SidMapping) []models.ContentRuleRow
	GetContentRuleMarkers(content models.SidMapping, rules []models.ContentRuleRow) string
	GetContentRowDisplayName(content models.SidMapping, rules []models.ContentRuleRow) string
	SortContentItemsByName(items []models.SidMapping) []models.SidMapping
	ClampContentCount(count int, maxCount int) int
}
