package spells_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/constants"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenSchoolTypeIsKnown_ReturnsItsAccentColor(t *testing.T) {
	t.Parallel()
	// Arrange
	spellSchoolValues := registry.GetSpellSchoolTypeValues()

	// Act
	schoolColor := constants.GetSpellSchoolColor(spellSchoolValues.Nightshade)

	// Assert
	assert.Equal(t, themes.ColorsSpellSchools.Nightshade, schoolColor)
}

func TestWhenSchoolTypeIsUnknown_FallsBackToTheBaseAccentColor(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	schoolColor := constants.GetSpellSchoolColor(gofakeit.UUID())

	// Assert
	assert.Equal(t, themes.ColorsBase.Accent, schoolColor)
}

func TestWhenEverySchoolTypeIsMapped_EachOneHasItsOwnColor(t *testing.T) {
	t.Parallel()
	// Arrange
	spellSchoolValues := registry.GetSpellSchoolTypeValues()
	schoolTypes := []string{
		spellSchoolValues.HighNeutral,
		spellSchoolValues.Daylight,
		spellSchoolValues.Nightshade,
		spellSchoolValues.Arcane,
		spellSchoolValues.Primal,
	}
	expected := []any{
		themes.ColorsSpellSchools.HighNeutral,
		themes.ColorsSpellSchools.Daylight,
		themes.ColorsSpellSchools.Nightshade,
		themes.ColorsSpellSchools.Arcane,
		themes.ColorsSpellSchools.Primal,
	}

	// Act
	actual := make([]any, 0, len(schoolTypes))
	for _, schoolType := range schoolTypes {
		actual = append(actual, constants.GetSpellSchoolColor(schoolType))
	}

	// Assert
	assert.Equal(t, expected, actual)
}
