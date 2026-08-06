package bannableItems_test

import (
	"slices"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/constants"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenNothingIsExcluded_ReturnsANonEmptyCatalog(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	items := constants.GetBannableItemsWithExclusions(nil)

	// Assert
	assert.NotEmpty(t, items)
}

func TestWhenCatalogIsBuilt_EveryEntryIsFullyPopulated(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	items := constants.GetBannableItemsWithExclusions(nil)

	// Assert
	for _, item := range items {
		require.NotEmpty(t, item.Sid, "entry %q has an empty Sid", item.Name)
		require.NotEmpty(t, item.Name, "entry %q has an empty Name", item.Sid)
		require.NotEmpty(t, item.Category, "entry %q has an empty Category", item.Sid)
	}
}

func TestWhenCatalogIsBuilt_EverySidIsUnique(t *testing.T) {
	t.Parallel()
	// Arrange
	items := constants.GetBannableItemsWithExclusions(nil)

	// Act
	sids := sidsOf(items)
	slices.Sort(sids)

	// Assert
	assert.Len(t, slices.Compact(sids), len(items))
}

func TestWhenCatalogIsBuilt_UsesOnlyTheKnownCategories(t *testing.T) {
	t.Parallel()
	// Arrange
	knownCategories := []string{"Combat", "Diplomacy", "Magic", "Misc", "Movement", "Set"}

	// Act
	items := constants.GetBannableItemsWithExclusions(nil)

	// Assert
	categories := make([]string, 0, len(items))
	for _, item := range items {
		categories = append(categories, item.Category)
	}
	slices.Sort(categories)
	assert.Equal(t, knownCategories, slices.Compact(categories))
}

func TestWhenCatalogIsBuilt_IsSortedByCategoryThenName(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	items := constants.GetBannableItemsWithExclusions(nil)

	// Assert
	assert.True(t, slices.IsSortedFunc(items, constants.CompareBannableItems))
}

func TestWhenSidsAreExcluded_TheyAreAbsentFromTheCatalog(t *testing.T) {
	t.Parallel()
	// Arrange
	full := constants.GetBannableItemsWithExclusions(nil)
	excluded := []string{full[0].Sid, full[len(full)/2].Sid}

	// Act
	items := constants.GetBannableItemsWithExclusions(excluded)

	// Assert
	assert.NotContains(t, sidsOf(items), excluded[0])
	assert.NotContains(t, sidsOf(items), excluded[1])
}

func TestWhenSidsAreExcluded_TheRestOfTheCatalogIsKept(t *testing.T) {
	t.Parallel()
	// Arrange
	full := constants.GetBannableItemsWithExclusions(nil)

	// Act
	items := constants.GetBannableItemsWithExclusions([]string{full[0].Sid})

	// Assert
	assert.Len(t, items, len(full)-1)
}

func TestWhenExclusionsDoNotMatchAnything_ReturnsTheWholeCatalog(t *testing.T) {
	t.Parallel()
	// Arrange
	full := constants.GetBannableItemsWithExclusions(nil)

	// Act
	items := constants.GetBannableItemsWithExclusions([]string{gofakeit.UUID()})

	// Assert
	assert.Equal(t, full, items)
}

func TestWhenExclusionSliceIsPassed_ItIsNotMutated(t *testing.T) {
	t.Parallel()
	// Arrange
	full := constants.GetBannableItemsWithExclusions(nil)
	excluded := []string{full[0].Sid, full[1].Sid}
	original := slices.Clone(excluded)

	// Act
	_ = constants.GetBannableItemsWithExclusions(excluded)

	// Assert
	assert.Equal(t, original, excluded)
}

func TestWhenCatalogIsRequestedTwice_TheFirstResultIsNotAliased(t *testing.T) {
	t.Parallel()
	// Arrange
	first := constants.GetBannableItemsWithExclusions(nil)

	// Act
	first[0] = constants.BannableItemEntry{Sid: gofakeit.UUID()}
	second := constants.GetBannableItemsWithExclusions(nil)

	// Assert
	assert.NotEqual(t, first[0], second[0])
}

func sidsOf(items []constants.BannableItemEntry) []string {
	sids := make([]string, 0, len(items))
	for _, item := range items {
		sids = append(sids, item.Sid)
	}
	return sids
}
