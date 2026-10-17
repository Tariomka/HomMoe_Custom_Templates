//go:build integration_test

package integration_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_generator"
	"github.com/Tariomka/hommoe_custom_templates/internal/composition"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	legacyMoveFixture    = "editorState_v2_square2p_move.gen.json"
	legacyQualityFixture = "editorState_v2_square2p_quality.gen.json"
)

func TestWhenALegacyFileWithEditsIsLoaded_TheFirstRegenerationOpensTheConversionNotice(t *testing.T) {
	// Arrange
	state := loadFixtureSession(t, legacyMoveFixture)

	// Act
	state.AutoRegenerate(time.Now())

	// Assert
	assert.True(t, state.GetDialogHost().IsOpen())
}

func TestWhenALegacyFileWithEditsIsConverted_TheDocumentBecomesUnsaved(t *testing.T) {
	// Arrange
	state := loadFixtureSession(t, legacyMoveFixture)

	// Act
	state.AutoRegenerate(time.Now())

	// Assert
	assert.True(t, state.IsUnsaved())
}

func TestWhenALegacyFileWithEditsIsConverted_NothingIsLeftToConvert(t *testing.T) {
	// Arrange
	state := loadFixtureSession(t, legacyMoveFixture)

	// Act
	state.AutoRegenerate(time.Now())

	// Assert
	stateData := state.GetStateData()
	assert.False(t, stateData.NeedsDerivation())
}

func TestWhenALegacyFileWithEditsIsConverted_TheLayoutIsTheSnapshot(t *testing.T) {
	// Arrange
	state := loadFixtureSession(t, legacyMoveFixture)
	snapshotZones := legacySnapshotZones(t, legacyMoveFixture)

	// Act
	state.AutoRegenerate(time.Now())

	// Assert
	assert.Equal(t, snapshotZones, state.GetLastTemplate().Variants[0].Zones)
}

func TestWhenTheConvertedLayoutIsSavedAndReloaded_TheLayoutIsTheSnapshot(t *testing.T) {
	// Arrange
	state := loadFixtureSession(t, legacyMoveFixture)
	snapshotZones := legacySnapshotZones(t, legacyMoveFixture)
	state.AutoRegenerate(time.Now())
	reloaded := saveAndReload(t, state)

	// Act
	reloaded.AutoRegenerate(time.Now())

	// Assert
	assert.Equal(t, snapshotZones, reloaded.GetLastTemplate().Variants[0].Zones)
}

func TestWhenALegacyFileIsSavedBeforeAnyRegeneration_TheReloadStillConverts(t *testing.T) {
	// Arrange
	state := loadFixtureSession(t, legacyMoveFixture)
	snapshotZones := legacySnapshotZones(t, legacyMoveFixture)
	reloaded := saveAndReload(t, state)

	// Act
	reloaded.AutoRegenerate(time.Now())

	// Assert
	assert.Equal(t, snapshotZones, reloaded.GetLastTemplate().Variants[0].Zones)
}

func TestWhenALegacyFileIsSavedBeforeAnyRegeneration_TheReloadIsStillPending(t *testing.T) {
	// Arrange
	state := loadFixtureSession(t, legacyMoveFixture)

	// Act
	reloaded := saveAndReload(t, state)

	// Assert
	stateData := reloaded.GetStateData()
	assert.True(t, stateData.NeedsDerivation())
}

func TestWhenAPendingReloadRegenerates_ItIsConverted(t *testing.T) {
	// Arrange
	reloaded := saveAndReload(t, loadFixtureSession(t, legacyMoveFixture))

	// Act
	reloaded.AutoRegenerate(time.Now())

	// Assert
	stateData := reloaded.GetStateData()
	assert.Equal(t, []any{false, common_generator.GeneratorOutputVersion},
		[]any{stateData.NeedsDerivation(), stateData.GeneratorVersion})
}

func TestWhenALegacyFileWithoutEditsIsLoaded_NoConversionNoticeOpens(t *testing.T) {
	// Arrange
	state := loadSessionFrom(t, legacyFileWithoutEdits(t))

	// Act
	state.AutoRegenerate(time.Now())

	// Assert
	assert.False(t, state.GetDialogHost().IsOpen())
}

func TestWhenALegacyFileWithoutEditsIsLoaded_TheDocumentStaysSaved(t *testing.T) {
	// Arrange
	state := loadSessionFrom(t, legacyFileWithoutEdits(t))

	// Act
	state.AutoRegenerate(time.Now())

	// Assert
	assert.False(t, state.IsUnsaved())
}

// A setting changed before the first regeneration must not pin the snapshot's
// old value: the conversion compares against the settings the file was saved
// with, so the untouched value follows the new setting.
func TestWhenADensityChangesBeforeTheFirstRegeneration_UntouchedValuesFollowIt(t *testing.T) {
	// Arrange
	state := loadFixtureSession(t, legacyMoveFixture)
	state.UpdateState(func(settings *editor_state_model.EditorState) { settings.ResourceDensityPercent = 180 })
	expected := freshLayoutOf(t, state)

	// Act
	state.AutoRegenerate(time.Now())

	// Assert
	assert.Equal(t, resourceValues(expected.Zones), resourceValues(state.GetLastTemplate().Variants[0].Zones))
}

// Saving while the conversion is pending loses the load-time settings, so the
// reloaded settings become the baseline: values that differ from them stay
// pinned (lossless and conservative).
func TestWhenADensityChangeIsSavedBeforeConversion_UntouchedValuesStayPinned(t *testing.T) {
	// Arrange
	state := loadFixtureSession(t, legacyMoveFixture)
	snapshotZones := legacySnapshotZones(t, legacyMoveFixture)
	state.UpdateState(func(settings *editor_state_model.EditorState) { settings.ResourceDensityPercent = 180 })
	reloaded := saveAndReload(t, state)

	// Act
	reloaded.AutoRegenerate(time.Now())

	// Assert
	assert.Equal(t, resourceValues(snapshotZones), resourceValues(reloaded.GetLastTemplate().Variants[0].Zones))
}

// Castle options changed between loading and the first regeneration reach the
// zones that own their castle count.
func TestWhenCastlesChangeBeforeTheFirstRegeneration_ReprofiledZonesGetTheNewCount(t *testing.T) {
	// Arrange
	state := loadFixtureSession(t, legacyQualityFixture)
	state.UpdateState(func(settings *editor_state_model.EditorState) {
		settings.AdvancedMode = false
		settings.NeutralZoneCastles = 3
	})

	// Act
	state.AutoRegenerate(time.Now())

	// Assert
	assert.Equal(t, []int{3}, reprofiledCastleCounts(t, state))
}

func TestWhenAFileRecordsAnotherGeneratorVersion_TheStatusSaysDerivedValuesMayDiffer(t *testing.T) {
	// Arrange
	state := loadFixtureSession(t, legacyMoveFixture)
	state.AutoRegenerate(time.Now())
	state.SaveStateToFile(filepath.Join(t.TempDir(), "converted.gen.json"))
	rewriteGeneratorVersion(t, state.GetCurrentPath(), 7)
	reloaded := newUIState()
	reloaded.LoadStateFromFile(state.GetCurrentPath())

	// Act
	reloaded.AutoRegenerate(time.Now())

	// Assert
	status, _ := reloaded.GetStatus()
	assert.Contains(t, status, "Saved with generator version 7")
}

func loadFixtureSession(t *testing.T, fixture string) *drivers.State {
	t.Helper()

	return loadSessionFrom(t, filepath.Join("..", "test_helpers", "testdata", fixture))
}

func loadSessionFrom(t *testing.T, path string) *drivers.State {
	t.Helper()
	state := newUIState()
	state.LoadStateFromFile(path)
	status, isError := state.GetStatus()
	require.False(t, isError, status)

	return state
}

// legacyFileWithoutEdits writes the move fixture without its manual section.
func legacyFileWithoutEdits(t *testing.T) string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join("..", "test_helpers", "testdata", legacyMoveFixture))
	require.NoError(t, err)
	var document map[string]any
	require.NoError(t, json.Unmarshal(content, &document))
	delete(document, "manualZones")
	delete(document, "manualConnections")
	stripped, err := json.Marshal(document)
	require.NoError(t, err)
	path := filepath.Join(t.TempDir(), "noEdits.gen.json")
	require.NoError(t, os.WriteFile(path, stripped, 0o600))

	return path
}

// saveAndReload saves the session and loads the written file into a new one.
func saveAndReload(t *testing.T, state *drivers.State) *drivers.State {
	t.Helper()
	state.SaveStateToFile(filepath.Join(t.TempDir(), "state.gen.json"))
	status, isError := state.GetStatus()
	require.False(t, isError, status)
	reloaded := newUIState()
	reloaded.LoadStateFromFile(state.GetCurrentPath())
	status, isError = reloaded.GetStatus()
	require.False(t, isError, status)

	return reloaded
}

func legacySnapshotZones(t *testing.T, fixture string) []template_model.Zone {
	t.Helper()
	loaded, err := composition.InitializeGuiHandler().
		LoadState(filepath.Join("..", "test_helpers", "testdata", fixture), false)
	require.NoError(t, err)
	zones, _ := loaded.State.ManualEditSettings.ToSnapshot()

	return zones
}

func freshLayoutOf(t *testing.T, state *drivers.State) template_model.Variant {
	t.Helper()
	settings := state.GetStateDto()
	settings.ManualEditSettings = editor_state_model.ManualEditSettings{}
	generated, err := composition.InitializeGuiHandler().GenerateTemplate(settings)
	require.NoError(t, err)

	return generated.Template.Variants[0]
}

func resourceValues(zones []template_model.Zone) map[string]int {
	values := make(map[string]int, len(zones))
	for _, zone := range zones {
		values[zone.Name] = zone.ResourcesValue
	}

	return values
}

func reprofiledCastleCounts(t *testing.T, state *drivers.State) []int {
	t.Helper()
	var counts []int
	for _, edit := range state.GetStateData().ManualZones {
		if _, recorded := edit.GetCastles(); !recorded || edit.IsAdded() {
			continue
		}
		for _, zone := range state.GetLastTemplate().Variants[0].Zones {
			if zone.Name == edit.Name {
				counts = append(counts, test_helpers.NewZoneEditorService().CountZoneCastles(zone))
			}
		}
	}

	return counts
}

func rewriteGeneratorVersion(t *testing.T, path string, version int) {
	t.Helper()
	content, err := os.ReadFile(path)
	require.NoError(t, err)
	rewritten := strings.Replace(string(content), `"generatorVersion": 1`, `"generatorVersion": `+strconv.Itoa(version), 1)
	require.NotEqual(t, string(content), rewritten, "the saved file holds no generator version to rewrite")
	require.NoError(t, os.WriteFile(path, []byte(rewritten), 0o600))
}
