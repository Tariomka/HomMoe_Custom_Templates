//go:build integration_test && gui

package gui_test

import (
	"image"
	"testing"

	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/constants"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/dialogs"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/interfaces"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// frameBonusPicker lays out one bonus-composer frame the way DialogHost does on
// every vsync and reports whether the dialog asked to close.
func frameBonusPicker(t *testing.T, dialog *dialogs.BonusPickerDialog, theme *material.Theme) bool {
	t.Helper()
	gtx, frameRouter := newDialogContext(image.Pt(460, 420))
	_, done := dialog.Body(gtx, theme)
	frameRouter.Frame(gtx.Ops)
	return done
}

// defaultReceiver mirrors the receiver the dialog uses before the user touches
// the receiver dropdown.
func defaultReceiver() string { return constants.GetBonusReceiverOptions()[0] }

// pickerCapture records the sub-picker the bonus dialog pushes onto the dialog
// stack, standing in for the real DialogHost.
type pickerCapture struct {
	picker *dialogs.MultiSelectPicker
}

func (this *pickerCapture) opener(t *testing.T) interfaces.DialogOpener {
	t.Helper()
	return func(dialog interfaces.IDialog) {
		picker, ok := dialog.(*dialogs.MultiSelectPicker)
		require.True(t, ok, "the bonus dialog must push a MultiSelectPicker")
		this.picker = picker
	}
}

// pickSpellsThroughSubPicker drives the real "Pick spells..." flow: the opener
// captures the spell picker the bonus dialog pushes, the wanted spells are
// clicked on it, and its Add button feeds them back into the bonus dialog.
func pickSpellsThroughSubPicker(
	t *testing.T,
	dialog *dialogs.BonusPickerDialog,
	capture *pickerCapture,
	spellSids []string,
	theme *material.Theme) {
	t.Helper()
	dialog.ClickPickSpells()
	require.False(t, frameBonusPicker(t, dialog, theme))
	require.NotNil(t, capture.picker, "the bonus dialog must push a spell picker")

	for _, sid := range spellSids {
		selectPickerEntry(t, capture.picker, sid, theme)
	}
	capture.picker.ClickAdd()
	require.True(t, framePicker(t, capture.picker, theme))
}

func TestWhenBonusPickerRenders_FillsTheAvailableSpace(t *testing.T) {
	t.Parallel()
	// Arrange
	dialog := dialogs.NewBonusPickerDialog(nil, nil, nil)
	gtx, frameRouter := newDialogContext(image.Pt(460, 420))

	// Act
	dimensions, closed := dialog.Body(gtx, themes.NewTheme())
	frameRouter.Frame(gtx.Ops)

	// Assert
	require.False(t, closed)
	assert.Equal(t, 460, dimensions.Size.X)
}

func TestWhenTheDefaultBonusIsAdded_EmitsAFreeTownPortalEntry(t *testing.T) {
	t.Parallel()
	// Arrange
	theme := themes.NewTheme()
	var applied []config.BonusEntry
	dialog := dialogs.NewBonusPickerDialog(nil, nil, func(entries []config.BonusEntry) { applied = entries })

	// Act
	dialog.ClickAdd()
	done := frameBonusPicker(t, dialog, theme)

	// Assert
	require.True(t, done)
	assert.Equal(t, []config.BonusEntry{{
		PresetType:     config.BonusTownPortalFree,
		ReceiverFilter: defaultReceiver(),
	}}, applied)
}

func TestWhenTheComposedBonusAlreadyExists_ReportsADuplicate(t *testing.T) {
	t.Parallel()
	// Arrange
	theme := themes.NewTheme()
	existing := []config.BonusEntry{{
		PresetType:     config.BonusTownPortalFree,
		ReceiverFilter: defaultReceiver(),
	}}
	applied := false
	dialog := dialogs.NewBonusPickerDialog(existing, nil, func([]config.BonusEntry) { applied = true })

	// Act
	dialog.ClickAdd()
	done := frameBonusPicker(t, dialog, theme)

	// Assert
	require.False(t, done)
	require.False(t, applied)
	assert.Equal(t, "That bonus already exists.", dialog.ErrorText())
}

func TestWhenTheMultiplierIsNotNumeric_ReportsAValidationError(t *testing.T) {
	t.Parallel()
	// Arrange
	theme := themes.NewTheme()
	dialog := dialogs.NewBonusPickerDialog(nil, nil, nil)
	require.True(t, dialog.SelectType("Unit Multiplier"))
	dialog.SetMultiplier("not a number")

	// Act
	dialog.ClickAdd()
	require.False(t, frameBonusPicker(t, dialog, theme))

	// Assert
	assert.Equal(t, "Enter a numeric multiplier.", dialog.ErrorText())
}

func TestWhenTheMultiplierIsNumeric_EmitsAUnitMultiplierEntry(t *testing.T) {
	t.Parallel()
	// Arrange
	theme := themes.NewTheme()
	var applied []config.BonusEntry
	dialog := dialogs.NewBonusPickerDialog(nil, nil, func(entries []config.BonusEntry) { applied = entries })
	require.True(t, dialog.SelectType("Unit Multiplier"))
	dialog.SetMultiplier(" 3.5 ")

	// Act
	dialog.ClickAdd()
	require.True(t, frameBonusPicker(t, dialog, theme))

	// Assert
	assert.Equal(t, []config.BonusEntry{{
		PresetType:     config.BonusUnitMultiplier,
		ReceiverFilter: defaultReceiver(),
		Param:          "3.5",
	}}, applied)
}

func TestWhenTheMovementValueIsNotNumeric_ReportsAValidationError(t *testing.T) {
	t.Parallel()
	// Arrange
	theme := themes.NewTheme()
	dialog := dialogs.NewBonusPickerDialog(nil, nil, nil)
	require.True(t, dialog.SelectType("Movement Bonus"))
	dialog.SetMovement("")

	// Act
	dialog.ClickAdd()
	require.False(t, frameBonusPicker(t, dialog, theme))

	// Assert
	assert.Equal(t, "Enter a numeric movement value.", dialog.ErrorText())
}

func TestWhenNoStartingItemWasPicked_ReportsAValidationError(t *testing.T) {
	t.Parallel()
	// Arrange
	theme := themes.NewTheme()
	dialog := dialogs.NewBonusPickerDialog(nil, nil, nil)
	require.True(t, dialog.SelectType("Starting Item"))

	// Act
	dialog.ClickAdd()
	require.False(t, frameBonusPicker(t, dialog, theme))

	// Assert
	assert.Equal(t, "Pick or enter an item.", dialog.ErrorText())
}

func TestWhenAStartingItemWasPicked_EmitsAStartingItemEntry(t *testing.T) {
	t.Parallel()
	// Arrange
	theme := themes.NewTheme()
	item := constants.GetBannableItemsWithExclusions(nil)[0].Sid
	var applied []config.BonusEntry
	dialog := dialogs.NewBonusPickerDialog(nil, nil, func(entries []config.BonusEntry) { applied = entries })
	require.True(t, dialog.SelectType("Starting Item"))
	dialog.SetItem(item)

	// Act
	dialog.ClickAdd()
	require.True(t, frameBonusPicker(t, dialog, theme))

	// Assert
	assert.Equal(t, []config.BonusEntry{{
		PresetType:     config.BonusStartingItem,
		ReceiverFilter: defaultReceiver(),
		Param:          item,
	}}, applied)
}

func TestWhenTheResourceAmountIsNotNumeric_ReportsAValidationError(t *testing.T) {
	t.Parallel()
	// Arrange
	theme := themes.NewTheme()
	dialog := dialogs.NewBonusPickerDialog(nil, nil, nil)
	require.True(t, dialog.SelectType("Starting Gold"))
	dialog.SetResourceAmount("plenty")

	// Act
	dialog.ClickAdd()
	require.False(t, frameBonusPicker(t, dialog, theme))

	// Assert
	assert.Equal(t, "Enter a numeric amount.", dialog.ErrorText())
}

func TestWhenTheResourceAmountIsNumeric_EmitsAResourceEntry(t *testing.T) {
	t.Parallel()
	// Arrange
	theme := themes.NewTheme()
	var applied []config.BonusEntry
	dialog := dialogs.NewBonusPickerDialog(nil, nil, func(entries []config.BonusEntry) { applied = entries })
	require.True(t, dialog.SelectType("Starting Gold"))
	dialog.SetResourceAmount("10000")

	// Act
	dialog.ClickAdd()
	require.True(t, frameBonusPicker(t, dialog, theme))

	// Assert
	assert.Equal(t, []config.BonusEntry{{
		PresetType:     config.BonusStartingGold,
		ReceiverFilter: defaultReceiver(),
		Param:          "10000",
	}}, applied)
}

func TestWhenNoSpellWasPicked_ReportsAValidationError(t *testing.T) {
	t.Parallel()
	// Arrange
	theme := themes.NewTheme()
	dialog := dialogs.NewBonusPickerDialog(nil, nil, nil)
	require.True(t, dialog.SelectType("Spell"))

	// Act
	dialog.ClickAdd()
	require.False(t, frameBonusPicker(t, dialog, theme))

	// Assert
	assert.Equal(t, "Pick at least one spell.", dialog.ErrorText())
}

func TestWhenSpellsArePickedAndMadeFree_EmitsOneFreeEntryPerSpell(t *testing.T) {
	t.Parallel()
	// Arrange
	theme := themes.NewTheme()
	spells := constants.GetKnownSpellsWithExclusions(nil)
	require.GreaterOrEqual(t, len(spells), 2)
	wanted := []string{spells[0].Sid, spells[1].Sid}

	capture := &pickerCapture{}
	var applied []config.BonusEntry
	dialog := dialogs.NewBonusPickerDialog(
		nil, capture.opener(t), func(entries []config.BonusEntry) { applied = entries })
	require.True(t, dialog.SelectType("Spell"))
	pickSpellsThroughSubPicker(t, dialog, capture, wanted, theme)
	dialog.SetMakeFree(true)

	// Act
	dialog.ClickAdd()
	require.True(t, frameBonusPicker(t, dialog, theme))

	// Assert
	assert.Equal(t, []config.BonusEntry{
		{PresetType: config.BonusSpell, ReceiverFilter: defaultReceiver(), Param: wanted[0], Param2: "1"},
		{PresetType: config.BonusSpell, ReceiverFilter: defaultReceiver(), Param: wanted[1], Param2: "1"},
	}, applied)
}

func TestWhenAPickedSpellIsRemoved_ItIsAbsentFromTheSelection(t *testing.T) {
	t.Parallel()
	// Arrange
	theme := themes.NewTheme()
	spells := constants.GetKnownSpellsWithExclusions(nil)
	require.GreaterOrEqual(t, len(spells), 2)
	wanted := []string{spells[0].Sid, spells[1].Sid}

	capture := &pickerCapture{}
	dialog := dialogs.NewBonusPickerDialog(nil, capture.opener(t), nil)
	require.True(t, dialog.SelectType("Spell"))
	pickSpellsThroughSubPicker(t, dialog, capture, wanted, theme)
	require.Equal(t, wanted, dialog.SelectedSpells())

	// Act
	require.True(t, dialog.ClickRemoveSpell(0))
	require.False(t, frameBonusPicker(t, dialog, theme))

	// Assert
	assert.Equal(t, []string{wanted[1]}, dialog.SelectedSpells())
}

func TestWhenSpellsArePicked_TheCaptionCountsThem(t *testing.T) {
	t.Parallel()
	// Arrange
	theme := themes.NewTheme()
	spells := constants.GetKnownSpellsWithExclusions(nil)
	require.GreaterOrEqual(t, len(spells), 2)

	capture := &pickerCapture{}
	dialog := dialogs.NewBonusPickerDialog(nil, capture.opener(t), nil)
	require.True(t, dialog.SelectType("Spell"))
	require.Equal(t, "Spells", dialog.SpellCountLabel())

	// Act
	pickSpellsThroughSubPicker(t, dialog, capture, []string{spells[0].Sid, spells[1].Sid}, theme)

	// Assert
	assert.Equal(t, "2 spells picked", dialog.SpellCountLabel())
}

func TestWhenAnExistingSpellBonusIsPresent_TheSpellPickerExcludesIt(t *testing.T) {
	t.Parallel()
	// Arrange
	theme := themes.NewTheme()
	spell := constants.GetKnownSpellsWithExclusions(nil)[0].Sid
	existing := []config.BonusEntry{{
		PresetType:     config.BonusSpell,
		ReceiverFilter: defaultReceiver(),
		Param:          spell,
	}}
	capture := &pickerCapture{}
	dialog := dialogs.NewBonusPickerDialog(existing, capture.opener(t), nil)
	require.True(t, dialog.SelectType("Spell"))

	// Act
	dialog.ClickPickSpells()
	require.False(t, frameBonusPicker(t, dialog, theme))

	// Assert
	require.NotNil(t, capture.picker)
	assert.NotContains(t, capture.picker.EntryIDs(), spell)
}

func TestWhenTheBonusPickerIsCancelled_NothingIsApplied(t *testing.T) {
	t.Parallel()
	// Arrange
	theme := themes.NewTheme()
	applied := false
	dialog := dialogs.NewBonusPickerDialog(nil, nil, func([]config.BonusEntry) { applied = true })

	// Act
	dialog.ClickCancel()
	require.True(t, frameBonusPicker(t, dialog, theme))

	// Assert
	assert.False(t, applied)
}

func TestWhenSeveralStartingItemsArePicked_EachBecomesItsOwnEntry(t *testing.T) {
	t.Parallel()
	// Arrange
	theme := themes.NewTheme()
	items := constants.GetBannableItemsWithExclusions(nil)
	require.GreaterOrEqual(t, len(items), 2)
	wanted := []string{items[0].Sid, items[1].Sid}

	capture := &pickerCapture{}
	var applied []config.BonusEntry
	dialog := dialogs.NewBonusPickerDialog(
		nil, capture.opener(t), func(entries []config.BonusEntry) { applied = entries })
	require.True(t, dialog.SelectType("Starting Item"))
	dialog.ClickPickItem()
	require.False(t, frameBonusPicker(t, dialog, theme))
	require.NotNil(t, capture.picker)
	for _, sid := range wanted {
		selectPickerEntry(t, capture.picker, sid, theme)
	}
	capture.picker.ClickAdd()
	require.True(t, framePicker(t, capture.picker, theme))

	// Act
	done := frameBonusPicker(t, dialog, theme)

	// Assert
	require.True(t, done)
	assert.Equal(t, []config.BonusEntry{
		{PresetType: config.BonusStartingItem, ReceiverFilter: defaultReceiver(), Param: wanted[0]},
		{PresetType: config.BonusStartingItem, ReceiverFilter: defaultReceiver(), Param: wanted[1]},
	}, applied)
}
