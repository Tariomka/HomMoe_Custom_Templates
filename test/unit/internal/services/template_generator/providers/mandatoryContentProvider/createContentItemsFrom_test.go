package mandatoryContentProvider_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers"
	"github.com/stretchr/testify/assert"
)

func TestWhenRowsAreEmpty_ReturnsNil(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := providers.NewMandatoryContentProvider()

	// Act
	actual := provider.CreateContentItemsFrom(nil)

	// Assert
	assert.Nil(t, actual)
}

func TestWhenRowSidIsEmpty_SkipsRow(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := providers.NewMandatoryContentProvider()
	rows := []models.ZoneContentRowSave{{Sid: "", Count: 2}}

	// Act
	actual := provider.CreateContentItemsFrom(rows)

	// Assert
	assert.Empty(t, actual)
}

func TestWhenRowCountIsThree_CreatesThreeIdenticalItems(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := providers.NewMandatoryContentProvider()
	rows := []models.ZoneContentRowSave{{Sid: "sawmill", Count: 3}}

	// Act
	actual := provider.CreateContentItemsFrom(rows)

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
	provider := providers.NewMandatoryContentProvider()
	rows := []models.ZoneContentRowSave{{Sid: "sawmill", Count: 0}}

	// Act
	actual := provider.CreateContentItemsFrom(rows)

	// Assert
	assert.Equal(t, []entities.MandatoryContentItem{{SID: "sawmill"}}, actual)
}

func TestWhenRowIsGroup_SetsIncludeListsInsteadOfSid(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := providers.NewMandatoryContentProvider()
	rows := []models.ZoneContentRowSave{{Sid: "include_list_dwellings", Count: 1, IsGroup: true}}

	// Act
	actual := provider.CreateContentItemsFrom(rows)

	// Assert
	assert.Equal(t, []entities.MandatoryContentItem{
		{IncludeLists: []string{"include_list_dwellings"}},
	}, actual)
}

func TestWhenRowIsMine_SetsIsMineOnItem(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := providers.NewMandatoryContentProvider()
	rows := []models.ZoneContentRowSave{{Sid: "gold_mine", Count: 1, IsMine: true}}

	// Act
	actual := provider.CreateContentItemsFrom(rows)

	// Assert
	assert.Equal(t, []entities.MandatoryContentItem{{SID: "gold_mine", IsMine: true}}, actual)
}

func TestWhenRowHasGuardedRule_AppliesGuardedFlagToItem(t *testing.T) {
	t.Parallel()
	// Arrange
	guarded := true
	provider := providers.NewMandatoryContentProvider()
	rows := []models.ZoneContentRowSave{{
		Sid:   "sawmill",
		Count: 1,
		Rules: []models.ContentRuleRowSave{{Name: "Guarded", IsGuarded: &guarded}},
	}}

	// Act
	actual := provider.CreateContentItemsFrom(rows)

	// Assert
	assert.Equal(t, []entities.MandatoryContentItem{{SID: "sawmill", IsGuarded: true}}, actual)
}
