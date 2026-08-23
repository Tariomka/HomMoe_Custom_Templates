package test_helpers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/mock"
)

// ZoneContentEditorServiceMock is a testify mock of
// zone_content.IZoneContentEditorService, used to unit-test collaborators
// without the real editor rules.
type ZoneContentEditorServiceMock struct {
	mock.Mock
}

func (this *ZoneContentEditorServiceMock) ComposeContentRule(
	request dtos.ContentRuleCompositionRequestDto) dtos.ContentRuleCompositionResultDto {
	arguments := this.Called(request)
	result, _ := arguments.Get(0).(dtos.ContentRuleCompositionResultDto)
	return result
}

func (this *ZoneContentEditorServiceMock) UpsertContentRule(
	rules []models.ContentRuleRow,
	rule models.ContentRuleRow) []models.ContentRuleRow {
	arguments := this.Called(rules, rule)
	merged, _ := arguments.Get(0).([]models.ContentRuleRow)
	return merged
}

func (this *ZoneContentEditorServiceMock) GetDefaultContentRules(
	options dtos.ContentRuleEditorOptionsDto) []models.ContentRuleRow {
	arguments := this.Called(options)
	rules, _ := arguments.Get(0).([]models.ContentRuleRow)
	return rules
}

func (this *ZoneContentEditorServiceMock) GetContentRuleMarkers(
	descriptions []dtos.ContentRuleDescriptionDto) string {
	arguments := this.Called(descriptions)
	return arguments.String(0)
}

func (this *ZoneContentEditorServiceMock) GetContentRowDisplayName(
	name string,
	descriptions []dtos.ContentRuleDescriptionDto) string {
	arguments := this.Called(name, descriptions)
	return arguments.String(0)
}

func (this *ZoneContentEditorServiceMock) SortContentItemsByName(
	items []models.SidMapping) []models.SidMapping {
	arguments := this.Called(items)
	sorted, _ := arguments.Get(0).([]models.SidMapping)
	return sorted
}

func (this *ZoneContentEditorServiceMock) ClampContentCount(count int, maxCount int) int {
	arguments := this.Called(count, maxCount)
	return arguments.Int(0)
}
