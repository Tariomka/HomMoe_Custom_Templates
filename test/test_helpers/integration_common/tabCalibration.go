//go:build integration_test

package integration_common

import (
	"image"
	"testing"
)

const (
	// The tab strip sits just below the toolbar. Calibration approaches it from
	// below (through the harmless panel area) so it never clicks the toolbar,
	// which holds the Exit button (os.Exit). These bounds bracket the tab row
	// for the fixed 1600x900 / 1px-per-dp layout.
	tabRowSearchTop    = 66
	tabRowSearchBottom = 150
)

// ProbeClick clicks probePoint and dismisses any modal dialog the click may
// have opened (e.g. a file picker), so calibration can keep probing the live UI.
func ProbeClick(runner *AppRunner, probePoint image.Point) {
	runner.ClickAt(probePoint)
	if runner.DialogsOpen() {
		runner.CloseTopDialog()
		runner.NextFrame()
	}
}

// CalibrateTabPoints discovers one click point inside every tab by probing the
// laid-out window. The tab strip is horizontally centered (the editor's vertical
// Flex uses Alignment: Middle) and its position depends on font metrics, so
// callers locate it dynamically instead of hard-coding fragile coordinates.
// It approaches the strip from below so it never clicks the toolbar above (which
// holds the Exit button), and records the points where a click changes the
// selected tab.
func CalibrateTabPoints(tb testing.TB, runner *AppRunner) []image.Point {
	tb.Helper()

	tabCount := runner.TabCount()
	points := make([]image.Point, tabCount)
	found := make([]bool, tabCount)
	remaining := tabCount

	// Warm-up frames register the input areas and let the default preview
	// auto-generate before probing.
	for range 3 {
		runner.NextFrame()
	}

	// Step 1: scan upward (panels -> tabs) for the first row whose click changes
	// the selected tab. That first change is always to a non-zero tab (the tab
	// strip starts on tab 0), giving a reliable anchor for step 2.
	stripY, anchor, anchorIdx := -1, image.Point{}, -1
	for y := tabRowSearchBottom; y >= tabRowSearchTop && anchorIdx < 0; y -= 2 {
		for x := 12; x <= WindowWidth-12; x += 6 {
			before := runner.SelectedTabIndex()
			probePoint := image.Pt(x, y)
			ProbeClick(runner, probePoint)
			if runner.SelectedTabIndex() == before {
				continue
			}
			stripY, anchor, anchorIdx = y, probePoint, runner.SelectedTabIndex()
			points[anchorIdx] = probePoint
			found[anchorIdx] = true
			remaining--
			break
		}
	}
	if stripY < 0 {
		tb.Fatalf("calibration failed: never located the tab strip in y=[%d,%d]", tabRowSearchTop, tabRowSearchBottom)
	}

	// Step 2: sweep the strip row. Each sweep first selects the non-zero anchor
	// tab so the initially selected tab also registers a change when hit.
	for attempt := 0; attempt < tabCount+1 && remaining > 0; attempt++ {
		ProbeClick(runner, anchor)
		for x := 12; x <= WindowWidth-12 && remaining > 0; x += 2 {
			before := runner.SelectedTabIndex()
			probePoint := image.Pt(x, stripY)
			ProbeClick(runner, probePoint)
			after := runner.SelectedTabIndex()
			if after != before && !found[after] {
				points[after] = probePoint
				found[after] = true
				remaining--
			}
		}
	}

	for tabIndex, ok := range found {
		if !ok {
			tb.Fatalf(
				"calibration failed: located tab strip at y=%d but could not find a click point for tab %d of %d",
				stripY,
				tabIndex,
				tabCount,
			)
		}
	}
	return points
}
