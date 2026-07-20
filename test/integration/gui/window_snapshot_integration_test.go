//go:build integration_test && gui

package gui_test

import (
	"image"
	"testing"
	"time"

	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers/integration_common"
)

var previewMask = image.Rect( //nolint:gochecknoglobals // shared test fixture rect.
	integration_common.WindowWidth-470, 0,
	integration_common.WindowWidth, integration_common.WindowHeight)

//nolint:paralleltest // Snapshots need exclusive access to the single headless GPU window.
func TestWindowSnapshots_TabClicksMatchGoldens(t *testing.T) {
	runner := integration_common.NewAppRunner()
	runner.Start()
	defer runner.Stop()

	runner.EnableSnapshots(t)
	runner.MaskRect(previewMask)
	if !integration_common.IsHeadless() {
		runner.SetRenderDelay(500 * time.Millisecond)
	}

	handler := integration_common.NewHandler(runner)

	handler.ClickGeneralTab().
		ClickLayoutAndZonesTab().
		ClickBonusesAndBansTab().
		ClickGeneralTab()
}
