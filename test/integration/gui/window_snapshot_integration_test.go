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
