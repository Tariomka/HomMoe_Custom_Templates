//go:build integration_test

package integration_test

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/internal/composition"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/zone_helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// writeTournamentFixture writes a real .gen.json carrying a hand-made layout at
// the given player count, then flips the tournament rule on in the file itself.
// The editor corrects such a count on the way in, so the only way to produce the
// file the user may already own is to patch it after saving.
func writeTournamentFixture(t *testing.T, playerCount int) string {
	t.Helper()
	directory := t.TempDir()
	author := newUIState()
	author.UpdateState(func(state *editor_state_model.EditorState) {
		state.TemplateName = "Tournament Fixture"
		state.PlayerCount = playerCount
	})
	author.AutoRegenerate(time.Now())
	template := author.GetLastTemplate()
	require.NotNil(t, template, "the fixture needs a generated layout to hand-edit")
	require.NotEmpty(t, template.Variants)

	zones := append([]template_model.Zone(nil), template.Variants[0].Zones...)
	zones[0].ManualPosition = new(data.NewVec2(0.3, 0.4))
	author.ApplyEditedZones(dtos.ZoneEditorZonesDto{
		Zones:       zones,
		Connections: append([]template_model.Connection(nil), template.Variants[0].Connections...),
	})
	authored := author.GetStateData()
	require.True(t, authored.HasManualEdits(), "the fixture must carry a manual layout")

	author.SaveStateToFile(filepath.Join(directory, "fixture.gen.json"))
	written := filepath.Join(directory, "Tournament Fixture.gen.json")
	enableTournamentInFile(t, written)
	return written
}

func enableTournamentInFile(t *testing.T, path string) {
	t.Helper()
	patchStateFile(t, path, func(document map[string]any) { document["tournament"] = true })
}

// selectTournamentVictoryInFile switches the file to the other tournament
// alias: the rule is off, the victory condition alone turns it on.
func selectTournamentVictoryInFile(t *testing.T, path string) {
	t.Helper()
	patchStateFile(t, path, func(document map[string]any) {
		document["tournament"] = false
		document["victoryCondition"] = registry.GetWinningConditionValues().Tournament
	})
}

func patchStateFile(t *testing.T, path string, patch func(document map[string]any)) {
	t.Helper()
	raw, err := os.ReadFile(path)
	require.NoError(t, err)
	document := map[string]any{}
	require.NoError(t, json.Unmarshal(raw, &document))
	patch(document)
	patched, err := json.Marshal(document)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(path, patched, 0o600))
}

// loadTournamentFixture loads a patched fixture into a fresh editor session.
func loadTournamentFixture(t *testing.T, playerCount int) (*drivers.State, string) {
	t.Helper()
	path := writeTournamentFixture(t, playerCount)
	state := newUIState()
	state.LoadStateFromFile(path)
	return state, path
}

func TestWhenATournamentFileHasTooManyPlayers_TheLoadedCountIsTwo(t *testing.T) {
	// Arrange
	state, _ := loadTournamentFixture(t, 6)

	// Act
	loaded := state.GetStateData()

	// Assert
	assert.Equal(t, 2, loaded.PlayerCount)
}

func TestWhenATournamentFileHasTooManyPlayers_TheManualLayoutIsDiscarded(t *testing.T) {
	// Arrange
	state, _ := loadTournamentFixture(t, 6)

	// Act
	loaded := state.GetStateData()

	// Assert
	assert.False(t, loaded.HasManualEdits())
}

func TestWhenATournamentFileHasTooManyPlayers_TheDocumentIsUnsaved(t *testing.T) {
	// Arrange
	state, _ := loadTournamentFixture(t, 6)

	// Assert
	assert.True(t, state.IsUnsaved())
}

func TestWhenATournamentFileHasTooManyPlayers_ExitAsksForConfirmation(t *testing.T) {
	// Arrange
	state, _ := loadTournamentFixture(t, 6)
	exitCalled := false
	state.SetOnExit(func() { exitCalled = true })

	// Act
	state.Exit()

	// Assert
	assert.False(t, exitCalled)
}

func TestWhenATournamentFileHasTooManyPlayers_TheStatusReportsTheCorrection(t *testing.T) {
	// Arrange
	state, _ := loadTournamentFixture(t, 6)

	// Act
	message, isError := state.GetStatus()

	// Assert
	require.False(t, isError)
	assert.Contains(t, message, "Tournament mode requires 2 players")
}

// Saving stays the user's decision: the file still holds what it held.
func TestWhenATournamentFileHasTooManyPlayers_TheFileOnDiskIsUnchanged(t *testing.T) {
	// Arrange
	path := writeTournamentFixture(t, 6)
	before, err := os.ReadFile(path)
	require.NoError(t, err)
	state := newUIState()

	// Act
	state.LoadStateFromFile(path)

	// Assert
	after, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, before, after)
}

// The load is followed by the first automatic regeneration, which writes its own
// status message; the correction has to still be readable afterwards.
func TestWhenTheFirstRegenerationFollowsACorrectedLoad_TheNoticeIsStillVisible(t *testing.T) {
	// Arrange
	state, _ := loadTournamentFixture(t, 6)

	// Act
	state.AutoRegenerate(time.Now())

	// Assert
	message, isError := state.GetStatus()
	require.False(t, isError)
	assert.Contains(t, message, "Tournament mode requires 2 players")
}

// The notice belongs to the document that earned it, not to the next one.
func TestWhenAnotherDocumentReplacesACorrectedLoad_TheNoticeIsGone(t *testing.T) {
	// Arrange
	state, _ := loadTournamentFixture(t, 6)

	// Act
	state.Reset()
	state.AutoRegenerate(time.Now())

	// Assert
	message, _ := state.GetStatus()
	assert.NotContains(t, message, "Tournament mode requires 2 players")
}

func TestWhenTheCorrectedStateIsSavedAndReloaded_ItLoadsClean(t *testing.T) {
	// Arrange
	state, path := loadTournamentFixture(t, 6)
	state.SaveStateToFile(path)
	require.False(t, state.IsUnsaved(), "precondition: saving clears the correction's dirty flag")
	reloaded := newUIState()

	// Act
	reloaded.LoadStateFromFile(filepath.Join(filepath.Dir(path), "Tournament Fixture.gen.json"))

	// Assert
	assert.False(t, reloaded.IsUnsaved())
}

func TestWhenTheCorrectedStateIsSavedAndReloaded_ItStaysAtTwoPlayers(t *testing.T) {
	// Arrange
	state, path := loadTournamentFixture(t, 6)
	state.SaveStateToFile(path)
	reloaded := newUIState()

	// Act
	reloaded.LoadStateFromFile(filepath.Join(filepath.Dir(path), "Tournament Fixture.gen.json"))

	// Assert
	assert.Equal(t, 2, reloaded.GetStateData().PlayerCount)
}

func TestWhenAValidTwoPlayerTournamentIsLoaded_TheManualLayoutSurvives(t *testing.T) {
	// Arrange
	state, _ := loadTournamentFixture(t, 2)

	// Act
	loaded := state.GetStateData()

	// Assert
	assert.True(t, loaded.HasManualEdits())
}

func TestWhenAValidTwoPlayerTournamentIsLoaded_TheDocumentStaysSaved(t *testing.T) {
	// Arrange
	state, _ := loadTournamentFixture(t, 2)

	// Assert
	assert.False(t, state.IsUnsaved())
}

// The victory condition is the other way into the tournament rules, and a file
// carrying it has to be corrected exactly like one carrying the checkbox.
func TestWhenATournamentVictoryFileHasTooManyPlayers_TheLoadedCountIsTwo(t *testing.T) {
	// Arrange
	path := writeTournamentFixture(t, 6)
	selectTournamentVictoryInFile(t, path)
	state := newUIState()

	// Act
	state.LoadStateFromFile(path)

	// Assert
	assert.Equal(t, 2, state.GetStateData().PlayerCount)
}

func TestWhenATournamentVictoryFileHasTooManyPlayers_TheManualLayoutIsDiscarded(t *testing.T) {
	// Arrange
	path := writeTournamentFixture(t, 6)
	selectTournamentVictoryInFile(t, path)
	state := newUIState()

	// Act
	state.LoadStateFromFile(path)

	// Assert
	loaded := state.GetStateData()
	assert.False(t, loaded.HasManualEdits())
}

func TestWhenATournamentVictoryFileHasTooManyPlayers_TheDocumentIsUnsaved(t *testing.T) {
	// Arrange
	path := writeTournamentFixture(t, 6)
	selectTournamentVictoryInFile(t, path)
	state := newUIState()

	// Act
	state.LoadStateFromFile(path)

	// Assert
	assert.True(t, state.IsUnsaved())
}

func TestWhenAValidTwoPlayerTournamentVictoryIsLoaded_TheManualLayoutSurvives(t *testing.T) {
	// Arrange
	path := writeTournamentFixture(t, 2)
	selectTournamentVictoryInFile(t, path)
	state := newUIState()

	// Act
	state.LoadStateFromFile(path)

	// Assert
	loaded := state.GetStateData()
	assert.True(t, loaded.HasManualEdits())
}

func TestWhenALoadFails_TheCurrentDocumentIsUnchanged(t *testing.T) {
	// Arrange
	state, _ := loadTournamentFixture(t, 2)
	before := state.GetStateData()

	// Act
	state.LoadStateFromFile(filepath.Join(t.TempDir(), "missing.gen.json"))

	// Assert
	assert.Equal(t, before, state.GetStateData())
}

// The live editor flow: a hand-made layout, then a mode change, then the map
// that mode generates instead.
func TestWhenAModeChangeFollowsAManualLayout_TheGeneratedTemplateReplacesIt(t *testing.T) {
	// Arrange
	state := newEditedSessionAtDefaults(t)

	// Act
	state.UpdateState(func(current *editor_state_model.EditorState) { current.GladiatorArena = true })
	state.AutoRegenerate(time.Now())

	// Assert
	regenerated := state.GetLastTemplate()
	require.NotNil(t, regenerated)
	require.NotEmpty(t, regenerated.Variants)
	for _, zone := range regenerated.Variants[0].Zones {
		require.Nilf(t, zone.ManualPosition, "zone %s kept a position from the discarded layout", zone.Name)
	}
}

func TestWhenAModeChangeFollowsAManualLayout_TheSnapshotIsDropped(t *testing.T) {
	// Arrange
	state := newEditedSessionAtDefaults(t)

	// Act
	state.UpdateState(func(current *editor_state_model.EditorState) {
		current.VictoryCondition = registry.GetWinningConditionValues().Tournament
	})

	// Assert
	edited := state.GetStateData()
	assert.False(t, edited.HasManualEdits())
}

// Changing the mode's own timing options is not a mode change, so the layout
// the user built by hand stays.
func TestWhenOnlyTournamentTimingChanges_TheManualLayoutSurvives(t *testing.T) {
	// Arrange
	state := newEditedSessionAtDefaults(t)

	// Act
	state.UpdateState(func(current *editor_state_model.EditorState) { current.TournamentInterval += 3 })

	// Assert
	edited := state.GetStateData()
	assert.True(t, edited.HasManualEdits())
}

// The application generator validates before it maps, so a tournament state that
// asks for more players still produces a two-player map.
func TestWhenATournamentStateIsGenerated_TheTemplateHasTwoPlayerZones(t *testing.T) {
	for _, testCase := range tournamentAliasCases() {
		t.Run(testCase.subtestName+"_TheTemplateHasTwoPlayerZones", func(t *testing.T) {
			// Arrange
			state := newOversizedTournamentState(testCase.mutate)

			// Act
			generated := generateTemplateThroughHandler(t, state)

			// Assert
			assert.Equal(t, 2, countPlayerZones(generated.Template.Variants[0].Zones))
		})
	}
}

// The exported player labels come from the player index the generator counted
// to, so they are the second place a disagreement about the count would show up.
func TestWhenATournamentStateIsGenerated_ThePlayerLabelsStopAtTwoPlayers(t *testing.T) {
	for _, testCase := range tournamentAliasCases() {
		t.Run(testCase.subtestName+"_ThePlayerLabelsStopAtTwoPlayers", func(t *testing.T) {
			// Arrange
			state := newOversizedTournamentState(testCase.mutate)

			// Act
			generated := generateTemplateThroughHandler(t, state)

			// Assert
			assert.Equal(t, []string{"Player1", "Player2"}, playerLabels(generated.Template.Variants[0].Zones))
		})
	}
}

// The rules the template carries have to agree with the layout it carries.
func TestWhenATournamentStateIsGenerated_TheTournamentRulesAreExported(t *testing.T) {
	for _, testCase := range tournamentAliasCases() {
		t.Run(testCase.subtestName+"_TheTournamentRulesAreExported", func(t *testing.T) {
			// Arrange
			state := newOversizedTournamentState(testCase.mutate)

			// Act
			generated := generateTemplateThroughHandler(t, state)

			// Assert
			assert.True(t, generated.Template.GameRules.WinConditions.Tournament)
		})
	}
}

// The handler corrects a clone. The caller still owns the state it passed in,
// and gets to decide what to do about it.
func TestWhenATournamentStateIsGenerated_TheCallersStateIsUntouched(t *testing.T) {
	for _, testCase := range tournamentAliasCases() {
		t.Run(testCase.subtestName+"_TheCallersStateIsUntouched", func(t *testing.T) {
			// Arrange
			state := newOversizedTournamentState(testCase.mutate)

			// Act
			generateTemplateThroughHandler(t, state)

			// Assert
			assert.Equal(t, 6, state.PlayerCount)
		})
	}
}

// The correction is reported rather than performed silently.
func TestWhenATournamentStateIsGenerated_TheCorrectionIsWarnedAbout(t *testing.T) {
	// Arrange
	state := newOversizedTournamentState(func(state *editor_state_model.EditorState) { state.Tournament = true })

	// Act
	generated := generateTemplateThroughHandler(t, state)

	// Assert
	assert.Contains(t, strings.Join(generated.Warnings, "; "), "playerCount 6")
}

type tournamentAliasCase struct {
	subtestName string
	mutate      func(state *editor_state_model.EditorState)
}

// tournamentAliasCases returns the two ways a state asks for the tournament
// rules, so each invariant below is proven through both of them.
func tournamentAliasCases() []tournamentAliasCase {
	return []tournamentAliasCase{
		{"WhenTheTournamentCheckboxIsSet", func(state *editor_state_model.EditorState) { state.Tournament = true }},
		{
			"WhenTheTournamentVictoryConditionIsChosen",
			func(state *editor_state_model.EditorState) {
				state.VictoryCondition = registry.GetWinningConditionValues().Tournament
			},
		},
	}
}

func newOversizedTournamentState(
	mutate func(state *editor_state_model.EditorState)) editor_state_model.EditorState {
	state := editor_state_model.NewDefaultEditorStateModel()
	mutate(&state)
	state.PlayerCount = 6
	return state
}

func generateTemplateThroughHandler(
	t *testing.T, state editor_state_model.EditorState) dtos.TemplateLoadDto {
	t.Helper()
	generated, err := composition.InitializeGuiHandler().
		GenerateTemplate(editor_state_dto.EditorStateDto{EditorState: state})
	require.NoError(t, err)
	require.NotNil(t, generated.Template)
	require.NotEmpty(t, generated.Template.Variants)
	return generated
}

// playerLabels returns every distinct player the generated zones name, sorted.
// Spawns carry the label, and a player-owned castle repeats it as its owner;
// both are exported, so both are collected.
func playerLabels(zones []template_model.Zone) []string {
	labels := []string{}
	for _, zone := range zones {
		for _, mainObject := range zone.MainObjects {
			for _, label := range []string{mainObject.Spawn, mainObject.Owner} {
				if label != "" && !slices.Contains(labels, label) {
					labels = append(labels, label)
				}
			}
		}
	}
	slices.Sort(labels)
	return labels
}

func countPlayerZones(zones []template_model.Zone) int {
	playerZones := 0
	for _, zone := range zones {
		if zone_helpers.IsZoneNamePlayer(zone.Name) {
			playerZones++
		}
	}
	return playerZones
}

// newEditedSessionAtDefaults returns a generated session carrying a hand-made
// layout, ready for the mode change that must discard it.
func newEditedSessionAtDefaults(t *testing.T) *drivers.State {
	t.Helper()
	state := newUIState()
	state.AutoRegenerate(time.Now())
	template := state.GetLastTemplate()
	require.NotNil(t, template)
	require.NotEmpty(t, template.Variants)

	zones := append([]template_model.Zone(nil), template.Variants[0].Zones...)
	zones[0].ManualPosition = new(data.NewVec2(0.25, 0.75))
	state.ApplyEditedZones(dtos.ZoneEditorZonesDto{
		Zones:       zones,
		Connections: append([]template_model.Connection(nil), template.Variants[0].Connections...),
	})
	edited := state.GetStateData()
	require.True(t, edited.HasManualEdits(), "precondition: the session must carry a manual layout")
	return state
}

// tournamentCorrectionNotice is the part of the correction message that names
// the rule, shared by the one-line and the two-line wording.
const tournamentCorrectionNotice = "Tournament mode requires 2 players"

// failingGenerationHandler is the real handler with a switch on its generation.
// Everything else is delegated, so the driver runs against production code and
// only the one failure the test is about is injected.
type failingGenerationHandler struct {
	handler_interfaces.IGuiHandler

	failGeneration bool
}

func (this *failingGenerationHandler) GenerateTemplate(
	state editor_state_dto.EditorStateDto) (dtos.TemplateLoadDto, error) {
	if this.failGeneration {
		return dtos.TemplateLoadDto{}, errors.New("generator unavailable")
	}

	return this.IGuiHandler.GenerateTemplate(state)
}

// newFailedRegenerationSession loads a corrected file into a session whose
// generator is broken, so the editor ends up showing the generation failure
// while the load's correction notice is still owed to the next good status.
func newFailedRegenerationSession(t *testing.T) (*drivers.State, *failingGenerationHandler) {
	t.Helper()
	handler := &failingGenerationHandler{IGuiHandler: composition.InitializeGuiHandler()}
	state := drivers.NewUIState(
		handler,
		composition.InitializeFileSystemHandler(),
		composition.InitializeRegenerationHandler(),
		true)

	state.LoadStateFromFile(writeTournamentFixture(t, 6))
	loadMessage, loadFailed := state.GetStatus()
	require.False(t, loadFailed, "precondition: the corrected load is not a failure")
	require.Contains(t, loadMessage, tournamentCorrectionNotice, "precondition: the load owes a notice")

	handler.failGeneration = true
	state.AutoRegenerate(time.Now())
	failureMessage, generationFailed := state.GetStatus()
	require.True(t, generationFailed, "precondition: the first regeneration must fail")
	require.Contains(t, failureMessage, "Generation failed")
	return state, handler
}

// idleFrameCount is enough frames to prove the status is not being rewritten
// per frame rather than once.
const idleFrameCount = 5

// runIdleFrames reproduces the window's per-frame loop while the user does
// nothing: every panel writes the unchanged state back and the frame asks for a
// regeneration that is not due. Such a frame must leave the status alone.
func runIdleFrames(state *drivers.State) {
	frameTime := time.Now()
	for frame := range idleFrameCount {
		state.UpdateState(func(*editor_state_model.EditorState) {})
		state.AutoRegenerate(frameTime.Add(time.Duration(frame) * 16 * time.Millisecond))
	}
}

// saveCorrectedSession saves the loaded document back over the file it came
// from, which is the user's answer to the correction.
func saveCorrectedSession(t *testing.T, state *drivers.State) {
	t.Helper()
	state.SaveStateToFile(state.GetCurrentPath())
	savedMessage, saveFailed := state.GetStatus()
	require.False(t, saveFailed, "precondition: the save must succeed")
	require.Contains(t, savedMessage, "Saved ")
}

// A failed generation leaves the notice pending, and the idle frames that
// follow used to republish it over the failure the user still has to read.
func TestWhenIdleFramesFollowAFailedRegeneration_TheFailureIsStillReported(t *testing.T) {
	// Arrange
	state, _ := newFailedRegenerationSession(t)

	// Act
	runIdleFrames(state)

	// Assert
	message, _ := state.GetStatus()
	assert.Contains(t, message, "Generation failed")
}

func TestWhenIdleFramesFollowAFailedRegeneration_TheStatusStaysAnError(t *testing.T) {
	// Arrange
	state, _ := newFailedRegenerationSession(t)

	// Act
	runIdleFrames(state)

	// Assert
	_, isError := state.GetStatus()
	assert.True(t, isError)
}

func TestWhenIdleFramesFollowASaveAfterAFailedRegeneration_TheSavedStatusRemains(t *testing.T) {
	// Arrange
	state, _ := newFailedRegenerationSession(t)
	saveCorrectedSession(t, state)

	// Act
	runIdleFrames(state)

	// Assert
	message, _ := state.GetStatus()
	assert.Contains(t, message, "Saved ")
}

func TestWhenIdleFramesFollowASaveAfterAFailedRegeneration_TheStatusIsNotAnError(t *testing.T) {
	// Arrange
	state, _ := newFailedRegenerationSession(t)
	saveCorrectedSession(t, state)

	// Act
	runIdleFrames(state)

	// Assert
	_, isError := state.GetStatus()
	assert.False(t, isError)
}

// The notice survives the failure and the frames in between: the first
// generation that succeeds is the one that finally carries it.
func TestWhenAGenerationSucceedsAfterAFailureAndASave_TheCorrectionIsCarried(t *testing.T) {
	// Arrange
	state := newRecoveredRegenerationSession(t)

	// Act
	message, isError := state.GetStatus()

	// Assert
	require.False(t, isError)
	assert.Contains(t, message, tournamentCorrectionNotice)
}

func TestWhenAGenerationSucceedsAfterAFailureAndASave_TheDiscardedLayoutIsCarried(t *testing.T) {
	// Arrange
	state := newRecoveredRegenerationSession(t)

	// Act
	message, _ := state.GetStatus()

	// Assert
	assert.Contains(t, message, "the manual zone layout was discarded")
}

// Carried once: the idle frames after the recovery show the generation status,
// not an empty one left behind by a notice that has already been spent.
func TestWhenIdleFramesFollowTheRecoveredGeneration_TheGenerationStatusRemains(t *testing.T) {
	// Arrange
	state := newRecoveredRegenerationSession(t)
	before, _ := state.GetStatus()

	// Act
	runIdleFrames(state)

	// Assert
	message, _ := state.GetStatus()
	assert.Equal(t, before, message)
}

// newRecoveredRegenerationSession walks the whole story: a corrected load, a
// failed regeneration, idle frames, an explicit save, more idle frames, and
// finally the retry that succeeds.
func newRecoveredRegenerationSession(t *testing.T) *drivers.State {
	t.Helper()
	state, handler := newFailedRegenerationSession(t)
	runIdleFrames(state)
	saveCorrectedSession(t, state)
	runIdleFrames(state)

	handler.failGeneration = false
	state.UpdateState(func(current *editor_state_model.EditorState) {
		current.GenerateRoads = !current.GenerateRoads
	})
	state.AutoRegenerate(time.Now())
	require.NotNil(t, state.GetLastTemplate(), "precondition: the retry must generate a template")
	return state
}
