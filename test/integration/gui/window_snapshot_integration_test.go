//go:build integration_test && gui

package gui_test

import (
	"image"
	"testing"
	"time"

	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers/integration_common"
)

// previewMask blanks the right-hand preview column (preview panel is capped at
// 440dp wide and anchored to the right edge). The preview renders a randomly
// generated template that rerolls every run, so it must never reach a snapshot.
var previewMask = image.Rect( //nolint:gochecknoglobals // shared test fixture rect.
	integration_common.WindowWidth-470, 0,
	integration_common.WindowWidth, integration_common.WindowHeight)

// TestWindowSnapshots_TabClicksMatchGoldens clicks through every editor tab with
// snapshot verification enabled: each click captures a headless screenshot,
// masks the preview column and compares against the committed golden
// (regenerate with `go test -tags=integration_test ./test/integration/... -run
// TestWindowSnapshots -args -update`).
//
//nolint:paralleltest // Snapshots need exclusive access to the single headless GPU window.
func TestWindowSnapshots_TabClicksMatchGoldens(t *testing.T) {
	runner := integration_common.NewAppRunner()
	runner.Start()
	defer runner.Stop()

	runner.EnableSnapshots(t)
	runner.MaskRect(previewMask)
	runner.SetRenderDelay(500 * time.Millisecond)

	handler := integration_common.NewHandler(runner)

	handler.ClickGeneralTab().
		ClickLayoutAndZonesTab().
		ClickBonusesAndBansTab().
		ClickGeneralTab()
}
