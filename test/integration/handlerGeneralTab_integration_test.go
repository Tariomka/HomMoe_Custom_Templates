//go:build integration_test

package integration_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers/integration_common"
	"github.com/stretchr/testify/assert"
)

// TestHandlerGeneralTab_SelectsAnOfficialMapSize pins the unshifted baseline that
// the experimental case below is measured against.
func TestHandlerGeneralTab_SelectsAnOfficialMapSize(t *testing.T) {
	// Arrange
	runner := integration_common.NewAppRunner(t)
	handler := integration_common.NewHandler(runner).ClickGeneralTab()
	runner.NextFrame()

	// Act
	handler.SelectMapSize(3)

	// Assert
	assert.Equal(t, 112, runner.CurrentState().MapSize)
}

// TestHandlerGeneralTab_SelectsAnExperimentalMapSize proves a coordinate
// computed after a layout shift lands on the widget it was meant for: the 11
// official sizes occupy rows 0-10, so row 12 exists at all only because the
// checkbox grew the dropdown to 28 rows.
func TestHandlerGeneralTab_SelectsAnExperimentalMapSize(t *testing.T) {
	// Arrange
	runner := integration_common.NewAppRunner(t)
	handler := integration_common.NewHandler(runner).
		ClickGeneralTab().
		ToggleExperimentalMapSizes()
	runner.NextFrame()

	// Act
	handler.SelectMapSize(12)

	// Assert
	assert.Equal(t, 272, runner.CurrentState().MapSize)
}

func TestHandlerGeneralTab_SelectsSingleHeroGameMode(t *testing.T) {
	// Arrange
	runner := integration_common.NewAppRunner(t)
	handler := integration_common.NewHandler(runner).ClickGeneralTab()
	runner.NextFrame()

	// Act
	handler.SelectGameMode(true)

	// Assert
	assert.Equal(t, "SingleHero", runner.CurrentState().GameMode)
}

func TestHandlerGeneralTab_TogglingExperimentalSizesUpdatesState(t *testing.T) {
	// Arrange
	runner := integration_common.NewAppRunner(t)
	handler := integration_common.NewHandler(runner).ClickGeneralTab()
	runner.NextFrame()

	// Act
	handler.ToggleExperimentalMapSizes()

	// Assert
	assert.True(t, runner.CurrentState().ExperimentalMapSizes)
}
