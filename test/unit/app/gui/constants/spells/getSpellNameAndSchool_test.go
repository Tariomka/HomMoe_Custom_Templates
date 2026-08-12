package spells_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenSidIsKnown_ReturnsSpellDisplayName(t *testing.T) {
	t.Parallel()
	// Arrange
	expected, found := constants.FindSpell(registry.GetPrimalSpellSidValues().Armageddon)
	require.True(t, found)

	// Act
	name, _ := constants.GetSpellNameAndSchool(expected.Sid)

	// Assert
	assert.Equal(t, expected.Name, name)
}

func TestWhenSidIsKnown_ReturnsSpellSchoolDisplayName(t *testing.T) {
	t.Parallel()
	// Arrange
	sid := registry.GetPrimalSpellSidValues().Armageddon
	expected, found := constants.FindSpell(sid)
	require.True(t, found)

	// Act
	_, school := constants.GetSpellNameAndSchool(sid)

	// Assert
	assert.Equal(t, constants.GetSpellSchoolDisplayName(expected.School), school)
}

func TestWhenSidIsUnknown_ReturnsSentenceCasedSid(t *testing.T) {
	t.Parallel()
	// Arrange
	sid := "not_a_real_spell"

	// Act
	name, _ := constants.GetSpellNameAndSchool(sid)

	// Assert
	assert.Equal(t, "Not a real spell", name)
}

func TestWhenSidIsUnknown_ReturnsGenericSchoolLabel(t *testing.T) {
	t.Parallel()
	// Arrange
	sid := "not_a_real_spell"

	// Act
	_, school := constants.GetSpellNameAndSchool(sid)

	// Assert
	assert.Equal(t, "Spell", school)
}
