//go:build integration_test && gui

package gui_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers/integration_common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// bansAndOverrides flattens the three lists the Bonuses & Bans tab keeps into
// one slice, so a round trip is compared by value and by order in one go.
func bansAndOverrides(state editor_state_model.EditorState) []string {
	flattened := make([]string, 0)
	for _, ban := range state.BannedItems {
		flattened = append(flattened, "item:"+ban.Sid)
	}
	for _, ban := range state.BannedMagics {
		flattened = append(flattened, "spell:"+ban.Sid)
	}
	for _, override := range state.ValueOverrides {
		flattened = append(flattened, fmt.Sprintf("override:%s=%d", override.Sid, override.GuardValue))
	}

	return flattened
}

// savedBonusesSession drives the Bonuses & Bans tab of a fresh window through
// populate, saves the result and returns the path of the saved file and the
// lists the editor held when saving.
func savedBonusesSession(
	t *testing.T, populate func(tab *integration_common.BonusesAndBansTabHandler)) (string, []string) {
	t.Helper()
	runner := integration_common.NewAppRunner(t)
	populate(integration_common.NewHandler(runner).ClickBonusesAndBansTab())
	saved := bansAndOverrides(runner.CurrentState())
	runner.SaveStateToFile(filepath.Join(t.TempDir(), "bonuses.gen.json"))
	status, isError := runner.Status()
	require.False(t, isError, status)

	return runner.CurrentPath(), saved
}

// loadedBansAndOverrides loads a file into a fresh window and reads the lists
// the editor then holds.
func loadedBansAndOverrides(t *testing.T, path string) []string {
	t.Helper()
	runner := integration_common.NewAppRunner(t)
	integration_common.NewHandler(runner).LoadFile(path)

	return bansAndOverrides(runner.CurrentState())
}

// pickableSids picks the sids the tests ban and override: two items, one of them
// listed after the other so an order that follows the picking is told apart from
// one that follows the pickers' own lists, a spell and a guard value override.
func pickableSids(t *testing.T) (laterItem string, earlierItem string, spell string, override string) {
	t.Helper()
	items := constants.GetBannableItemsWithExclusions(nil)
	spells := constants.GetKnownSpellsWithExclusions(nil)
	overrides := constants.GetValueOverrideSidsWithExclusions(nil)
	require.GreaterOrEqual(t, len(items), 2)
	require.NotEmpty(t, spells)
	require.NotEmpty(t, overrides)

	return items[1].Sid, items[0].Sid, spells[0].Sid, overrides[0]
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenBansAndOverridesArePickedSavedAndLoaded_TheyComeBackInThePickedOrder(t *testing.T) {
	// Arrange
	laterItem, earlierItem, spell, override := pickableSids(t)
	path, saved := savedBonusesSession(t, func(tab *integration_common.BonusesAndBansTabHandler) {
		tab.BanItem(laterItem).BanItem(earlierItem).BanSpell(spell).OverrideGuardValue(override)
	})
	require.Len(t, saved, 4)

	// Act
	loaded := loadedBansAndOverrides(t, path)

	// Assert
	assert.Equal(t,
		[]string{"item:" + laterItem, "item:" + earlierItem, "spell:" + spell, saved[3]},
		loaded)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAnOverrideIsSavedAndLoaded_ItsGuardValueIsKept(t *testing.T) {
	// Arrange
	_, _, _, override := pickableSids(t)
	path, saved := savedBonusesSession(t, func(tab *integration_common.BonusesAndBansTabHandler) {
		tab.OverrideGuardValue(override)
	})

	// Act
	loaded := loadedBansAndOverrides(t, path)

	// Assert
	assert.Equal(t, saved, loaded)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenABannedSpellIsRemovedSavedAndLoaded_ItStaysRemoved(t *testing.T) {
	// Arrange
	_, _, spell, _ := pickableSids(t)
	path, _ := savedBonusesSession(t, func(tab *integration_common.BonusesAndBansTabHandler) {
		tab.BanSpell(spell).RemoveEntry(0)
	})

	// Act
	loaded := loadedBansAndOverrides(t, path)

	// Assert
	assert.Empty(t, loaded)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAnOverrideIsRemovedSavedAndLoaded_ItStaysRemoved(t *testing.T) {
	// Arrange
	_, _, _, override := pickableSids(t)
	path, _ := savedBonusesSession(t, func(tab *integration_common.BonusesAndBansTabHandler) {
		tab.OverrideGuardValue(override).RemoveEntry(0)
	})

	// Act
	loaded := loadedBansAndOverrides(t, path)

	// Assert
	assert.Empty(t, loaded)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenABannedItemIsRemovedSavedAndLoaded_ItStaysRemoved(t *testing.T) {
	// Arrange
	laterItem, earlierItem, _, _ := pickableSids(t)
	path, _ := savedBonusesSession(t, func(tab *integration_common.BonusesAndBansTabHandler) {
		tab.BanItem(laterItem).BanItem(earlierItem).RemoveEntry(0)
	})

	// Act
	loaded := loadedBansAndOverrides(t, path)

	// Assert
	assert.Equal(t, []string{"item:" + earlierItem}, loaded)
}
