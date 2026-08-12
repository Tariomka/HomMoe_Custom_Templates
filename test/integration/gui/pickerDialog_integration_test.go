//go:build integration_test && gui

package gui_test

import (
	"image"
	"testing"

	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/constants"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/dialogs"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
	"github.com/Tariomka/hommoe_custom_templates/internal/composition"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// framePicker lays out one picker frame the way DialogHost does on every vsync
// - clicks queued with widget.Clickable.Click are consumed here - and reports
// whether the picker asked to close.
func framePicker(t *testing.T, picker *dialogs.MultiSelectPicker, theme *material.Theme) bool {
	t.Helper()
	gtx, frameRouter := newDialogContext(image.Pt(560, 560))
	_, done := picker.Body(gtx, theme)
	frameRouter.Frame(gtx.Ops)
	return done
}

// selectPickerEntry narrows the list to the given id so its row is laid out,
// then clicks it and lets the next frame consume the click.
func selectPickerEntry(t *testing.T, picker *dialogs.MultiSelectPicker, id string, theme *material.Theme) {
	t.Helper()
	picker.SetSearch(id)
	require.True(t, picker.ClickEntry(id), "the picker must list %q", id)
	require.False(t, framePicker(t, picker, theme))
}

func TestWhenItemPickerRenders_ListsEveryBannableItem(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := make([]string, 0)
	for _, item := range constants.GetBannableItemsWithExclusions(nil) {
		expected = append(expected, item.Sid)
	}
	picker := dialogs.NewItemPickerDialog("Pick Starting Item", nil, composition.InitializeGuiHandler(), nil)

	// Act
	require.False(t, framePicker(t, picker, themes.NewTheme()))

	// Assert
	assert.Equal(t, expected, picker.EntryIDs())
}

func TestWhenItemPickerExcludesAnItem_ItIsAbsentFromTheEntries(t *testing.T) {
	t.Parallel()
	// Arrange
	excluded := constants.GetBannableItemsWithExclusions(nil)[0].Sid

	// Act
	picker := dialogs.NewItemPickerDialog("Pick Starting Item", []string{excluded}, composition.InitializeGuiHandler(), nil)

	// Assert
	assert.NotContains(t, picker.EntryIDs(), excluded)
}

func TestWhenItemPickerSelectionIsApplied_EmitsTheSelectedIds(t *testing.T) {
	t.Parallel()
	// Arrange
	theme := themes.NewTheme()
	wanted := constants.GetBannableItemsWithExclusions(nil)[0].Sid
	var applied []string
	picker := dialogs.NewItemPickerDialog("Pick Starting Item", nil, composition.InitializeGuiHandler(), func(ids []string) { applied = ids })
	selectPickerEntry(t, picker, wanted, theme)

	// Act
	picker.ClickAdd()
	done := framePicker(t, picker, theme)

	// Assert
	require.True(t, done)
	assert.Equal(t, []string{wanted}, applied)
}

func TestWhenPickerIsCancelled_NothingIsApplied(t *testing.T) {
	t.Parallel()
	// Arrange
	theme := themes.NewTheme()
	applied := false
	picker := dialogs.NewItemPickerDialog("Pick Starting Item", nil, composition.InitializeGuiHandler(), func([]string) { applied = true })
	selectPickerEntry(t, picker, constants.GetBannableItemsWithExclusions(nil)[0].Sid, theme)

	// Act
	picker.ClickCancel()
	require.True(t, framePicker(t, picker, theme))

	// Assert
	assert.False(t, applied)
}

func TestWhenSearchFilterMatchesNothing_NoRowsAreProduced(t *testing.T) {
	t.Parallel()
	// Arrange
	theme := themes.NewTheme()
	picker := dialogs.NewItemPickerDialog("Pick Starting Item", nil, composition.InitializeGuiHandler(), nil)
	picker.SetSearch("no such item exists anywhere")

	// Act
	require.False(t, framePicker(t, picker, theme))

	// Assert
	assert.Equal(t, 0, picker.RowCount(theme))
}

func TestWhenSearchFilterMatchesAnEntry_TheNonMatchingEntriesAreDropped(t *testing.T) {
	t.Parallel()
	// Arrange
	wanted := constants.GetBannableItemsWithExclusions(nil)[0].Sid
	picker := dialogs.NewItemPickerDialog("Pick Starting Item", nil, composition.InitializeGuiHandler(), nil)
	require.Greater(t, len(picker.EntryIDs()), 1)

	// Act
	picker.SetSearch(wanted)

	// Assert
	assert.Contains(t, picker.MatchingEntryIDs(), wanted)
}

func TestWhenPickerIsFlat_EveryRowIsALeafRow(t *testing.T) {
	t.Parallel()
	// Arrange
	theme := themes.NewTheme()
	picker := dialogs.NewValueOverridePickerDialog(nil, composition.InitializeGuiHandler(), nil)
	picker.SetSearch(constants.GetValueOverrideSidsWithExclusions(nil)[0])
	require.NotEmpty(t, picker.MatchingEntryIDs())

	// Act
	rowCount := picker.RowCount(theme)

	// Assert
	assert.Equal(t, len(picker.MatchingEntryIDs()), rowCount)
}

func TestWhenPickerIsGrouped_MatchingGroupsAddHeaderRows(t *testing.T) {
	t.Parallel()
	// Arrange
	theme := themes.NewTheme()
	spell := constants.GetKnownSpellsWithExclusions(nil)[0]
	picker := dialogs.NewSpellPickerDialog(nil, false, composition.InitializeGuiHandler(), nil)
	picker.SetSearch(spell.Sid)
	require.NotEmpty(t, picker.MatchingEntryIDs())

	// Act
	rowCount := picker.RowCount(theme)

	// Assert
	assert.Greater(t, rowCount, len(picker.MatchingEntryIDs()))
}

func TestWhenSpellPickerSelectionIsApplied_ReportsTheMakeFreeFlag(t *testing.T) {
	t.Parallel()
	// Arrange
	theme := themes.NewTheme()
	spell := constants.GetKnownSpellsWithExclusions(nil)[0]
	var free bool
	applied := false
	picker := dialogs.NewSpellPickerDialog(nil, true, composition.InitializeGuiHandler(), func(_ []string, makeFree bool) {
		applied = true
		free = makeFree
	})
	selectPickerEntry(t, picker, spell.Sid, theme)

	// Act
	picker.ClickAdd()
	require.True(t, framePicker(t, picker, theme))

	// Assert
	require.True(t, applied)
	assert.False(t, free)
}

func TestWhenValueOverrideSelectionIsApplied_EmitsSidEqualsDefaultGuardLines(t *testing.T) {
	t.Parallel()
	// Arrange
	theme := themes.NewTheme()
	sid := constants.GetValueOverrideSidsWithExclusions(nil)[0]
	var lines []string
	picker := dialogs.NewValueOverridePickerDialog(nil, composition.InitializeGuiHandler(), func(applied []string) { lines = applied })
	selectPickerEntry(t, picker, sid, theme)

	// Act
	picker.ClickAdd()
	require.True(t, framePicker(t, picker, theme))

	// Assert
	assert.Equal(t, []string{sid + "=5000"}, lines)
}
