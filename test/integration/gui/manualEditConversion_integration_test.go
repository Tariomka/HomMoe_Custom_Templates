//go:build integration_test && gui

package gui_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers/integration_common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	legacyMoveFixture     = "editorState_v2_square2p_move.gen.json"
	conversionNoticeTitle = "Older file format converted"

	// postConversionFrameCount is how many input-free frames the tests let the
	// window draw to prove a dialog that was dismissed is not opened again.
	postConversionFrameCount = 5
)

// legacyMoveFixturePath is where the committed legacy v2 fixture with a moved zone
// lives, relative to this package.
func legacyMoveFixturePath() string {
	return filepath.Join("..", "..", "test_helpers", "testdata", legacyMoveFixture)
}

// readStateDocument reads an editor state file as a generic document, keeping
// numbers exactly as written.
func readStateDocument(t *testing.T, path string) map[string]any {
	t.Helper()
	content, err := os.ReadFile(path)
	require.NoError(t, err)
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.UseNumber()
	var document map[string]any
	require.NoError(t, decoder.Decode(&document))

	return document
}

// writeStateDocument writes a document into a throwaway directory and returns
// the path of the file.
func writeStateDocument(t *testing.T, document map[string]any) string {
	t.Helper()
	content, err := json.Marshal(document)
	require.NoError(t, err)
	path := filepath.Join(t.TempDir(), "variant.gen.json")
	require.NoError(t, os.WriteFile(path, content, 0o600))

	return path
}

// legacyFileWithoutManualEdits is the legacy move fixture with its manual
// section removed.
func legacyFileWithoutManualEdits(t *testing.T) string {
	t.Helper()
	document := readStateDocument(t, legacyMoveFixturePath())
	delete(document, "manualZones")
	delete(document, "manualConnections")

	return writeStateDocument(t, document)
}

// legacyFileWithValueOverrideText is the legacy move fixture carrying the given
// value override text.
func legacyFileWithValueOverrideText(t *testing.T, text string) string {
	t.Helper()
	document := readStateDocument(t, legacyMoveFixturePath())
	document["valueOverrides"] = text

	return writeStateDocument(t, document)
}

// convertedFileEditedBy loads the legacy move fixture, lets the window convert
// it, saves the result in the current format and returns a copy of that file
// after edit has rewritten its document.
func convertedFileEditedBy(t *testing.T, edit func(document map[string]any)) string {
	t.Helper()
	runner := integration_common.NewAppRunner(t)
	integration_common.NewHandler(runner).
		LoadFile(legacyMoveFixturePath()).
		ClickMessageDialogOK()
	runner.SaveStateToFile(filepath.Join(t.TempDir(), "converted.gen.json"))
	document := readStateDocument(t, runner.CurrentPath())
	edit(document)

	return writeStateDocument(t, document)
}

// loadedStatus loads a file into a fresh window, lets it regenerate and reports
// the status line.
func loadedStatus(t *testing.T, path string) string {
	t.Helper()
	runner := integration_common.NewAppRunner(t)
	integration_common.NewHandler(runner).LoadFile(path)
	status, _ := runner.Status()

	return status
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenALegacyFileWithEditsIsLoaded_TheFirstRegenerationOpensTheConversionNotice(t *testing.T) {
	// Arrange
	runner := integration_common.NewAppRunner(t)
	handler := integration_common.NewHandler(runner)

	// Act
	handler.LoadFile(legacyMoveFixturePath())

	// Assert
	assert.Equal(t, conversionNoticeTitle, runner.TopDialogTitle())
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenTheConversionNoticeIsConfirmed_TheNoticeCloses(t *testing.T) {
	// Arrange
	runner := integration_common.NewAppRunner(t)
	handler := integration_common.NewHandler(runner).LoadFile(legacyMoveFixturePath())
	require.True(t, runner.DialogsOpen(), "the notice must be open for there to be anything to confirm")

	// Act
	handler.ClickMessageDialogOK()

	// Assert
	assert.False(t, runner.DialogsOpen())
}

// The notice is announced once per load: later regenerations must not reopen it.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenTheConversionNoticeIsConfirmedAndFramesFollow_TheNoticeDoesNotReopen(t *testing.T) {
	// Arrange
	runner := integration_common.NewAppRunner(t)
	integration_common.NewHandler(runner).
		LoadFile(legacyMoveFixturePath()).
		ClickMessageDialogOK()

	// Act
	for range postConversionFrameCount {
		runner.NextFrame()
	}

	// Assert
	assert.False(t, runner.DialogsOpen())
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenALegacyFileWithEditsIsConverted_TheDocumentBecomesUnsaved(t *testing.T) {
	// Arrange
	runner := integration_common.NewAppRunner(t)
	handler := integration_common.NewHandler(runner)

	// Act
	handler.LoadFile(legacyMoveFixturePath())

	// Assert
	assert.True(t, runner.IsUnsaved())
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenALegacyFileWithoutEditsIsLoaded_NoNoticeOpens(t *testing.T) {
	// Arrange
	runner := integration_common.NewAppRunner(t)
	handler := integration_common.NewHandler(runner)

	// Act
	handler.LoadFile(legacyFileWithoutManualEdits(t))

	// Assert
	assert.False(t, runner.DialogsOpen())
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenALegacyFileWithoutEditsIsLoaded_TheDocumentStaysSaved(t *testing.T) {
	// Arrange
	runner := integration_common.NewAppRunner(t)
	handler := integration_common.NewHandler(runner)

	// Act
	handler.LoadFile(legacyFileWithoutManualEdits(t))

	// Assert
	assert.False(t, runner.IsUnsaved())
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAFileWithEditsRecordsAnotherGeneratorVersion_TheStatusSaysUneditedValuesMayDiffer(t *testing.T) {
	// Arrange
	path := convertedFileEditedBy(t, func(document map[string]any) {
		document["generatorVersion"] = json.Number("7")
	})

	// Act
	status := loadedStatus(t, path)

	// Assert
	assert.Contains(t, status, "Saved with generator version 7")
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAFileWithoutEditsRecordsAnotherGeneratorVersion_TheStatusIsSilentAboutIt(t *testing.T) {
	// Arrange
	path := convertedFileEditedBy(t, func(document map[string]any) {
		document["generatorVersion"] = json.Number("7")
		delete(document, "manualEdits")
	})

	// Act
	status := loadedStatus(t, path)

	// Assert
	assert.NotContains(t, status, "Saved with generator version")
}

// The warning is part of the load status, which the regeneration that follows
// overwrites, so it is read straight after the load.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenALegacyFileHoldsAMalformedValueOverrideLine_TheLoadStatusWarnsAboutIt(t *testing.T) {
	// Arrange
	runner := integration_common.NewAppRunner(t)
	path := legacyFileWithValueOverrideText(t, "broken")

	// Act
	runner.LoadStateFromFile(path)

	// Assert
	status, _ := runner.Status()
	assert.Contains(t, status, "value override line 1 'broken' is not sid=value and was dropped")
}

//nolint:paralleltest // Snapshots need exclusive access to the single headless GPU window.
func TestWhenALegacyFileWithEditsIsConverted_TheConversionNoticeMatchesGolden(t *testing.T) {
	// Arrange
	runner := integration_common.NewAppRunner(t)
	handler := integration_common.NewHandler(runner).WithFixtureDirectory().WithSnapshots()

	// Act
	handler.LoadFile(legacyMoveFixturePath())

	// Assert
	assert.Equal(t, conversionNoticeTitle, runner.TopDialogTitle())
}
