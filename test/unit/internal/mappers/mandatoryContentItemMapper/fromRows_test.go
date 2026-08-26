package mandatoryContentItemMapper_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
)

func TestWhenRowsAreEmpty_ReturnsNil(t *testing.T) {
	t.Parallel()
	// Arrange
	mapper := mappers.NewMandatoryContentItemMapper(content_rules.NewContentRuleService())

	// Act
	actual := mapper.FromRows(nil)

	// Assert
	assert.Nil(t, actual)
}

func TestWhenRowSidIsEmpty_SkipsRow(t *testing.T) {
	t.Parallel()
	// Arrange
	mapper := mappers.NewMandatoryContentItemMapper(content_rules.NewContentRuleService())
	rows := []editor_state_model.ZoneContentRow{{Sid: "", Count: 2}}

	// Act
	actual := mapper.FromRows(rows)

	// Assert
	assert.Empty(t, actual)
}

func TestWhenRowCountIsThree_CreatesThreeIdenticalItems(t *testing.T) {
	t.Parallel()
	// Arrange
	mapper := mappers.NewMandatoryContentItemMapper(content_rules.NewContentRuleService())
	rows := []editor_state_model.ZoneContentRow{{Sid: "sawmill", Count: 3}}

	// Act
	actual := mapper.FromRows(rows)

	// Assert
	assert.Equal(t, []entities.MandatoryContentItem{
		{SID: "sawmill"},
		{SID: "sawmill"},
		{SID: "sawmill"},
	}, actual)
}

func TestWhenRowCountIsBelowOne_NormalizesToSingleItem(t *testing.T) {
	t.Parallel()
	// Arrange
	mapper := mappers.NewMandatoryContentItemMapper(content_rules.NewContentRuleService())
	rows := []editor_state_model.ZoneContentRow{{Sid: "sawmill", Count: 0}}

	// Act
	actual := mapper.FromRows(rows)

	// Assert
	assert.Equal(t, []entities.MandatoryContentItem{{SID: "sawmill"}}, actual)
}

func TestWhenRowIsGroup_SetsIncludeListsInsteadOfSid(t *testing.T) {
	t.Parallel()
	// Arrange
	mapper := mappers.NewMandatoryContentItemMapper(content_rules.NewContentRuleService())
	rows := []editor_state_model.ZoneContentRow{{Sid: "include_list_dwellings", Count: 1, IsGroup: true}}

	// Act
	actual := mapper.FromRows(rows)

	// Assert
	assert.Equal(t, []entities.MandatoryContentItem{
		{IncludeLists: []string{"include_list_dwellings"}},
	}, actual)
}

func TestWhenRowIsMine_SetsIsMineOnItem(t *testing.T) {
	t.Parallel()
	// Arrange
	mapper := mappers.NewMandatoryContentItemMapper(content_rules.NewContentRuleService())
	rows := []editor_state_model.ZoneContentRow{{Sid: "gold_mine", Count: 1, IsMine: true}}

	// Act
	actual := mapper.FromRows(rows)

	// Assert
	assert.Equal(t, []entities.MandatoryContentItem{{SID: "gold_mine", IsMine: true}}, actual)
}

func TestWhenRowHasGuardedRule_AppliesGuardedFlagToItem(t *testing.T) {
	t.Parallel()
	// Arrange
	guarded := true
	mapper := mappers.NewMandatoryContentItemMapper(content_rules.NewContentRuleService())
	rows := []editor_state_model.ZoneContentRow{{
		Sid:   "sawmill",
		Count: 1,
		Rules: []editor_state_model.ContentRuleRow{{Name: "Guarded", IsGuarded: &guarded}},
	}}

	// Act
	actual := mapper.FromRows(rows)

	// Assert
	assert.Equal(t, []entities.MandatoryContentItem{{SID: "sawmill", IsGuarded: true}}, actual)
}
