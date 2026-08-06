package spells_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/constants"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenDisplayNameIsKnown_ReturnsItsAccentColor(t *testing.T) {
	t.Parallel()
	// Arrange
	displayName := constants.GetSpellSchoolDisplayName(registry.GetSpellSchoolTypeValues().Primal)

	// Act
	schoolColor := constants.GetSpellSchoolColorFromDisplayName(displayName)

	// Assert
	assert.Equal(t, themes.ColorsSpellSchools.Primal, schoolColor)
}

func TestWhenDisplayNameIsUnknown_FallsBackToTheBaseAccentColor(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	schoolColor := constants.GetSpellSchoolColorFromDisplayName(gofakeit.UUID())

	// Assert
	assert.Equal(t, themes.ColorsBase.Accent, schoolColor)
}

func TestWhenEveryDisplayNameIsMapped_MatchesTheSchoolTypeMapping(t *testing.T) {
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

	// Act
	byDisplayName := make([]any, 0, len(schoolTypes))
	bySchoolType := make([]any, 0, len(schoolTypes))
	for _, schoolType := range schoolTypes {
		displayName := constants.GetSpellSchoolDisplayName(schoolType)
		byDisplayName = append(byDisplayName, constants.GetSpellSchoolColorFromDisplayName(displayName))
		bySchoolType = append(bySchoolType, constants.GetSpellSchoolColor(schoolType))
	}

	// Assert
	assert.Equal(t, bySchoolType, byDisplayName)
}
