package spells_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenComparedSpellsHaveDifferentSchools_OrdersByTheRegistrySchoolOrder(t *testing.T) {
	t.Parallel()
	// Arrange
	schools := registry.GetSpellSchoolTypeList()
	require.GreaterOrEqual(t, len(schools), 2)
	first := constants.SpellEntry{School: schools[0], Tier: 9}
	second := constants.SpellEntry{School: schools[1], Tier: 1}

	// Act
	comparison := constants.CompareSpellEntries(first, second)

	// Assert
	assert.Negative(t, comparison)
}

func TestWhenComparedSpellsShareASchool_OrdersByTier(t *testing.T) {
	t.Parallel()
	// Arrange
	school := registry.GetSpellSchoolTypeList()[0]
	first := constants.SpellEntry{School: school, Tier: 1, Name: "Zephyr"}
	second := constants.SpellEntry{School: school, Tier: 2, Name: "Amulet"}

	// Act
	comparison := constants.CompareSpellEntries(first, second)

	// Assert
	assert.Negative(t, comparison)
}

func TestWhenComparedSpellsShareASchoolAndTier_OrdersByName(t *testing.T) {
	t.Parallel()
	// Arrange
	school := registry.GetSpellSchoolTypeList()[0]
	first := constants.SpellEntry{School: school, Tier: 1, Name: "Amulet"}
	second := constants.SpellEntry{School: school, Tier: 1, Name: "Zephyr"}

	// Act
	comparison := constants.CompareSpellEntries(first, second)

	// Assert
	assert.Negative(t, comparison)
}

func TestWhenComparedSpellsHaveUnknownSchools_TreatsThemAsEqualRanked(t *testing.T) {
	t.Parallel()
	// Arrange
	name := gofakeit.Word()
	first := constants.SpellEntry{School: gofakeit.UUID(), Tier: 1, Name: name}
	second := constants.SpellEntry{School: gofakeit.UUID(), Tier: 1, Name: name}

	// Act
	comparison := constants.CompareSpellEntries(first, second)

	// Assert
	assert.Zero(t, comparison)
}
