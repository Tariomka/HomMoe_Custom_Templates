package zone_content

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/content_rule_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
)

// IZoneContentEditorService holds the decision-making behind the zone-content
// editor: composing and merging content rules, and shaping the rows and
// catalogue the section presents.
type IZoneContentEditorService interface {
	ComposeContentRule(composition content_rule_model.ContentRuleComposition) (editor_state_model.ContentRuleRow, bool)
	UpsertContentRule(
		rules []editor_state_model.ContentRuleRow,
		rule editor_state_model.ContentRuleRow) []editor_state_model.ContentRuleRow
	GetDefaultContentRules(options []content_rule_model.ContentRuleOption) []editor_state_model.ContentRuleRow
	GetContentRuleMarkers(descriptions []content_rule_model.ContentRuleDescription) string
	GetContentRowDisplayName(name string, descriptions []content_rule_model.ContentRuleDescription) string
	SortContentItemsByName(items []models.SidMapping) []models.SidMapping
	ClampContentCount(count int, maxCount int) int
}
