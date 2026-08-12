package gameModes_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/stretchr/testify/assert"
)

func TestWhenGameModesAreRequested_ReturnsTheClassicAndSingleHeroSids(t *testing.T) {
	t.Parallel()
	// Arrange
	gameModes := registry.GetGameModeValues()

	// Act
	modes := constants.GetGameModes()

	// Assert
	assert.Equal(t, []string{gameModes.Classic, gameModes.SingleHero}, modes)
}
