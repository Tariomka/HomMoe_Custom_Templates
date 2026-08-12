package valueOverrideSids_test

import (
	"slices"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/constants"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenNothingIsExcluded_ReturnsANonEmptyList(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	sids := constants.GetValueOverrideSidsWithExclusions(nil)

	// Assert
	assert.NotEmpty(t, sids)
}

func TestWhenListIsBuilt_ContainsNoEmptySids(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	sids := constants.GetValueOverrideSidsWithExclusions(nil)

	// Assert
	assert.NotContains(t, sids, "")
}

func TestWhenListIsBuilt_ContainsNoDuplicates(t *testing.T) {
	t.Parallel()
	// Arrange
	sids := constants.GetValueOverrideSidsWithExclusions(nil)

	// Act
	compacted := slices.Compact(slices.Clone(sids))

	// Assert
	assert.Len(t, compacted, len(sids))
}

func TestWhenListIsBuilt_IsSortedAlphabetically(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	sids := constants.GetValueOverrideSidsWithExclusions(nil)

	// Assert
	assert.True(t, slices.IsSorted(sids))
}

func TestWhenSidIsExcluded_ItIsAbsentFromTheList(t *testing.T) {
	t.Parallel()
	// Arrange
	full := constants.GetValueOverrideSidsWithExclusions(nil)

	// Act
	sids := constants.GetValueOverrideSidsWithExclusions([]string{full[0]})

	// Assert
	assert.NotContains(t, sids, full[0])
}

func TestWhenSidIsExcluded_TheRestOfTheListIsKept(t *testing.T) {
	t.Parallel()
	// Arrange
	full := constants.GetValueOverrideSidsWithExclusions(nil)

	// Act
	sids := constants.GetValueOverrideSidsWithExclusions([]string{full[0]})

	// Assert
	assert.Len(t, sids, len(full)-1)
}

func TestWhenExclusionsDoNotMatchAnything_ReturnsTheWholeList(t *testing.T) {
	t.Parallel()
	// Arrange
	full := constants.GetValueOverrideSidsWithExclusions(nil)

	// Act
	sids := constants.GetValueOverrideSidsWithExclusions([]string{gofakeit.UUID()})

	// Assert
	assert.Equal(t, full, sids)
}

func TestWhenExclusionSliceIsPassed_ItIsNotMutated(t *testing.T) {
	t.Parallel()
	// Arrange
	full := constants.GetValueOverrideSidsWithExclusions(nil)
	excluded := []string{full[0], full[1]}
	original := slices.Clone(excluded)

	// Act
	_ = constants.GetValueOverrideSidsWithExclusions(excluded)

	// Assert
	assert.Equal(t, original, excluded)
}
