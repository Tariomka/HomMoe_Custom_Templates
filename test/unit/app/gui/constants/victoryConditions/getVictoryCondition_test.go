package victoryConditions_test

import (
	"strings"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/constants"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenIdIsKnown_ReturnsMatchingVictoryCondition(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := constants.GetVictoryConditionValues().HoldCity

	// Act
	victory, _ := constants.GetVictoryCondition(expected.ID)

	// Assert
	assert.Equal(t, expected, victory)
}

func TestWhenIdIsKnown_ReportsFound(t *testing.T) {
	t.Parallel()
	// Arrange
	knownID := constants.GetVictoryConditionValues().Tournament.ID

	// Act
	_, found := constants.GetVictoryCondition(knownID)

	// Assert
	assert.True(t, found)
}

func TestWhenIdDiffersOnlyInCase_ReturnsMatchingVictoryCondition(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := constants.GetVictoryConditionValues().GuardianArena

	// Act
	victory, _ := constants.GetVictoryCondition(strings.ToUpper(expected.ID))

	// Assert
	assert.Equal(t, expected, victory)
}

func TestWhenIdIsUnknown_ReportsNotFound(t *testing.T) {
	t.Parallel()
	// Arrange
	unknownID := "not-a-victory-" + gofakeit.LetterN(12)

	// Act
	_, found := constants.GetVictoryCondition(unknownID)

	// Assert
	assert.False(t, found)
}

func TestWhenIdIsUnknown_ReturnsZeroVictory(t *testing.T) {
	t.Parallel()
	// Arrange
	unknownID := "not-a-victory-" + gofakeit.LetterN(12)

	// Act
	victory, _ := constants.GetVictoryCondition(unknownID)

	// Assert
	assert.Equal(t, constants.Victory{}, victory)
}
