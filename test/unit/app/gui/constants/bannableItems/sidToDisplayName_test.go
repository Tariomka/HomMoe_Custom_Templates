package bannableItems_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/constants"
	"github.com/stretchr/testify/assert"
)

func TestWhenSidHasTheArtifactSuffix_DropsIt(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	displayName := constants.SidToDisplayName("pole_star_artifact")

	// Assert
	assert.Equal(t, "Pole star", displayName)
}

func TestWhenSidHasUnderscores_ReplacesThemWithSpaces(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	displayName := constants.SidToDisplayName("seven_league_boots")

	// Assert
	assert.Equal(t, "Seven league boots", displayName)
}

func TestWhenSidIsEmpty_ReturnsItUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	displayName := constants.SidToDisplayName("")

	// Assert
	assert.Empty(t, displayName)
}

func TestWhenSidIsOnlyTheArtifactSuffix_ReturnsTheOriginalSid(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	displayName := constants.SidToDisplayName("_artifact")

	// Assert
	assert.Equal(t, "_artifact", displayName)
}
