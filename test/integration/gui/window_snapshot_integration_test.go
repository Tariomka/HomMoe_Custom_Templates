//go:build integration_test && gui

package gui_test

import (
	"testing"
	"time"

	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers/integration_common"
)

//nolint:paralleltest // Snapshots need exclusive access to the single headless GPU window.
func TestWindowSnapshots_TabClicksMatchGoldens(t *testing.T) {
	runner := integration_common.NewAppRunner(t)

	if !integration_common.IsHeadless() {
		runner.SetRenderDelay(500 * time.Millisecond)
	}

	handler := integration_common.NewHandler(runner).WithSnapshots()

	handler.ClickGeneralTab().
		ClickLayoutAndZonesTab().
		ClickBonusesAndBansTab().
		ClickGeneralTab()
}

// TestWindowSnapshots_ScrollMatchesGoldens captures the Layout & Zones panel
// before and after a wheel event. Advanced zone control is enabled first because
// the panel otherwise overflows by only about 18px, which is too little to tell
// a real scroll from a stuck one.
//
//nolint:paralleltest // Snapshots need exclusive access to the single headless GPU window.
func TestWindowSnapshots_ScrollMatchesGoldens(t *testing.T) {
	runner := integration_common.NewAppRunner(t)

	if !integration_common.IsHeadless() {
		runner.SetRenderDelay(500 * time.Millisecond)
	}

	handler := integration_common.NewHandler(runner).WithSnapshots()

	handler.ClickLayoutAndZonesTab().
		ToggleAdvancedZoneControl().
		ScrollPanel(400)
}

// TestWindowSnapshots_MapSizeShiftMatchesGoldens captures the General tab with
// the map size dropdown open before and after the experimental sizes are
// allowed. The dropdown renders inline, so this is where the checkbox's layout
// shift is actually visible.
//
//nolint:paralleltest // Snapshots need exclusive access to the single headless GPU window.
func TestWindowSnapshots_MapSizeShiftMatchesGoldens(t *testing.T) {
	runner := integration_common.NewAppRunner(t)

	if !integration_common.IsHeadless() {
		runner.SetRenderDelay(500 * time.Millisecond)
	}

	handler := integration_common.NewHandler(runner).WithSnapshots()

	handler.ClickGeneralTab().
		ToggleExperimentalMapSizes().
		OpenMapSizeSelector()
}
