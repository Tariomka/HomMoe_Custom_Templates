package content_rules

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

// IContentRuleService resolves and applies the rules attached to mandatory
// content, letting callers substitute a stub catalogue in tests.
type IContentRuleService interface {
	GetRules() []IContentRule
	ApplyRulesToItem(item *entities.MandatoryContentItem, rules []IContentRule)
	CreateRuleFromSavedRule(saved models.ContentRuleRowSave, content models.SidMapping) IContentRule
	RestoreRulesFromRow(row models.ZoneContentRowSave, content models.SidMapping) []IContentRule
	GetDistanceDisplayNames() []string
	GetVariantsForContent(content models.SidMapping) []models.VariantMapping
	GetVariantForContentByID(content models.SidMapping, variantID int) (models.VariantMapping, bool)
}
