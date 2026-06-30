package performance_test

import (
	"image"
	"testing"
)

const (
	tabCyclesPerOp = 3

	// The tab strip sits just below the toolbar. Calibration approaches it from
	// below (through the harmless panel area) so it never clicks the toolbar,
	// which holds the Exit button (os.Exit). These bounds bracket the tab row
	// for the fixed 1600x900 / 1px-per-dp benchmark layout.
	tabRowSearchTop    = 66
	tabRowSearchBottom = 150
)

// probeClick clicks p and dismisses any modal dialog the click may have opened
// (e.g. a file picker), so calibration can keep probing the live UI.
func probeClick(r *AppRunner, p image.Point) {
	r.ClickAt(p)
	if r.DialogsOpen() {
		r.CloseTopDialog()
		r.NextFrame()
	}
}

// calibrateTabPoints discovers one click point inside every tab by probing the
// laid-out window. The tab strip is horizontally centred (the editor's vertical
// Flex uses Alignment: Middle) and its position depends on font metrics, so the
// benchmark locates it dynamically instead of hard-coding fragile coordinates.
// It approaches the strip from below so it never clicks the toolbar above (which
// holds the Exit button), and records the points where a click changes the
// selected tab.
func calibrateTabPoints(tb testing.TB, r *AppRunner) []image.Point {
	tb.Helper()

	tabCount := r.TabCount()
	points := make([]image.Point, tabCount)
	found := make([]bool, tabCount)
	remaining := tabCount

	// Warm-up frames register the input areas and let the default preview
	// auto-generate before probing.
	for range 3 {
		r.NextFrame()
	}

	// Step 1: scan upward (panels -> tabs) for the first row whose click changes
	// the selected tab. That first change is always to a non-zero tab (the tab
	// strip starts on tab 0), giving a reliable anchor for step 2.
	stripY, anchor, anchorIdx := -1, image.Point{}, -1
	for y := tabRowSearchBottom; y >= tabRowSearchTop && anchorIdx < 0; y -= 2 {
		for x := 12; x <= benchWindowWidth-12; x += 6 {
			before := r.SelectedTabIndex()
			p := image.Pt(x, y)
			probeClick(r, p)
			if r.SelectedTabIndex() == before {
				continue
			}
			stripY, anchor, anchorIdx = y, p, r.SelectedTabIndex()
			points[anchorIdx] = p
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
		probeClick(r, anchor)
		for x := 12; x <= benchWindowWidth-12 && remaining > 0; x += 2 {
			before := r.SelectedTabIndex()
			p := image.Pt(x, stripY)
			probeClick(r, p)
			after := r.SelectedTabIndex()
			if after != before && !found[after] {
				points[after] = p
				found[after] = true
				remaining--
			}
		}
	}

	for i, ok := range found {
		if !ok {
			tb.Fatalf("calibration failed: located tab strip at y=%d but could not find a click point for tab %d of %d", stripY, i, tabCount)
		}
	}
	return points
}

// BenchmarkEditorWindow_TabCycling starts the editor, clicks through every tab
// several times and gracefully shuts down. It runs identically headless or with
// a real on-screen window (the latter additionally renders the UI); select the
// mode with `go test ... -args headless`.
func BenchmarkEditorWindow_TabCycling(b *testing.B) {
	runner := NewAppRunner()
	runner.Start()
	defer runner.Stop()

	points := calibrateTabPoints(b, runner)

	// Warm-up frames let the default state auto-generate its preview so the timed
	// loop measures steady-state tab navigation, not first-frame setup.
	for range 5 {
		runner.NextFrame()
	}

	b.ReportAllocs()

	for b.Loop() {
		for range tabCyclesPerOp {
			for idx := range points {
				runner.ClickAt(points[idx])
				if got := runner.SelectedTabIndex(); got != idx {
					b.Fatalf("clicking tab %d selected tab %d instead", idx, got)
				}
				// Render an extra frame so the freshly selected panel is laid out
				// fully, exercising each tab's content every cycle.
				runner.NextFrame()
			}
		}
	}
}
