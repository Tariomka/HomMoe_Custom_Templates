//go:build integration_test

package integration_test

import (
	"image"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers/integration_common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	// Calibrate BEFORE enabling snapshots: probing performs many throwaway
	// clicks that must not produce snapshots.
	points := integration_common.CalibrateTabPoints(t, runner)
	require.NotEmpty(t, points)

	runner.EnableSnapshots(t)
	runner.MaskRect(previewMask)

	for tabIndex, tabPoint := range points {
		runner.ClickAt(tabPoint)
		assert.Equal(t, tabIndex, runner.SelectedTabIndex())
	}
}
