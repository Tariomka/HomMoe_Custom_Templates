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
//
// The remaining coordinates were measured the same way, but widgets that do not
// emit semantic.Button (checkboxes and dropdown triggers) have no button bounds
// to read, so their positions were taken from the bounds.Min of the neighbouring
// label node and then confirmed by driving them and asserting the state change
// they cause. Each constant below records that confirmation. To re-measure, see
// plans/gui-handler-framework.md phase 0.
const (
	// tabStripCenterY is the vertical center of the tab row.
	tabStripCenterY = 60

	generalTabX        = 672
	layoutAndZonesTabX = 789
	bonusesAndBansTabX = 933

	// toolbarCenterY is the vertical center of the toolbar button row. The Save
	// and Exit buttons are deliberately not exposed: Save writes a real file to
	// the detected templates directory, and Exit calls os.Exit, which would kill
	// the test process.
	toolbarCenterY = 23

	newButtonX    = 32
	loadButtonX   = 84
	saveToButtonX = 200

	// mapSizeSelectorTrigger* opens the General tab's map size dropdown.
	// Confirmed by clicking it and observing the option rows appear.
	mapSizeSelectorTriggerX = 360
	mapSizeSelectorTriggerY = 199

	// mapSizeOptionFirstCenterY and mapSizeOptionPitch address the rows of the
	// open map size dropdown: row i sits at mapSizeOptionFirstCenterY +
	// i*mapSizeOptionPitch. These rows do emit semantic.Button, so both values
	// were read directly from their bounds.
	mapSizeOptionX            = 300
	mapSizeOptionFirstCenterY = 225
	mapSizeOptionPitch        = 25

	// experimentalMapSizesCheckbox* toggles "Allow non official larger map
	// sizes (>240)", which grows the map size dropdown from 11 rows to 28 and so
	// shifts everything the open dropdown pushes down. Confirmed by clicking it
	// and asserting EditorStateDto.ExperimentalMapSizes flipped.
	experimentalMapSizesCheckboxX = 100
	experimentalMapSizesCheckboxY = 234

	// gameMode*X select the General tab's game mode segment buttons at
	// gameModeCenterY. SingleHero removes the three hero count rows below them,
	// shifting the rest of the right column up. Confirmed by clicking each and
	// asserting EditorStateDto.GameMode.
	gameModeCenterY     = 126
	gameModeClassicX    = 761
	gameModeSingleHeroX = 827

	// zoneEditorButton* opens the Layout & Zones tab's "Manual zone editor..."
	// dialog. Confirmed by clicking it and asserting DialogsOpen.
	zoneEditorButtonX = 846
	zoneEditorButtonY = 129

	// advancedZoneControlCheckbox* toggles "Advanced zone control", which adds
	// enough rows to give the Layout & Zones panel a ~386px scroll range (it
	// overflows by only ~18px otherwise). Confirmed by clicking it and asserting
	// EditorStateDto.AdvancedMode.
	advancedZoneControlCheckboxX = 660
	advancedZoneControlCheckboxY = 355

	// panelScroll* is a point over the left settings column that carries no
	// interactive widget, so a wheel event there reaches the panel's list rather
	// than a slider. Confirmed by scrolling and observing the rows move.
	panelScrollX = 300
	panelScrollY = 650

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
