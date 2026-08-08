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
	UpsertContentRule(
		rules []models.ContentRuleRowSave,
		rule models.ContentRuleRowSave) []models.ContentRuleRowSave
	GetDefaultContentRules(options dtos.ContentRuleEditorOptionsDto) []models.ContentRuleRowSave
	GetContentRuleMarkers(descriptions []dtos.ContentRuleDescriptionDto) string
	GetContentRowDisplayName(name string, descriptions []dtos.ContentRuleDescriptionDto) string
	SortContentItemsByName(items []models.SidMapping) []models.SidMapping
	ClampContentCount(count int, maxCount int) int
}
