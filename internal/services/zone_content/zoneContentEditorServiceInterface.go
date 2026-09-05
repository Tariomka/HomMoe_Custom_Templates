package zone_content

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
)

// IZoneContentEditorService holds the decision-making behind the zone-content
// editor: composing and merging content rules, and shaping the rows and
// catalogue the section presents.
type IZoneContentEditorService interface {
	ComposeContentRule(request dtos.ContentRuleCompositionRequestDto) dtos.ContentRuleCompositionResultDto
	UpsertContentRule(
		rules []editor_state_model.ContentRuleRow,
		rule editor_state_model.ContentRuleRow) []editor_state_model.ContentRuleRow
	GetDefaultContentRules(options dtos.ContentRuleEditorOptionsDto) []editor_state_model.ContentRuleRow
	GetContentRuleMarkers(descriptions []dtos.ContentRuleDescriptionDto) string
	GetContentRowDisplayName(name string, descriptions []dtos.ContentRuleDescriptionDto) string
	SortContentItemsByName(items []models.SidMapping) []models.SidMapping
	ClampContentCount(count int, maxCount int) int
}
