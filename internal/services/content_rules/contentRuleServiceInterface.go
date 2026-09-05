package content_rules

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

// IContentRuleService resolves and applies the rules attached to mandatory
// content, letting callers substitute a stub catalogue in tests.
type IContentRuleService interface {
	GetRules() []IContentRule
	ApplyRulesToItem(item *template_model.MandatoryContentItem, rules []IContentRule)
	CreateRuleFromSavedRule(saved editor_state_model.ContentRuleRow, content models.SidMapping) IContentRule
	RestoreRulesFromRow(row editor_state_model.ZoneContentRow, content models.SidMapping) []IContentRule
	GetDistanceDisplayNames() []string
	GetVariantsForContent(content models.SidMapping) []models.VariantMapping
	GetVariantForContentByID(content models.SidMapping, variantID int) (models.VariantMapping, bool)
}
