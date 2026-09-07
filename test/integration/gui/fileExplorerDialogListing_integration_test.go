//go:build integration_test && gui

package gui_test

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers/integration_common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The fixture entries the listing tests are driven against. Names are fixed, not
// fuzzed, because every action is also compared against a golden image.
const (
	visibleFile   = "alpha.gen.json"
	visibleFolder = "archive"
	hiddenFile    = ".alpha-hidden.gen.json"
	hiddenFolder  = ".archive-hidden"
)

// scrollFixtureCount is enough rows to overflow the listing several times over.
const scrollFixtureCount = 40

// newListingHandler opens the Load dialog on a fixture directory holding one
// visible and one hidden entry of each kind.
func newListingHandler(t *testing.T) *integration_common.FileExplorerHandler {
	t.Helper()
	runner := integration_common.NewAppRunner(t)
	if !integration_common.IsHeadless() {
		runner.SetRenderDelay(500 * time.Millisecond)
	}

	return integration_common.NewHandler(runner).
		WithFixtureFiles(visibleFile, hiddenFile).
		WithFixtureFolders(visibleFolder, hiddenFolder).
		WithSnapshots().
		ClickLoad()
}

// TestWhenShowHiddenIsToggledOn_HiddenEntriesAppearInTheListing covers the
// toggle the review found untested: hidden entries are filtered out until it is
// pressed, on Windows and Linux alike.
//
//nolint:paralleltest // Snapshots need exclusive access to the single headless GPU window.
func TestWhenShowHiddenIsToggledOn_HiddenEntriesAppearInTheListing(t *testing.T) {
	// Arrange
	handler := newListingHandler(t)
	require.NotContains(t, handler.Dialog().EntryNames(), hiddenFile)

	// Act
	handler.ClickShowHidden()

	// Assert
	assert.Contains(t, handler.Dialog().EntryNames(), hiddenFile)
}

//nolint:paralleltest // Snapshots need exclusive access to the single headless GPU window.
func TestWhenShowHiddenIsToggledOff_HiddenEntriesDisappearAgain(t *testing.T) {
	// Arrange
	handler := newListingHandler(t).ClickShowHidden()
	require.Contains(t, handler.Dialog().EntryNames(), hiddenFolder)

	// Act
	handler.ClickShowHidden()

	// Assert
	assert.NotContains(t, handler.Dialog().EntryNames(), hiddenFolder)
}

// TestWhenARowIsClicked_ThatEntryBecomesTheSelection drives the selection with a
// real pointer press and release instead of a queued click. Only the open dialog
// is exercised: a row click in save mode is deliberately inert, so that the
// resolved name cannot be silently retargeted.
//
//nolint:paralleltest // Snapshots need exclusive access to the single headless GPU window.
func TestWhenARowIsClicked_ThatEntryBecomesTheSelection(t *testing.T) {
	// Arrange
	handler := newListingHandler(t)

	// Act
	handler.ClickRow(visibleFile)

	// Assert
	assert.Equal(t, filepath.Join(handler.FixtureDirectory(), visibleFile), handler.Dialog().SelectedPath())
}

//nolint:paralleltest // Snapshots need exclusive access to the single headless GPU window.
func TestWhenADirectoryRowIsClicked_TheListingDescendsIntoIt(t *testing.T) {
	// Arrange
	handler := newListingHandler(t)

	// Act
	handler.ClickRow(visibleFolder)

	// Assert
	assert.Equal(t, filepath.Join(handler.FixtureDirectory(), visibleFolder), handler.Dialog().CurrentDir())
}

// TestWhenTheListingIsScrolled_TheFirstVisibleRowAdvances proves the listing
// itself scrolls, rather than the dialog resizing around a fixed viewport.
//
//nolint:paralleltest // Snapshots need exclusive access to the single headless GPU window.
func TestWhenTheListingIsScrolled_TheFirstVisibleRowAdvances(t *testing.T) {
	// Arrange
	runner := integration_common.NewAppRunner(t)
	if !integration_common.IsHeadless() {
		runner.SetRenderDelay(500 * time.Millisecond)
	}

	names := make([]string, 0, scrollFixtureCount)
	for index := range scrollFixtureCount {
		names = append(names, fmt.Sprintf("entry-%02d.gen.json", index))
	}

	handler := integration_common.NewHandler(runner).
		WithFixtureFiles(names...).
		WithSnapshots().
		ClickLoad()

	firstBefore, _ := handler.Dialog().ScrollPosition()
	require.Zero(t, firstBefore, "the listing must start at the top")

	// Act
	handler.Scroll(400)

	// Assert
	firstAfter, _ := handler.Dialog().ScrollPosition()
	assert.Greater(t, firstAfter, firstBefore)
}
