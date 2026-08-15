//go:build integration_test

package integration_common

import "image"

// Window-pixel coordinates used by the snapshot handlers. They are measured at
// WindowWidth x WindowHeight with PxPerDp = 1, not derived from
// app/gui/constants/ui.go, because the toolbar, the tab strip and the preview
// panel are laid out by flex weights and rendered label widths, which no padding
// constant expresses.
//
// The tab coordinates were measured by replaying a frame through
// utils.ButtonPositionLogger. The mask rectangles were measured by capturing two
// unmasked runs and diffing them: only the preview canvas interior, the status
// message (it carries a timestamp) and the output directory path (it is
// per-machine, see AGENTS.md 2.7) differ. Re-measure the same way if the preview
// panel layout changes.
const (
	// tabStripCenterY is the vertical center of the tab row.
	tabStripCenterY = 60

	generalTabX        = 672
	layoutAndZonesTabX = 789
	bonusesAndBansTabX = 933

	// previewPanelContentLeft/Right bound the preview panel's content column,
	// inside its frame.
	previewPanelContentLeft  = 1157
	previewPanelContentRight = 1583

	// previewCanvasBorder* bound the bordered square the map is drawn into. The
	// border itself is deterministic and is deliberately left unmasked.
	previewCanvasBorderLeft   = 1162
	previewCanvasBorderTop    = 202
	previewCanvasBorderRight  = 1578
	previewCanvasBorderBottom = 628

	// statusMessage* bound the fixed-height block that reports the last
	// generation, ending with a wall-clock timestamp.
	statusMessageTop    = 726
	statusMessageBottom = 775

	// outputDirectoryPath* bound the read-only textbox showing the detected
	// Olden Era templates directory.
	outputDirectoryPathTop    = 809
	outputDirectoryPathBottom = 838
)

// previewCanvasMask covers the interior of the preview canvas, whose contents
// are regenerated with fresh randomness whenever the topology is Random.
func previewCanvasMask() image.Rectangle {
	return image.Rect(
		previewCanvasBorderLeft+1, previewCanvasBorderTop+1,
		previewCanvasBorderRight-1, previewCanvasBorderBottom-1)
}

// statusMessageMask covers the status line, which embeds a timestamp.
func statusMessageMask() image.Rectangle {
	return image.Rect(
		previewPanelContentLeft, statusMessageTop,
		previewPanelContentRight, statusMessageBottom)
}

// outputDirectoryMask covers the resolved output directory, which differs
// between machines, operating systems and game installations.
func outputDirectoryMask() image.Rectangle {
	return image.Rect(
		previewPanelContentLeft, outputDirectoryPathTop,
		previewPanelContentRight, outputDirectoryPathBottom)
}
