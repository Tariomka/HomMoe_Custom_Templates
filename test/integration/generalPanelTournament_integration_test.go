//go:build integration_test

package integration_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers/integration_common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The victory dropdown rows are addressed by what they say; these are the two
// labels these tests switch between, as app/gui/constants spells them.
const (
	tournamentVictoryLabel = "Tournament"
	standardVictoryLabel   = "Standard"
)

// writeGeneralPanelFixture writes a settings file the General tab can be loaded
// from, so the panel's widgets start from the state the test is about.
func writeGeneralPanelFixture(t *testing.T, mutate func(state *editor_state_model.EditorState)) string {
	t.Helper()
	directory := t.TempDir()
	author := newUIState()
	author.UpdateState(func(state *editor_state_model.EditorState) {
		state.TemplateName = "General Fixture"
		mutate(state)
	})
	author.SaveStateToFile(filepath.Join(directory, "fixture.gen.json"))
	return filepath.Join(directory, "General Fixture.gen.json")
}

// TestGeneralTab_PlayerSliderRespondsWithoutTournament is the control for the
// disabled cases below: the very same drag has to move an unlocked slider, or
// the coordinate would prove nothing.
func TestGeneralTab_PlayerSliderRespondsWithoutTournament(t *testing.T) {
	// Arrange
	runner := integration_common.NewAppRunner(t)
	handler := integration_common.NewHandler(runner).ClickGeneralTab()
	runner.NextFrame()

	// Act
	handler.DragPlayerCountToMaximum()

	// Assert
	assert.Equal(t, 8, runner.CurrentState().PlayerCount)
}

func TestGeneralTab_PlayerSliderIsDisabledUnderTheTournamentCheckbox(t *testing.T) {
	// Arrange
	path := writeGeneralPanelFixture(t, func(state *editor_state_model.EditorState) { state.Tournament = true })
	runner := integration_common.NewAppRunner(t)
	handler := integration_common.NewHandler(runner).ClickGeneralTab()
	runner.LoadStateFromFile(path)
	runner.NextFrame()
	require.Equal(t, 2, runner.CurrentState().PlayerCount, "precondition: the fixture loads at two players")

	// Act
	handler.DragPlayerCountToMaximum()

	// Assert
	assert.Equal(t, 2, runner.CurrentState().PlayerCount)
}

func TestGeneralTab_PlayerSliderIsDisabledUnderTheTournamentVictoryCondition(t *testing.T) {
	// Arrange
	path := writeGeneralPanelFixture(t, func(state *editor_state_model.EditorState) {
		state.VictoryCondition = registry.GetWinningConditionValues().Tournament
	})
	runner := integration_common.NewAppRunner(t)
	handler := integration_common.NewHandler(runner).ClickGeneralTab()
	runner.LoadStateFromFile(path)
	runner.NextFrame()

	// Act
	handler.DragPlayerCountToMaximum()

	// Assert
	assert.Equal(t, 2, runner.CurrentState().PlayerCount)
}

// A file may carry the tournament rule with a count the editor cannot run; the
// panel has to show, and keep writing back, the corrected count.
func TestGeneralTab_ACorrectedTournamentCountSurvivesIdleFrames(t *testing.T) {
	// Arrange
	path := writeGeneralPanelFixture(t, func(state *editor_state_model.EditorState) { state.PlayerCount = 6 })
	enableTournamentInFile(t, path)
	runner := integration_common.NewAppRunner(t)
	integration_common.NewHandler(runner).ClickGeneralTab()
	runner.LoadStateFromFile(path)

	// Act
	for range 5 {
		runner.NextFrame()
	}

	// Assert
	assert.Equal(t, 2, runner.CurrentState().PlayerCount)
}

// The correction is announced once. If the panel kept resubmitting the stale
// count, every idle frame would queue the notice again and overwrite whatever
// the user is reading.
func TestGeneralTab_ACorrectedTournamentCountIsNotReannouncedEveryFrame(t *testing.T) {
	// Arrange
	path := writeGeneralPanelFixture(t, func(state *editor_state_model.EditorState) { state.PlayerCount = 6 })
	enableTournamentInFile(t, path)
	runner := integration_common.NewAppRunner(t)
	integration_common.NewHandler(runner).ClickGeneralTab()
	runner.LoadStateFromFile(path)
	runner.NextFrame()
	runner.SetStatus("idle", false)

	// Act
	for range 5 {
		runner.NextFrame()
	}

	// Assert
	message, _ := runner.Status()
	assert.Equal(t, "idle", message)
}

// Leaving the tournament rules unlocks the control but remembers nothing: the
// count stays where the rules left it.
func TestGeneralTab_LeavingTournamentKeepsTwoPlayers(t *testing.T) {
	// Arrange
	runner := integration_common.NewAppRunner(t)
	handler := integration_common.NewHandler(runner).ClickGeneralTab()
	handler.SelectVictoryCondition(tournamentVictoryLabel)
	handler.DragPlayerCountToMaximum()
	require.Equal(t, 2, runner.CurrentState().PlayerCount, "precondition: the locked slider ignored the drag")

	// Act
	handler.SelectVictoryCondition(standardVictoryLabel)

	// Assert
	assert.Equal(t, 2, runner.CurrentState().PlayerCount)
}

func TestGeneralTab_LeavingTournamentUnlocksTheSlider(t *testing.T) {
	// Arrange
	runner := integration_common.NewAppRunner(t)
	handler := integration_common.NewHandler(runner).ClickGeneralTab()
	handler.SelectVictoryCondition(tournamentVictoryLabel)
	handler.SelectVictoryCondition(standardVictoryLabel)

	// Act
	handler.DragPlayerCountToMaximum()

	// Assert
	assert.Equal(t, 8, runner.CurrentState().PlayerCount)
}

// Choosing the tournament victory condition locks the control on the spot, with
// no file involved.
func TestGeneralTab_ChoosingTheTournamentVictoryLocksTheSlider(t *testing.T) {
	// Arrange
	runner := integration_common.NewAppRunner(t)
	handler := integration_common.NewHandler(runner).ClickGeneralTab()
	handler.DragPlayerCountToMaximum()
	require.Equal(t, 8, runner.CurrentState().PlayerCount, "precondition: the unlocked slider took the drag")

	// Act
	handler.SelectVictoryCondition(tournamentVictoryLabel)

	// Assert
	assert.Equal(t, 2, runner.CurrentState().PlayerCount)
}

// Only one alias has to be on: unchecking the rule while the victory condition
// still selects a tournament leaves the mode, and the lock, in place.
func TestGeneralTab_UncheckingTheRuleUnderTheTournamentVictoryKeepsTheLock(t *testing.T) {
	// Arrange
	runner := integration_common.NewAppRunner(t)
	handler := integration_common.NewHandler(runner).ClickGeneralTab()
	handler.SelectVictoryCondition(tournamentVictoryLabel)
	require.True(t, runner.CurrentState().Tournament, "precondition: the condition switched the rule on")

	// Act
	handler.ToggleConditionRule()
	handler.DragPlayerCountToMaximum()

	// Assert
	assert.Equal(t, 2, runner.CurrentState().PlayerCount)
}

// Only one alias has to be on: unchecking the rule while the victory condition
// still selects a tournament leaves the mode, and the lock, in place.
func TestGeneralTab_TheVictoryConditionAloneKeepsTheSliderLocked(t *testing.T) {
	// Arrange
	path := writeGeneralPanelFixture(t, func(state *editor_state_model.EditorState) {
		state.VictoryCondition = registry.GetWinningConditionValues().Tournament
	})
	disableTournamentCheckboxInFile(t, path)
	runner := integration_common.NewAppRunner(t)
	handler := integration_common.NewHandler(runner).ClickGeneralTab()
	runner.LoadStateFromFile(path)
	runner.NextFrame()

	// Act
	handler.DragPlayerCountToMaximum()

	// Assert
	assert.Equal(t, 2, runner.CurrentState().PlayerCount)
}

func disableTournamentCheckboxInFile(t *testing.T, path string) {
	t.Helper()
	raw, err := os.ReadFile(path)
	require.NoError(t, err)
	document := map[string]any{}
	require.NoError(t, json.Unmarshal(raw, &document))
	document["tournament"] = false
	patched, err := json.Marshal(document)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(path, patched, 0o600))
}
