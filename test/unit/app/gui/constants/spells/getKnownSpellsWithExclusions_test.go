package spells_test

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
	spells := constants.GetKnownSpellsWithExclusions(nil)

	// Assert
	assert.NotEmpty(t, spells)
}

func TestWhenCatalogIsBuilt_EveryEntryIsFullyPopulated(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	spells := constants.GetKnownSpellsWithExclusions(nil)

	// Assert
	for _, spell := range spells {
		require.NotEmpty(t, spell.Sid, "entry %q has an empty Sid", spell.Name)
		require.NotEmpty(t, spell.Name, "entry %q has an empty Name", spell.Sid)
		require.NotEmpty(t, spell.School, "entry %q has an empty School", spell.Sid)
		require.Positive(t, spell.Tier, "entry %q has a non-positive Tier", spell.Sid)
	}
}

func TestWhenCatalogIsBuilt_EverySidIsUnique(t *testing.T) {
	t.Parallel()
	// Arrange
	spells := constants.GetKnownSpellsWithExclusions(nil)

	// Act
	sids := make([]string, 0, len(spells))
	for _, spell := range spells {
		sids = append(sids, spell.Sid)
	}
	slices.Sort(sids)

	// Assert
	assert.Len(t, slices.Compact(sids), len(spells))
}

func TestWhenCatalogIsBuilt_IsSortedBySchoolThenTier(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	spells := constants.GetKnownSpellsWithExclusions(nil)

	// Assert
	assert.True(t, slices.IsSortedFunc(spells, constants.CompareSpellEntries))
}

func TestWhenSidIsExcluded_ItIsAbsentFromTheCatalog(t *testing.T) {
	t.Parallel()
	// Arrange
	full := constants.GetKnownSpellsWithExclusions(nil)

	// Act
	spells := constants.GetKnownSpellsWithExclusions([]string{full[0].Sid})

	// Assert
	assert.NotContains(t, spells, full[0])
}

func TestWhenSidIsExcluded_TheRestOfTheCatalogIsKept(t *testing.T) {
	t.Parallel()
	// Arrange
	full := constants.GetKnownSpellsWithExclusions(nil)

	// Act
	spells := constants.GetKnownSpellsWithExclusions([]string{full[0].Sid})

	// Assert
	assert.Len(t, spells, len(full)-1)
}

func TestWhenExclusionsDoNotMatchAnything_ReturnsTheWholeCatalog(t *testing.T) {
	t.Parallel()
	// Arrange
	full := constants.GetKnownSpellsWithExclusions(nil)

	// Act
	spells := constants.GetKnownSpellsWithExclusions([]string{gofakeit.UUID()})

	// Assert
	assert.Equal(t, full, spells)
}

func TestWhenCatalogIsRequestedTwice_TheFirstResultIsNotAliased(t *testing.T) {
	t.Parallel()
	// Arrange
	first := constants.GetKnownSpellsWithExclusions(nil)

	// Act
	first[0] = constants.SpellEntry{Sid: gofakeit.UUID()}
	second := constants.GetKnownSpellsWithExclusions(nil)

	// Assert
	assert.NotEqual(t, first[0], second[0])
}
