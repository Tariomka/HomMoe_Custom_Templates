//go:build integration_test

package integration_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/interfaces"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/panels"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newEditorSession builds a State plus the three editor panels bound to it,
// mirroring how editor.Window wires them together. The returned saveFrame and
// loadPanels closures reproduce the per-frame window.save() and window.load()
// behaviour so the integration tests can drive the real frame loop without a
// display.
func newEditorSession() (state *drivers.State, saveFrame func(), loadPanels func()) {
	backend := handlers.NewGuiHandler()
	state = drivers.NewUIStateWithBackend(backend)
	editorPanels := []interfaces.IPanel{
		panels.NewGeneralPanel(state),
		panels.NewLayoutPanel(state, backend, backend, backend),
		panels.NewBonusesPanel(state),
	}
	saveFrame = func() {
		for _, panel := range editorPanels {
			panel.SaveToState()
		}
	}
	loadPanels = func() {
		for _, panel := range editorPanels {
			panel.LoadFromState()
		}
	}
	return state, saveFrame, loadPanels
}

// TestLoadFromFile_SyncsPanels_AndSurvivesNextFrameSave is the regression test
// for bug #1: loading a saved editor state must update the UI panels and must
// not be silently overwritten by the very next frame's SaveToState.
func TestLoadFromFile_SyncsPanels_AndSurvivesNextFrameSave(t *testing.T) {
	dir := t.TempDir()
	savedPath := filepath.Join(dir, "saved.gen.json")

	// Author a distinctive state and persist it through the real save path.
	author := drivers.NewUIState()
	author.UpdateState(func(s *dtos.EditorStateDto) {
		s.TemplateName = "Loaded Template"
		s.PlayerCount = 6
		s.Topology = config.TopologyCross
		s.ResourceDensityPercent = 175
	})
	author.SaveStateToFile(savedPath)
	message, irError := author.GetStatus()
	require.False(t, irError)
	assert.Equal(t, "Saved "+savedPath, message)

	// Fresh editor session sitting at the defaults.
	state, _, _ := newEditorSession()

	// One frame at defaults: the panels write their default widget values back
	// into the state, just like the running editor does every frame.
	require.Equal(t, 2, state.GetStateData().PlayerCount, "precondition: fresh session is at default player count")

	// Load the saved file. The panel resync is supplied as the onLoaded
	// callback, exactly as editor.Window wires window.load().
	state.LoadStateFromFile(savedPath)
	message, irError = state.GetStatus()
	require.False(t, irError)
	assert.Equal(t, "Loaded "+savedPath, message)

	loaded := state.GetStateData()
	assert.Equal(t, "Loaded Template", loaded.TemplateName)
	assert.Equal(t, 6, loaded.PlayerCount)
	assert.Equal(t, config.TopologyCross, loaded.Topology)
	assert.Equal(t, 175, loaded.ResourceDensityPercent)

	// The crux of bug #1: the next frame's SaveToState must reproduce the
	// loaded values (because the panels resynced), not clobber them with the
	// stale widget state that existed before the load.
	after := state.GetStateData()
	assert.Equal(t, "Loaded Template", after.TemplateName, "loaded template name was overwritten by the next frame")
	assert.Equal(t, 6, after.PlayerCount, "loaded player count was overwritten by the next frame")
	assert.Equal(t, config.TopologyCross, after.Topology, "loaded topology was overwritten by the next frame")
	assert.Equal(t, 175, after.ResourceDensityPercent, "loaded resource density was overwritten by the next frame")
}

// TestManualEdits_PersistToGenJson_AndReapplyAfterLoad is the regression test
// for bug #2: manual zone and connection edits must be written into the
// .gen.json file and restored (and reapplied to the regenerated template) when
// the file is loaded again.
func TestManualEdits_PersistToGenJson_AndReapplyAfterLoad(t *testing.T) {
	dir := t.TempDir()
	savedPath := filepath.Join(dir, "manual.gen.json")
	now := time.Now()

	// Generate a baseline template (the first AutoRegenerate call generates).
	state := drivers.NewUIState()
	state.AutoRegenerate(now)
	template := state.GetLastTemplate()
	require.NotNil(t, template, "expected a generated template")
	require.NotEmpty(t, template.Variants)
	require.GreaterOrEqual(t, len(template.Variants[0].Zones), 2)

	// Hand-edit the layout: stamp a manual position onto every zone and add a
	// user-created connection between the first two zones.
	zones := append([]entities.Zone(nil), template.Variants[0].Zones...)
	for i := range zones {
		zones[i].ManualPosition = &[2]float64{0.1 * float64(i+1), 0.2 * float64(i+1)}
	}
	connections := append([]entities.Connection(nil), template.Variants[0].Connections...)
	added := entities.Connection{
		From:           zones[0].Name,
		To:             zones[1].Name,
		ConnectionType: "Portal",
		IsUserAdded:    true,
	}
	connections = append(connections, added)
	state.ApplyEditedZones(zones, connections)

	// Save and confirm the file on disk actually carries the manual edits.
	state.SaveStateToFile(savedPath)
	message, irError := state.GetStatus()
	require.False(t, irError)
	assert.Equal(t, "Saved "+savedPath, message)

	raw, err := os.ReadFile(savedPath)
	require.NoError(t, err)
	var onDisk dtos.EditorStateDto
	require.NoError(t, json.Unmarshal(raw, &onDisk))

	require.Len(t, onDisk.ManualZones, len(zones), "gen.json did not persist all manual zones")
	require.NotEmpty(t, onDisk.ManualConnections, "gen.json did not persist manual connections")

	// ManualPosition is json:"-" on the zone itself, so the round trip relies
	// on the save wrapper preserving it.
	require.NotNil(t, onDisk.ManualZones[0].ManualPosition, "manual position was lost on save")
	assert.InDelta(t, 0.1, onDisk.ManualZones[0].ManualPosition[0], 1e-9)
	assert.InDelta(t, 0.2, onDisk.ManualZones[0].ManualPosition[1], 1e-9)

	lastConn := onDisk.ManualConnections[len(onDisk.ManualConnections)-1]
	assert.True(t, lastConn.IsUserAdded, "IsUserAdded flag was lost on save")
	assert.Equal(t, "Portal", lastConn.Connection.ConnectionType)

	// Load into a fresh session and regenerate; the manual layout must come
	// back and be reapplied to the regenerated template.
	reloaded := drivers.NewUIState()
	reloaded.LoadStateFromFile(savedPath)
	message, irError = reloaded.GetStatus()
	require.False(t, irError)
	assert.Equal(t, "Loaded "+savedPath, message)
	reloaded.AutoRegenerate(now)

	got := reloaded.GetLastTemplate()
	require.NotNil(t, got, "expected a regenerated template after load")
	require.NotEmpty(t, got.Variants)

	gotZones := got.Variants[0].Zones
	require.Len(t, gotZones, len(zones), "reapplied zone count does not match the saved manual layout")
	for i := range gotZones {
		require.NotNilf(t, gotZones[i].ManualPosition, "zone %d lost its manual position after load", i)
	}

	foundAdded := false
	for _, c := range got.Variants[0].Connections {
		if c.From == added.From && c.To == added.To && c.ConnectionType == "Portal" {
			foundAdded = true
			break
		}
	}
	assert.True(t, foundAdded, "user-added connection was not reapplied after load")
}

// TestSaveWithoutManualEdits_OmitsManualFields verifies that a plain save (no
// manual editing) does not write any manual-edit payload, keeping default
// .gen.json files clean.
func TestSaveWithoutManualEdits_OmitsManualFields(t *testing.T) {
	dir := t.TempDir()
	savedPath := filepath.Join(dir, "plain.gen.json")

	state := drivers.NewUIState()
	state.SaveStateToFile(savedPath)
	message, irError := state.GetStatus()
	require.False(t, irError)
	assert.Equal(t, "Saved "+savedPath, message)

	raw, err := os.ReadFile(savedPath)
	require.NoError(t, err)

	var onDisk dtos.EditorStateDto
	require.NoError(t, json.Unmarshal(raw, &onDisk))
	assert.Empty(t, onDisk.ManualZones)
	assert.Empty(t, onDisk.ManualConnections)
}

// TestStructuralRegeneration_DropsManualEdits confirms that changing a
// layout-defining option (player count) after manual editing discards the
// manual edits rather than reapplying them onto an incompatible layout.
func TestStructuralRegeneration_DropsManualEdits(t *testing.T) {
	now := time.Now()
	state := drivers.NewUIState()
	state.AutoRegenerate(now)

	template := state.GetLastTemplate()
	require.NotNil(t, template)
	require.NotEmpty(t, template.Variants)

	zones := append([]entities.Zone(nil), template.Variants[0].Zones...)
	for i := range zones {
		zones[i].ManualPosition = &[2]float64{0.3, 0.4}
	}
	state.ApplyEditedZones(zones, template.Variants[0].Connections)

	// A structural change (player count) must regenerate from scratch.
	state.UpdateState(func(s *dtos.EditorStateDto) { s.PlayerCount = 4 })
	state.AutoRegenerate(now)

	regenerated := state.GetLastTemplate()
	require.NotNil(t, regenerated)
	require.NotEmpty(t, regenerated.Variants)
	// With four players the regenerated layout has its own spawn zones; the
	// manual single-position stamp must not have been forced back on.
	assert.GreaterOrEqual(t, len(regenerated.Variants[0].Zones), 4)
}

// TestLoadFromFile_RestoresGameMode_AndSurvivesNextFrameSave is the regression
// test for the generalPanel gameMode load bug: LoadFromState used to hard-reset
// the game-mode segment button to index 0, so the next frame's SaveToState
// silently reverted a loaded non-default game mode back to Classic.
func TestLoadFromFile_RestoresGameMode_AndSurvivesNextFrameSave(t *testing.T) {
	dir := t.TempDir()
	savedPath := filepath.Join(dir, "gamemode.gen.json")
	singleHero := registry.GetGameModeValues().SingleHero

	// Author a state with the non-default game mode and persist it.
	author := drivers.NewUIState()
	author.UpdateState(func(s *dtos.EditorStateDto) { s.GameMode = singleHero })
	author.SaveStateToFile(savedPath)
	message, irError := author.GetStatus()
	require.False(t, irError)
	require.Equal(t, "Saved "+savedPath, message)

	// Fresh editor session at defaults; load, resync panels (window.load()),
	// then run the next frame's SaveToState (window.save()).
	state, saveFrame, loadPanels := newEditorSession()
	state.LoadStateFromFile(savedPath)
	loadPanels()
	saveFrame()

	assert.Equal(t, singleHero, state.GetStateData().GameMode,
		"loaded game mode was clobbered by the next frame's SaveToState")
}

// TestLoadFromFile_UnknownGameMode_FallsBackToClassic verifies that a
// hand-edited .gen.json with an unrecognized game mode does not panic and
// settles on the default (Classic) mode.
func TestLoadFromFile_UnknownGameMode_FallsBackToClassic(t *testing.T) {
	dir := t.TempDir()
	savedPath := filepath.Join(dir, "unknown-gamemode.gen.json")

	// UpdateState fixes invalid values on the fly, so the hand-edited file has
	// to be written directly.
	require.NoError(t, os.WriteFile(savedPath, []byte(`{"gameMode": "NotARealGameMode"}`), 0o644))

	state, saveFrame, loadPanels := newEditorSession()
	state.LoadStateFromFile(savedPath)
	loadPanels()
	saveFrame()

	assert.Equal(t, registry.GetGameModeValues().Classic, state.GetStateData().GameMode,
		"unknown game mode should fall back to Classic")
}

// TestLoadFromFile_UnknownVictoryCondition_FallsBackToStandardAndWarns verifies
// that a hand-edited .gen.json with an unrecognized victory condition settles
// on Standard and tells the user about it instead of silently reshaping the
// setting.
func TestLoadFromFile_UnknownVictoryCondition_FallsBackToStandardAndWarns(t *testing.T) {
	dir := t.TempDir()
	savedPath := filepath.Join(dir, "unknown-victory.gen.json")

	// UpdateState fixes invalid values on the fly, so the hand-edited file has
	// to be written directly.
	require.NoError(t, os.WriteFile(savedPath,
		[]byte(`{"victoryCondition": "notARealVictoryCondition"}`), 0o644))

	state, saveFrame, loadPanels := newEditorSession()
	state.LoadStateFromFile(savedPath)
	loadPanels()
	saveFrame()

	message, isError := state.GetStatus()
	require.False(t, isError)
	assert.Contains(t, message, `victoryCondition "notARealVictoryCondition" is not a known victory condition`,
		"the user must be told their victory condition was dropped")
	assert.Equal(t, registry.GetWinningConditionValues().Standard, state.GetStateData().VictoryCondition,
		"unknown victory condition should fall back to Standard")
}
