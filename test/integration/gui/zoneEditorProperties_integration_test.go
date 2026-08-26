//go:build integration_test && gui

package gui_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers/integration_common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The option labels the zone editor's side panel offers. They are spelled out
// rather than read back from the production tables, so a table that silently
// changes shows up here as a missing button instead of as a test that still
// agrees with whatever the code now says.
const (
	directConnectionTypeLabel = "Direct"
	spawnAGuardZoneLabel      = "Spawn-A"
	lowGuardPresetLabel       = "Low (52000)"
	fastWeeklyIncrementLabel  = "Fast (20%)"
	goldZoneQualityLabel      = "Gold"
	threeCastlesLabel         = "3"
	goldOnlyContentPoolEntry  = "classic_template_pool_random_t5_item"
	probeMatchGroup           = "rnd_guard_batch_h"
)

// The zone name is a read-only material label, not an editor: zonePropertyRows
// draws it with material.Body1 and the dialog offers no way to rename a zone, so
// the backlog's TestWhenAZoneNameIsTyped_... has nothing to drive and is not
// written here.
//
// Two further layout facts the tests below depend on. First, Gio inserts typed
// text at the caret and the caret sits at the start of a freshly focused field,
// so typing "1" into a field showing "0.2" leaves "10.2" - the expectations are
// written for insertion, not replacement. Second, the side panel's row
// coordinates were measured on a zone carrying a one-line note, which is what a
// player spawn and a neutral zone have; the shared Hub's note wraps differently
// and its rows do not line up, so the textbox tests drive a spawn and the
// dropdown tests drive a placed neutral zone.

// editedZone finds a zone as the open editor currently holds it, which is what
// the property widgets write back into every frame.
func editedZone(
	t *testing.T,
	zoneEditor *integration_common.ZoneEditorHandler,
	name string) entities.Zone {
	t.Helper()
	for _, zone := range zoneEditor.Dialog().EditedZones() {
		if zone.Name == name {
			return zone
		}
	}
	t.Fatalf("the editor is not holding a zone called %q", name)

	return entities.Zone{}
}

// manualConnectionSave finds the committed manual record of a connection, which
// is what Apply writes and the only place a connection's edited properties can
// be read back from.
func manualConnectionSave(
	t *testing.T,
	runner *integration_common.AppRunner,
	name string) editor_state_model.ManualConnectionSave {
	t.Helper()
	for _, save := range runner.CurrentState().ManualConnections {
		if save.Connection.Name == name {
			return save
		}
	}
	t.Fatalf("the editor state committed no manual connection called %q", name)

	return editor_state_model.ManualConnectionSave{}
}

// selectPlacedNeutralZone places a zone on empty canvas and selects it, which is
// the only way to reach the quality and castle dropdowns: they are drawn for
// neutral zones only, and the Geometric Hub layout ships none.
func selectPlacedNeutralZone(zoneEditor *integration_common.ZoneEditorHandler) {
	zoneEditor.ClickAddZone().ClickCanvasAt(emptyCanvasSpot).ClickZone(placedZoneName)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAZoneSizeIsTyped_TheZoneRecordsIt(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditor(t, geometricHubLayout, true)
	zoneEditor.ClickZone(spawnAZoneName)

	// Act
	zoneEditor.TypeZoneSize(".5")

	// Assert
	assert.InDelta(t, 0.51, editedZone(t, zoneEditor, spawnAZoneName).Size, 1e-9)
}

// A size above the range the field advertises is pulled down to it rather than
// refused, so the zone can never carry a size the generator will not honour.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAZoneSizeIsAboveTheMaximum_ItIsClampedToTwo(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditor(t, geometricHubLayout, true)
	zoneEditor.ClickZone(spawnAZoneName)

	// Act
	zoneEditor.TypeZoneSize("9")

	// Assert
	assert.InDelta(t, 2.0, editedZone(t, zoneEditor, spawnAZoneName).Size, 1e-9)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAZoneSizeIsBelowTheMinimum_ItIsClampedToATenth(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditor(t, geometricHubLayout, true)
	zoneEditor.ClickZone(spawnAZoneName)

	// Act
	zoneEditor.TypeZoneSize("0.05")

	// Assert
	assert.InDelta(t, 0.1, editedZone(t, zoneEditor, spawnAZoneName).Size, 1e-9)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAZoneSizeHasMoreThanTwoDecimals_ItIsRoundedToTwo(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditor(t, geometricHubLayout, true)
	zoneEditor.ClickZone(spawnAZoneName)

	// Act
	zoneEditor.TypeZoneSize("0.456")

	// Assert
	assert.InDelta(t, 0.46, editedZone(t, zoneEditor, spawnAZoneName).Size, 1e-9)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAZoneGuardMultiplierIsTyped_TheZoneRecordsIt(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditor(t, geometricHubLayout, true)
	zoneEditor.ClickZone(spawnAZoneName)

	// Act
	zoneEditor.TypeZoneGuard("2.")

	// Assert
	assert.InDelta(t, 2.1, editedZone(t, zoneEditor, spawnAZoneName).GuardMultiplier, 1e-9)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAZoneWeeklyIncrementIsTyped_TheZoneRecordsIt(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditor(t, geometricHubLayout, true)
	zoneEditor.ClickZone(spawnAZoneName)

	// Act
	zoneEditor.TypeZoneWeekly("1")

	// Assert
	assert.InDelta(t, 10.2, editedZone(t, zoneEditor, spawnAZoneName).GuardWeeklyIncrement, 1e-9)
}

// Picking a quality runs the zone back through ApplyZoneEditorQuality, so the
// content pool is rebuilt for the new tier rather than merely relabelled: a
// silver zone draws tier-three content, a gold one draws tier five.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenANeutralZoneQualityIsSelected_ItsContentIsReprofiled(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditor(t, geometricHubLayout, true)
	selectPlacedNeutralZone(zoneEditor)
	require.NotContains(t, editedZone(t, zoneEditor, placedZoneName).GuardedContentPool,
		goldOnlyContentPoolEntry, "the placed zone must start below gold for the reprofile to show")

	// Act
	zoneEditor.SelectZoneQuality(goldZoneQualityLabel)

	// Assert
	assert.Contains(t,
		editedZone(t, zoneEditor, placedZoneName).GuardedContentPool, goldOnlyContentPoolEntry)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenANeutralZoneCastleCountIsSelected_ItHoldsThatManyCastles(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditor(t, geometricHubLayout, true)
	selectPlacedNeutralZone(zoneEditor)

	// Act
	zoneEditor.SelectZoneCastles(threeCastlesLabel)

	// Assert
	assert.Len(t, editedZone(t, zoneEditor, placedZoneName).MainObjects, 3)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAConnectionGuardValueIsTyped_TheAppliedConnectionRecordsIt(t *testing.T) {
	// Arrange
	runner, zoneEditor := openZoneEditor(t, geometricHubLayout, true)
	zoneEditor.ClickConnection(hubToSpawnAName)

	// Act
	zoneEditor.TypeConnectionGuardValue("1").ClickApply()

	// Assert
	assert.Equal(t, 135000, manualConnectionSave(t, runner, hubToSpawnAName).Connection.GuardValue)
}

// A guard value that is not a number leaves the connection alone: the writeback
// only assigns what parses, so a half-typed field cannot zero a guard.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAConnectionGuardValueIsNotNumeric_ThePreviousValueIsKept(t *testing.T) {
	// Arrange
	runner, zoneEditor := openZoneEditor(t, geometricHubLayout, true)
	zoneEditor.ClickConnection(hubToSpawnAName)

	// Act
	zoneEditor.TypeConnectionGuardValue("abc").ClickApply()

	// Assert
	assert.Equal(t, 35000, manualConnectionSave(t, runner, hubToSpawnAName).Connection.GuardValue)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAConnectionIncrementIsTyped_TheAppliedConnectionRecordsIt(t *testing.T) {
	// Arrange
	runner, zoneEditor := openZoneEditor(t, geometricHubLayout, true)
	zoneEditor.ClickConnection(hubToSpawnAName)

	// Act
	zoneEditor.TypeConnectionIncrement("7").ClickApply()

	// Assert
	assert.InDelta(t, 70.15,
		manualConnectionSave(t, runner, hubToSpawnAName).Connection.GuardWeeklyIncrement, 1e-9)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAConnectionTypeIsSelected_TheAppliedConnectionRecordsIt(t *testing.T) {
	// Arrange
	runner, zoneEditor := openZoneEditor(t, geometricHubLayout, true)
	zoneEditor.ClickConnection(hubToSpawnAName)

	// Act
	zoneEditor.SelectConnectionType(directConnectionTypeLabel).ClickApply()

	// Assert
	assert.Equal(t, directConnectionTypeLabel,
		manualConnectionSave(t, runner, hubToSpawnAName).Connection.ConnectionType)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAConnectionGuardZoneIsSelected_TheAppliedConnectionRecordsIt(t *testing.T) {
	// Arrange
	runner, zoneEditor := openZoneEditor(t, geometricHubLayout, true)
	zoneEditor.ClickConnection(hubToSpawnAName)

	// Act
	zoneEditor.SelectConnectionGuardZone(spawnAGuardZoneLabel).ClickApply()

	// Assert
	assert.Equal(t, spawnAZoneName,
		manualConnectionSave(t, runner, hubToSpawnAName).Connection.GuardZone)
}

// The preset does not carry a value of its own: picking one rewrites the guard
// value field below it, and that field is what the connection is written from.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAConnectionGuardPresetIsSelected_TheGuardValueFollowsIt(t *testing.T) {
	// Arrange
	runner, zoneEditor := openZoneEditor(t, geometricHubLayout, true)
	zoneEditor.ClickConnection(hubToSpawnAName)

	// Act
	zoneEditor.SelectConnectionGuardPreset(lowGuardPresetLabel).ClickApply()

	// Assert
	assert.Equal(t, 52000, manualConnectionSave(t, runner, hubToSpawnAName).Connection.GuardValue)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAConnectionWeeklyIncrementIsSelected_TheIncrementFollowsIt(t *testing.T) {
	// Arrange
	runner, zoneEditor := openZoneEditor(t, geometricHubLayout, true)
	zoneEditor.ClickConnection(hubToSpawnAName)

	// Act
	zoneEditor.SelectConnectionWeekly(fastWeeklyIncrementLabel).ClickApply()

	// Assert
	assert.InDelta(t, 0.2,
		manualConnectionSave(t, runner, hubToSpawnAName).Connection.GuardWeeklyIncrement, 1e-9)
}

// The match group row only exists while the advanced options are shown, so this
// also covers the checkbox that reveals it.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAdvancedOptionsAreShownAndAMatchGroupIsTyped_TheAppliedConnectionRecordsIt(t *testing.T) {
	// Arrange
	runner, zoneEditor := openZoneEditor(t, geometricHubLayout, true)
	zoneEditor.ClickConnection(hubToSpawnAName).ToggleAdvancedOptions()

	// Act
	zoneEditor.TypeConnectionMatchGroup(probeMatchGroup).ClickApply()

	// Assert
	assert.Equal(t, probeMatchGroup,
		manualConnectionSave(t, runner, hubToSpawnAName).Connection.GuardMatchGroup)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenGuardEscapeIsToggled_TheAppliedConnectionRecordsIt(t *testing.T) {
	// Arrange
	runner, zoneEditor := openZoneEditor(t, geometricHubLayout, true)
	zoneEditor.ClickConnection(hubToSpawnAName).ToggleAdvancedOptions()

	// Act
	zoneEditor.ToggleGuardEscape().ClickApply()

	// Assert
	assert.True(t, manualConnectionSave(t, runner, hubToSpawnAName).Connection.GuardEscape)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenSimTurnSquadIsToggled_TheAppliedConnectionRecordsIt(t *testing.T) {
	// Arrange
	runner, zoneEditor := openZoneEditor(t, geometricHubLayout, true)
	zoneEditor.ClickConnection(hubToSpawnAName).ToggleAdvancedOptions()

	// Act
	zoneEditor.ToggleSimTurnSquad().ClickApply()

	// Assert
	assert.True(t, manualConnectionSave(t, runner, hubToSpawnAName).Connection.SimTurnSquad)
}
