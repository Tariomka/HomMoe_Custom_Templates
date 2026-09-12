//go:build integration_test && gui

package gui_test

import (
	"path/filepath"
	"strconv"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers/integration_common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// idleFrameCount is how many input-free frames the idle tests let the window
// draw. The panel writes its widgets back on every one of them, so anything
// that re-applies an edit rather than merely redrawing it shows up within a
// handful.
const idleFrameCount = 4

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

// The guard preset tables a placed neutral zone's incident connection is
// offered, spelled out the same way. A placed zone starts Silver and the tests
// below raise it to Gold, so both tables are named here, together with the
// Custom option that stands for a number in neither of them.
const (
	customGuardPresetLabel        = "Custom"
	silverDefaultGuardPresetLabel = "Default (20000)"
	silverMediumGuardPresetLabel  = "Medium (24000)"
	silverHighGuardPresetLabel    = "High (27000)"
	goldMediumGuardPresetLabel    = "Medium (48000)"

	silverDefaultGuardValue = 20000
	silverMediumGuardValue  = 24000
	silverHighGuardValue    = 27000
	goldDefaultGuardValue   = 25000
	goldMediumGuardValue    = 48000

	// customGuardValue is in no preset table of any quality, so it stays Custom
	// however the zone it hangs off is re-tiered.
	customGuardValue = 12345
	// editedGuardValue is typed after a quality change, to prove the edit lands
	// in the connection the editor kept rather than in a discarded copy.
	editedGuardValue = 51000
)

// The two victory conditions the combined-flow tests switch to, as
// app/gui/constants spells them, and the notice the discard they cause writes.
// Guardian Arena is the victory that turns the effective arena on; Tournament is
// the one that turns the effective tournament on.
const (
	guardianArenaVictoryLabel = "Guardian Arena"
	tournamentVictoryLabel    = "Tournament"

	discardedLayoutNotice = "The manual zone layout was discarded because " +
		"the game mode change regenerates the map."
	generatedStatusFragment = "generated with latest changes"
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
	name string) template_model.Zone {
	t.Helper()
	for _, zone := range zoneEditor.Dialog().EditedZones() {
		if zone.Name == name {
			return zone
		}
	}
	t.Fatalf("the editor is not holding a zone called %q", name)

	return template_model.Zone{}
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

// openEditorWithNeutralEdge places a neutral zone, draws a connection from it to
// a player spawn and selects that connection.
//
// It is the only reachable setup for the guard propagation tests. A connection's
// guard table comes from its stronger endpoint, a player spawn resolves as
// Unknown, and the layout's only other neutral is the Hub - whose tier is fixed
// by its name whatever the quality dropdown says. So the placed zone, which
// starts Silver, is the one endpoint a test can move, and this edge is the one
// that follows it. Snapshots stay off: these tests assert what the panel says
// and what the connection carries, not what the frame looks like.
func openEditorWithNeutralEdge(t *testing.T) (
	*integration_common.AppRunner, *integration_common.ZoneEditorHandler) {
	t.Helper()
	runner, zoneEditor := openZoneEditor(t, geometricHubLayout, false)
	zoneEditor.ClickAddZone().
		ClickCanvasAt(emptyCanvasSpot).
		ClickAddConnection().
		DragFromZoneTo(placedZoneName, spawnAZoneName).
		ClickAddConnection().
		ClickConnectionBetween(placedZoneName, spawnAZoneName)
	require.Equal(t, silverDefaultGuardValue,
		pendingConnection(t, zoneEditor, placedZoneName, spawnAZoneName).GuardValue,
		"the drawn edge must start on the placed zone's Silver table for these tests to mean anything")
	require.True(t, zoneEditor.Dialog().SelectedConnectionIsUserAdded(),
		"the drawn edge must be the selection before the property panel can be driven")

	return runner, zoneEditor
}

// neutralEdgeGuardValue reads the guard the drawn edge currently carries inside
// the open editor.
func neutralEdgeGuardValue(
	t *testing.T, zoneEditor *integration_common.ZoneEditorHandler) int {
	t.Helper()
	return pendingConnection(t, zoneEditor, placedZoneName, spawnAZoneName).GuardValue
}

// manualConnectionSaveBetween finds the committed manual record of the edge
// joining two zones. A user-added connection carries no name, so it is found by
// its endpoints.
func manualConnectionSaveBetween(
	t *testing.T,
	runner *integration_common.AppRunner,
	from string,
	to string) editor_state_model.ManualConnectionSave {
	t.Helper()
	for _, save := range runner.CurrentState().ManualConnections {
		if (save.Connection.From == from && save.Connection.To == to) ||
			(save.Connection.From == to && save.Connection.To == from) {
			return save
		}
	}
	t.Fatalf("the editor state committed no manual connection between %q and %q", from, to)

	return editor_state_model.ManualConnectionSave{}
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

// Raising a neutral zone's tier carries its incident guards to the same named
// tier of the new table rather than leaving the old numbers behind: a Silver
// Medium guard becomes a Gold Medium one.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenANeutralZoneQualityRises_ItsIncidentGuardKeepsItsNamedTier(t *testing.T) {
	// Arrange
	_, zoneEditor := openEditorWithNeutralEdge(t)
	zoneEditor.SelectConnectionGuardPreset(silverMediumGuardPresetLabel)
	require.Equal(t, silverMediumGuardValue, neutralEdgeGuardValue(t, zoneEditor))

	// Act
	zoneEditor.ClickZone(placedZoneName).SelectZoneQuality(goldZoneQualityLabel)

	// Assert
	assert.Equal(t, goldMediumGuardValue, neutralEdgeGuardValue(t, zoneEditor))
}

// Default is a named tier like any other, not a fixed number, so it moves with
// the table too.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenANeutralZoneQualityRises_ItsDefaultGuardBecomesTheNewDefault(t *testing.T) {
	// Arrange
	_, zoneEditor := openEditorWithNeutralEdge(t)

	// Act
	zoneEditor.ClickZone(placedZoneName).SelectZoneQuality(goldZoneQualityLabel)

	// Assert
	assert.Equal(t, goldDefaultGuardValue, neutralEdgeGuardValue(t, zoneEditor))
}

// A number that is in no preset table is the user's own, so a re-tier leaves it
// exactly where they put it.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenANeutralZoneQualityRises_ACustomIncidentGuardIsLeftAlone(t *testing.T) {
	// Arrange
	_, zoneEditor := openEditorWithNeutralEdge(t)
	zoneEditor.SetConnectionGuardValue(strconv.Itoa(customGuardValue))
	require.Equal(t, customGuardValue, neutralEdgeGuardValue(t, zoneEditor))

	// Act
	zoneEditor.ClickZone(placedZoneName).SelectZoneQuality(goldZoneQualityLabel)

	// Assert
	assert.Equal(t, customGuardValue, neutralEdgeGuardValue(t, zoneEditor))
}

// The panel has to agree with the number: after the re-tier the dropdown names
// the tier the new value belongs to, read off the control the user sees.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenARemappedConnectionIsReselected_TheDropdownNamesItsNewTier(t *testing.T) {
	// Arrange
	_, zoneEditor := openEditorWithNeutralEdge(t)
	zoneEditor.SelectConnectionGuardPreset(silverMediumGuardPresetLabel).
		ClickZone(placedZoneName).
		SelectZoneQuality(goldZoneQualityLabel)

	// Act
	zoneEditor.ClickConnectionBetween(placedZoneName, spawnAZoneName)

	// Assert
	assert.Equal(t, goldMediumGuardPresetLabel, zoneEditor.Dialog().SelectedGuardPresetLabel())
}

// Custom is an option of its own, drawn after the six numeric presets, so a
// guard that matches none of them has something to display.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAConnectionIsSelected_CustomIsOfferedAfterEveryPreset(t *testing.T) {
	// Arrange
	_, zoneEditor := openEditorWithNeutralEdge(t)

	// Act
	labels := zoneEditor.Dialog().GuardPresetLabels()

	// Assert
	assert.Equal(t, []string{
		silverDefaultGuardPresetLabel,
		"Weakest (18000)",
		"Low (21000)",
		silverMediumGuardPresetLabel,
		silverHighGuardPresetLabel,
		"Very High (30000)",
		customGuardPresetLabel,
	}, labels)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAGuardValueMatchesNoPreset_TheDropdownShowsCustom(t *testing.T) {
	// Arrange
	_, zoneEditor := openEditorWithNeutralEdge(t)

	// Act
	zoneEditor.SetConnectionGuardValue(strconv.Itoa(customGuardValue))

	// Assert
	assert.Equal(t, customGuardPresetLabel, zoneEditor.Dialog().SelectedGuardPresetLabel())
}

// Custom carries no number of its own, so picking it is not an edit: the guard
// the user typed stays exactly as it was.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenCustomIsPicked_TheGuardValueIsUnchanged(t *testing.T) {
	// Arrange
	_, zoneEditor := openEditorWithNeutralEdge(t)
	zoneEditor.SetConnectionGuardValue(strconv.Itoa(customGuardValue))

	// Act
	zoneEditor.SelectConnectionGuardPreset(customGuardPresetLabel)

	// Assert
	assert.Equal(t, customGuardValue, neutralEdgeGuardValue(t, zoneEditor))
}

// Preset identity is the number and nothing else: a typed value that happens to
// equal a table entry is that tier, with no remembered selection involved.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenATypedGuardValueEqualsAPreset_TheDropdownNamesThatPreset(t *testing.T) {
	// Arrange
	_, zoneEditor := openEditorWithNeutralEdge(t)
	zoneEditor.SetConnectionGuardValue(strconv.Itoa(customGuardValue))

	// Act
	zoneEditor.SetConnectionGuardValue(strconv.Itoa(silverHighGuardValue))

	// Assert
	assert.Equal(t, silverHighGuardPresetLabel, zoneEditor.Dialog().SelectedGuardPresetLabel())
}

// Nothing remembers that the user was shown Custom, so leaving the connection
// and coming back has to work the display out from the number again.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAConnectionWithACustomGuardIsReselected_TheDropdownStillShowsCustom(t *testing.T) {
	// Arrange
	_, zoneEditor := openEditorWithNeutralEdge(t)
	zoneEditor.SetConnectionGuardValue(strconv.Itoa(customGuardValue)).
		ClickConnection(hubToSpawnAName)

	// Act
	zoneEditor.ClickConnectionBetween(placedZoneName, spawnAZoneName)

	// Assert
	assert.Equal(t, customGuardPresetLabel, zoneEditor.Dialog().SelectedGuardPresetLabel())
}

// The same holds across a reopen: the applied edit persists as a plain number,
// and the freshly built panel infers Custom from it.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenTheEditorIsReopened_ACustomGuardStillShowsCustom(t *testing.T) {
	// Arrange
	_, zoneEditor := openEditorWithNeutralEdge(t)
	reopened := zoneEditor.SetConnectionGuardValue(strconv.Itoa(customGuardValue)).
		ClickApply().
		ClickLayoutAndZonesTab().
		OpenZoneEditor()

	// Act
	reopened.ClickConnectionBetween(placedZoneName, spawnAZoneName)

	// Assert
	assert.Equal(t, customGuardPresetLabel, reopened.Dialog().SelectedGuardPresetLabel())
}

// The quality edit hands the editor a whole new set of connections, and the old
// ones are thrown away. Editing the connection afterwards has to reach the set
// the editor kept - if the selection still pointed into the discarded one, this
// guard would apply as the remapped 48000 and the user's 51000 would vanish.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAConnectionIsEditedAfterAQualityChange_TheAppliedGuardIsTheLaterEdit(t *testing.T) {
	// Arrange
	runner, zoneEditor := openEditorWithNeutralEdge(t)
	zoneEditor.SelectConnectionGuardPreset(silverMediumGuardPresetLabel).
		ClickZone(placedZoneName).
		SelectZoneQuality(goldZoneQualityLabel).
		ClickConnectionBetween(placedZoneName, spawnAZoneName)
	require.Equal(t, goldMediumGuardValue, neutralEdgeGuardValue(t, zoneEditor))

	// Act
	zoneEditor.SetConnectionGuardValue(strconv.Itoa(editedGuardValue)).ClickApply()

	// Assert
	assert.Equal(t, editedGuardValue,
		manualConnectionSaveBetween(t, runner, placedZoneName, spawnAZoneName).Connection.GuardValue)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAQualityChangeIsApplied_TheCommittedGuardCarriesTheNewValue(t *testing.T) {
	// Arrange
	runner, zoneEditor := openEditorWithNeutralEdge(t)
	zoneEditor.SelectConnectionGuardPreset(silverMediumGuardPresetLabel).
		ClickZone(placedZoneName).
		SelectZoneQuality(goldZoneQualityLabel)

	// Act
	zoneEditor.ClickApply()

	// Assert
	assert.Equal(t, goldMediumGuardValue,
		manualConnectionSaveBetween(t, runner, placedZoneName, spawnAZoneName).Connection.GuardValue)
}

// A quality edit is pending dialog work until Apply, so cancelling has to leave
// the document with no manual edit at all - not with the remapped guard.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAQualityChangeIsCancelled_NothingIsCommitted(t *testing.T) {
	// Arrange
	runner, zoneEditor := openEditorWithNeutralEdge(t)
	zoneEditor.SelectConnectionGuardPreset(silverMediumGuardPresetLabel).
		ClickZone(placedZoneName).
		SelectZoneQuality(goldZoneQualityLabel)

	// Act
	zoneEditor.ClickCancel()

	// Assert
	assert.Empty(t, runner.CurrentState().ManualConnections)
}

// idleFrames draws further frames with no input at all, which is what the window
// does between two user actions.
func idleFrames(runner *integration_common.AppRunner) {
	for range idleFrameCount {
		runner.NextFrame()
	}
}

// manualZoneNames reports the zones the editor state currently records as
// manually laid out, which is the layout a reopened dialog is built from.
func manualZoneNames(runner *integration_common.AppRunner) []string {
	names := make([]string, 0)
	for _, zone := range runner.CurrentState().ManualZones {
		names = append(names, zone.Name)
	}

	return names
}

// Custom is not a value, so picking it cannot overwrite a guard that happens to
// sit on a preset either.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenCustomIsPickedOnAPresetGuard_TheGuardValueIsUnchanged(t *testing.T) {
	// Arrange
	_, zoneEditor := openEditorWithNeutralEdge(t)
	zoneEditor.SetConnectionGuardValue(strconv.Itoa(silverHighGuardValue))

	// Act
	zoneEditor.SelectConnectionGuardPreset(customGuardPresetLabel)

	// Assert
	assert.Equal(t, silverHighGuardValue, neutralEdgeGuardValue(t, zoneEditor))
}

// Nothing remembers that Custom was picked: the display is worked out from the
// number every frame, so asking for Custom over a number that names a tier snaps
// straight back to that tier instead of sticking on Custom.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenCustomIsPickedOnAPresetGuard_TheDropdownNamesThatPresetAgain(t *testing.T) {
	// Arrange
	_, zoneEditor := openEditorWithNeutralEdge(t)
	zoneEditor.SetConnectionGuardValue(strconv.Itoa(silverHighGuardValue))

	// Act
	zoneEditor.SelectConnectionGuardPreset(customGuardPresetLabel)

	// Assert
	assert.Equal(t, silverHighGuardPresetLabel, zoneEditor.Dialog().SelectedGuardPresetLabel())
}

// A half-typed field keeps the guard, and the dropdown has to keep saying so: it
// is driven by the number the connection still carries, not by the text.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAGuardValueIsNotNumeric_TheDropdownKeepsTheLastPresetLabel(t *testing.T) {
	// Arrange
	_, zoneEditor := openEditorWithNeutralEdge(t)
	zoneEditor.SetConnectionGuardValue(strconv.Itoa(silverHighGuardValue))

	// Act
	zoneEditor.SetConnectionGuardValue("abc")

	// Assert
	assert.Equal(t, silverHighGuardPresetLabel, zoneEditor.Dialog().SelectedGuardPresetLabel())
}

// The re-tier is a one-off edit, but the panel writes its widgets back on every
// frame afterwards. Those frames must leave the remapped guard exactly as the
// re-tier left it rather than re-running the preset that produced it.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenIdleFramesFollowAQualityChange_TheRemappedGuardIsKept(t *testing.T) {
	// Arrange
	runner, zoneEditor := openEditorWithNeutralEdge(t)
	zoneEditor.SelectConnectionGuardPreset(silverMediumGuardPresetLabel).
		ClickZone(placedZoneName).
		SelectZoneQuality(goldZoneQualityLabel).
		ClickConnectionBetween(placedZoneName, spawnAZoneName)
	require.Equal(t, goldMediumGuardValue, neutralEdgeGuardValue(t, zoneEditor))

	// Act
	idleFrames(runner)

	// Assert
	assert.Equal(t, goldMediumGuardValue, neutralEdgeGuardValue(t, zoneEditor))
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenIdleFramesFollowAQualityChange_TheDropdownLabelDoesNotDrift(t *testing.T) {
	// Arrange
	runner, zoneEditor := openEditorWithNeutralEdge(t)
	zoneEditor.SelectConnectionGuardPreset(silverMediumGuardPresetLabel).
		ClickZone(placedZoneName).
		SelectZoneQuality(goldZoneQualityLabel).
		ClickConnectionBetween(placedZoneName, spawnAZoneName)
	require.Equal(t, goldMediumGuardPresetLabel, zoneEditor.Dialog().SelectedGuardPresetLabel())

	// Act
	idleFrames(runner)

	// Assert
	assert.Equal(t, goldMediumGuardPresetLabel, zoneEditor.Dialog().SelectedGuardPresetLabel())
}

// The saved file holds a plain number and no display state at all, so the round
// trip through disk has to arrive at the same reading as the reopen does.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenTheStateIsSavedAndLoaded_ACustomGuardStillShowsCustom(t *testing.T) {
	// Arrange
	_, zoneEditor := openEditorWithNeutralEdge(t)
	explorer := zoneEditor.SetConnectionGuardValue(strconv.Itoa(customGuardValue)).
		ClickApply().
		ClickSaveTo().
		ClickSave()
	require.FileExists(t, filepath.Join(explorer.FixtureDirectory(), defaultSaveFile))
	reloaded := explorer.Editor().
		ClickLoad().
		ClickRow(defaultSaveFile).
		ClickOpen().
		Editor().
		ClickLayoutAndZonesTab().
		OpenZoneEditor()

	// Act
	reloaded.ClickConnectionBetween(placedZoneName, spawnAZoneName)

	// Assert
	assert.Equal(t, customGuardPresetLabel, reloaded.Dialog().SelectedGuardPresetLabel())
}

// Cancelling throws away the dialog's work, not the work an earlier Apply
// already committed: the guard the document is carrying stays the one that was
// applied, not the one the cancelled re-tier would have produced.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAQualityChangeIsCancelled_AnAlreadyCommittedGuardIsKept(t *testing.T) {
	// Arrange
	runner, zoneEditor := openEditorWithNeutralEdge(t)
	reopened := zoneEditor.SelectConnectionGuardPreset(silverMediumGuardPresetLabel).
		ClickApply().
		ClickLayoutAndZonesTab().
		OpenZoneEditor()
	require.Equal(t, silverMediumGuardValue,
		manualConnectionSaveBetween(t, runner, placedZoneName, spawnAZoneName).Connection.GuardValue)
	reopened.ClickZone(placedZoneName).SelectZoneQuality(goldZoneQualityLabel)

	// Act
	reopened.ClickCancel()

	// Assert
	assert.Equal(t, silverMediumGuardValue,
		manualConnectionSaveBetween(t, runner, placedZoneName, spawnAZoneName).Connection.GuardValue)
}

// The layout the cancelled dialog was opened over is committed work too, so the
// zone the user placed and applied is still there afterwards.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAQualityChangeIsCancelled_TheCommittedManualZoneSurvives(t *testing.T) {
	// Arrange
	runner, zoneEditor := openEditorWithNeutralEdge(t)
	reopened := zoneEditor.SelectConnectionGuardPreset(silverMediumGuardPresetLabel).
		ClickApply().
		ClickLayoutAndZonesTab().
		OpenZoneEditor()
	reopened.ClickZone(placedZoneName).SelectZoneQuality(goldZoneQualityLabel)

	// Act
	reopened.ClickCancel()

	// Assert
	assert.Contains(t, manualZoneNames(runner), placedZoneName)
}

// applyRemappedGuardThenSelectVictory drives the whole flow the mode contract is
// about, in one go and through real input: a pending quality edit remaps the
// drawn edge's guard onto the Gold table, Apply commits that as the document's
// manual layout, and a victory condition picked on the General tab turns an
// effective mode on. It returns the runner and the zone editor reopened over
// whatever the editor regenerated in place of the discarded layout.
func applyRemappedGuardThenSelectVictory(t *testing.T, victoryLabel string) (
	*integration_common.AppRunner, *integration_common.ZoneEditorHandler) {
	t.Helper()
	runner, zoneEditor := openEditorWithNeutralEdge(t)
	editorTabs := zoneEditor.SelectConnectionGuardPreset(silverMediumGuardPresetLabel).
		ClickZone(placedZoneName).
		SelectZoneQuality(goldZoneQualityLabel).
		ClickApply()
	require.Equal(t, goldMediumGuardValue,
		manualConnectionSaveBetween(t, runner, placedZoneName, spawnAZoneName).Connection.GuardValue,
		"precondition: Apply commits the remapped guard, so the mode change has real work to discard")
	require.Contains(t, manualZoneNames(runner), placedZoneName,
		"precondition: Apply commits the placed zone too")

	general := editorTabs.ClickGeneralTab()
	// The tab click is applied during the layout it is polled on, so the General
	// panel has not been drawn yet and its victory selector has no input area.
	runner.NextFrame()
	general.SelectVictoryCondition(victoryLabel)
	idleFrames(runner)

	return runner, general.ClickLayoutAndZonesTab().OpenZoneEditor()
}

// editedMainObjectTypes reports the main object types a zone the editor is
// holding carries.
func editedMainObjectTypes(zone template_model.Zone) []string {
	types := make([]string, 0)
	for _, mainObject := range zone.MainObjects {
		types = append(types, mainObject.Type)
	}

	return types
}

// editedEdgeEndpoints reports every edge the editor is holding, as an
// order-independent "From|To" pair, which is how a user-drawn edge is told apart
// from the generated ones.
func editedEdgeEndpoints(zoneEditor *integration_common.ZoneEditorHandler) []string {
	endpoints := make([]string, 0)
	for _, connection := range zoneEditor.Dialog().EditedConnectionRecords() {
		from, to := connection.From, connection.To
		if from > to {
			from, to = to, from
		}
		endpoints = append(endpoints, from+"|"+to)
	}

	return endpoints
}

// Turning the arena on is a layout-defining change, so the map is rebuilt for
// that mode. On a hub layout the arena is a main object inside the Hub, which is
// the visible proof that the graph now on screen is the one the new mode asked
// for rather than the one the discarded edits described.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenArenaModeFollowsAnAppliedQualityEdit_TheRegeneratedHubCarriesTheArena(t *testing.T) {
	// Arrange
	_, reopened := applyRemappedGuardThenSelectVictory(t, guardianArenaVictoryLabel)

	// Act
	hub := editedZone(t, reopened, hubZoneName)

	// Assert
	assert.Contains(t, editedMainObjectTypes(hub), registry.GetMainObjectTypeValues().GladiatorArena)
}

// The zone the user placed and applied belonged to the discarded layout, so the
// regenerated map must not carry it back.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenArenaModeFollowsAnAppliedQualityEdit_ThePlacedZoneIsNotReapplied(t *testing.T) {
	// Arrange
	_, reopened := applyRemappedGuardThenSelectVictory(t, guardianArenaVictoryLabel)

	// Act
	names := editedZoneNames(reopened)

	// Assert
	assert.NotContains(t, names, placedZoneName)
}

// Neither is the edge that carried the remapped guard: it was drawn by hand
// between the placed zone and a spawn, and nothing in a generated hub layout
// joins those two.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenArenaModeFollowsAnAppliedQualityEdit_TheRemappedEdgeIsNotReapplied(t *testing.T) {
	// Arrange
	_, reopened := applyRemappedGuardThenSelectVictory(t, guardianArenaVictoryLabel)

	// Act
	endpoints := editedEdgeEndpoints(reopened)

	// Assert
	assert.NotContains(t, endpoints, placedZoneName+"|"+spawnAZoneName)
}

// The discard is announced, and the regeneration it triggers writes its own
// status straight afterwards. The notice has to survive that, or the user would
// never learn their layout is gone.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenArenaModeFollowsAnAppliedQualityEdit_TheDiscardIsStillReported(t *testing.T) {
	// Arrange
	runner, _ := applyRemappedGuardThenSelectVictory(t, guardianArenaVictoryLabel)

	// Act
	message, _ := runner.Status()

	// Assert
	require.Contains(t, message, generatedStatusFragment,
		"precondition: the mode change regenerated, so the notice had a status to survive")
	assert.Contains(t, message, discardedLayoutNotice)
}

// The other effective mode goes through the same contract. Tournament replaces
// the chosen topology outright, so the Hub the Geometric Hub layout is built
// around is not in the map that mode generates.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenTournamentModeFollowsAnAppliedQualityEdit_TheHubLayoutIsReplaced(t *testing.T) {
	// Arrange
	_, reopened := applyRemappedGuardThenSelectVictory(t, tournamentVictoryLabel)
	require.NotEmpty(t, editedZoneNames(reopened),
		"precondition: the reopened editor is holding the regenerated map, not an empty canvas")

	// Act
	names := editedZoneNames(reopened)

	// Assert
	assert.NotContains(t, names, hubZoneName)
}

// And the zone the user placed and applied is not carried into it either.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenTournamentModeFollowsAnAppliedQualityEdit_ThePlacedZoneIsNotReapplied(t *testing.T) {
	// Arrange
	_, reopened := applyRemappedGuardThenSelectVictory(t, tournamentVictoryLabel)
	require.NotEmpty(t, editedZoneNames(reopened),
		"precondition: the reopened editor is holding the regenerated map, not an empty canvas")

	// Act
	names := editedZoneNames(reopened)

	// Assert
	assert.NotContains(t, names, placedZoneName)
}
