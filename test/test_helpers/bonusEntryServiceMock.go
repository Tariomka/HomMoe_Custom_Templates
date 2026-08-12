package test_helpers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/stretchr/testify/mock"
)

// BonusEntryServiceMock is a testify mock of bonuses.IBonusEntryService, used
// to unit-test collaborators without composing real bonus entries.
type BonusEntryServiceMock struct {
	mock.Mock
}

func (this *BonusEntryServiceMock) DescribeExistingBonuses(
	existing []config.BonusEntry) dtos.ExistingBonusesDto {
	arguments := this.Called(existing)
	summary, _ := arguments.Get(0).(dtos.ExistingBonusesDto)
	return summary
}

func (this *BonusEntryServiceMock) BuildBonusEntries(
	request dtos.BonusCompositionRequestDto) dtos.BonusCompositionResultDto {
	arguments := this.Called(request)
	result, _ := arguments.Get(0).(dtos.BonusCompositionResultDto)
	return result
}

func (this *BonusEntryServiceMock) FilterNewBonusEntries(
	entries []config.BonusEntry,
	existingKeys map[string]bool) []config.BonusEntry {
	arguments := this.Called(entries, existingKeys)
	fresh, _ := arguments.Get(0).([]config.BonusEntry)
	return fresh
}

func (this *BonusEntryServiceMock) GetSpellCountLabel(count int) string {
	arguments := this.Called(count)
	return arguments.String(0)
}
