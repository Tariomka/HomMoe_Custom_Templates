package test_helpers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/mock"
)

// ContentRuleServiceMock is a testify mock of content_rules.IContentRuleService,
// used to unit-test collaborators without the real rule catalogue.
type ContentRuleServiceMock struct {
	mock.Mock
}

func (this *ContentRuleServiceMock) GetRules() []content_rules.IContentRule {
	arguments := this.Called()
	rules, _ := arguments.Get(0).([]content_rules.IContentRule)
	return rules
}

func (this *ContentRuleServiceMock) ApplyRulesToItem(
	item *entities.MandatoryContentItem,
	rules []content_rules.IContentRule) {
	this.Called(item, rules)
}

func (this *ContentRuleServiceMock) CreateRuleFromSavedRule(
	saved models.ContentRuleRow,
	content models.SidMapping) content_rules.IContentRule {
	arguments := this.Called(saved, content)
	rule, _ := arguments.Get(0).(content_rules.IContentRule)
	return rule
}

func (this *ContentRuleServiceMock) RestoreRulesFromRow(
	row models.ZoneContentRow,
	content models.SidMapping) []content_rules.IContentRule {
	arguments := this.Called(row, content)
	rules, _ := arguments.Get(0).([]content_rules.IContentRule)
	return rules
}

func (this *ContentRuleServiceMock) GetDistanceDisplayNames() []string {
	arguments := this.Called()
	names, _ := arguments.Get(0).([]string)
	return names
}

func (this *ContentRuleServiceMock) GetVariantsForContent(content models.SidMapping) []models.VariantMapping {
	arguments := this.Called(content)
	variants, _ := arguments.Get(0).([]models.VariantMapping)
	return variants
}

func (this *ContentRuleServiceMock) GetVariantForContentByID(
	content models.SidMapping,
	variantID int) (models.VariantMapping, bool) {
	arguments := this.Called(content, variantID)
	variant, _ := arguments.Get(0).(models.VariantMapping)
	return variant, arguments.Bool(1)
}
