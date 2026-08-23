package zone_content

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

// IZoneContentEditorService holds the decision-making behind the zone-content
// editor: composing and merging content rules, and shaping the rows and
// catalogue the section presents.
type IZoneContentEditorService interface {
	ComposeContentRule(request dtos.ContentRuleCompositionRequestDto) dtos.ContentRuleCompositionResultDto
	UpsertContentRule(rules []models.ContentRuleRow, rule models.ContentRuleRow) []models.ContentRuleRow
	GetDefaultContentRules(options dtos.ContentRuleEditorOptionsDto) []models.ContentRuleRow
	GetContentRuleMarkers(descriptions []dtos.ContentRuleDescriptionDto) string
	GetContentRowDisplayName(name string, descriptions []dtos.ContentRuleDescriptionDto) string
	SortContentItemsByName(items []models.SidMapping) []models.SidMapping
	ClampContentCount(count int, maxCount int) int
}
