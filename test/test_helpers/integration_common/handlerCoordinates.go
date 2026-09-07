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

	// topologySelectorTrigger* opens the Layout & Zones tab's topology dropdown.
	// Confirmed by clicking it and then clicking a row, asserting
	// EditorStateDto.Topology.
	topologySelectorTriggerX = 366
	topologySelectorTriggerY = 124

	// topologyOption* bound the block the open topology dropdown draws its rows
	// in. The rows do emit semantic.Button, so they are addressed by label
	// inside this rectangle rather than by a row index.
	topologyOptionsLeft   = 178
	topologyOptionsTop    = 138
	topologyOptionsRight  = 556
	topologyOptionsBottom = 420

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

	// fileDialogPanel* bound the modal panel DialogHost centers on screen for
	// the file explorer's 720x560 PreferredSize. Buttons are looked up inside
	// this rectangle so the dialog's own Save and Cancel are told apart from the
	// toolbar's, which stay laid out behind the scrim.
	fileDialogPanelWidth  = 720
	fileDialogPanelHeight = 560
	fileDialogPanelLeft   = (WindowWidth - fileDialogPanelWidth) / 2
	fileDialogPanelTop    = (WindowHeight - fileDialogPanelHeight) / 2

	// fileDialogListScroll* is a point over the file explorer's listing. Every
	// row is clickable, but a row absorbs taps only - the wheel reaches the
	// enclosing material.List.
	fileDialogListScrollX = fileDialogPanelLeft + fileDialogPanelWidth/2
	fileDialogListScrollY = fileDialogPanelTop + 200

	// fileStatus* bound the toolbar's right-hand "File: ..." column. It reports
	// the path of the file being edited: the per-run temporary directory in a
	// fixture-backed test, and a per-machine path otherwise.
	fileStatusLeft   = 1000
	fileStatusBottom = 46

	// headerBarSlack widens the file explorer's path bar mask vertically,
	// because the mask is derived from the header buttons flanking the bar and
	// the textbox between them is not exactly as tall as they are.
	headerBarSlack = 6

	// zoneEditorPanel* bound the modal panel DialogHost centers on screen for
	// the zone editor's 1000x720 PreferredSize. Its buttons are looked up inside
	// this rectangle so they are told apart from the tab's own, which stay laid
	// out behind the scrim.
	zoneEditorPanelWidth  = 1000
	zoneEditorPanelHeight = 720
	zoneEditorPanelLeft   = (WindowWidth - zoneEditorPanelWidth) / 2
	zoneEditorPanelTop    = (WindowHeight - zoneEditorPanelHeight) / 2

	// zoneEditorCanvasBox* is where the space reserved for the canvas starts.
	// The canvas centres a square inside that space and reports the centring
	// offset through CanvasOrigin(), so only this origin has to be measured -
	// the square itself is derived. The calibration test presses every zone's
	// mapped point and asserts it became the selection.
	zoneEditorCanvasBoxLeft = 314
	zoneEditorCanvasBoxTop  = 181

	// zoneEditorSnapCheckbox* is the toolbar's "Snap" checkbox, which emits a
	// semantic.CheckBox rather than a labelled button. Confirmed by clicking it
	// and asserting SnapEnabled.
	zoneEditorSnapCheckboxX = 809
	zoneEditorSnapCheckboxY = 153

	// zoneEditorSidePanel* bound the properties column to the right of the
	// canvas. Its editors and dropdown triggers are centered on
	// zoneEditorSidePanelFieldX, its full-width buttons on
	// zoneEditorSidePanelButtonX.
	zoneEditorSidePanelLeft     = 996
	zoneEditorSidePanelRight    = 1266
	zoneEditorSidePanelFieldX   = 1186
	zoneEditorSidePanelButtonX  = 1131
	zoneEditorSidePanelNoteDrop = 29

	// zoneEditorZone*Y are the property rows of a selected zone, measured on a
	// player spawn. A neutral zone carries no "spawn" note row, so every one of
	// its rows sits zoneEditorSidePanelNoteDrop higher - see zoneRowY.
	zoneEditorZoneSizeY    = 259
	zoneEditorZoneGuardY   = 290
	zoneEditorZoneWeeklyY  = 317
	zoneEditorZoneQualityY = 348
	zoneEditorZoneCastlesY = 375

	// zoneEditorConnection*Y are the property rows of a selected connection with
	// its advanced options collapsed.
	zoneEditorConnectionTypeY        = 230
	zoneEditorConnectionGuardZoneY   = 257
	zoneEditorConnectionGuardPresetY = 288
	zoneEditorConnectionGuardValueY  = 315
	zoneEditorConnectionWeeklyY      = 346
	zoneEditorConnectionIncrementY   = 373
	zoneEditorConnectionAdvancedY    = 413

	// zoneEditorConnectionAdvanced*Y are the rows the advanced options reveal.
	zoneEditorConnectionMatchGroupY   = 446
	zoneEditorConnectionGuardEscapeY  = 480
	zoneEditorConnectionSimTurnSquadY = 520
)

// zoneEditorRect is the modal panel the zone editor is drawn into.
func zoneEditorRect() image.Rectangle {
	return image.Rect(
		zoneEditorPanelLeft, zoneEditorPanelTop,
		zoneEditorPanelLeft+zoneEditorPanelWidth, zoneEditorPanelTop+zoneEditorPanelHeight)
}

// zoneEditorSidePanelRect is the properties column, used to look up the rows an
// open dropdown draws. A dropdown pushes every row below it down, so an option
// is only ever addressed by label inside this rectangle, never by coordinate.
func zoneEditorSidePanelRect() image.Rectangle {
	return image.Rect(
		zoneEditorSidePanelLeft, zoneEditorPanelTop,
		zoneEditorSidePanelRight, zoneEditorPanelTop+zoneEditorPanelHeight)
}

// topologyOptionsRect is the block the open topology dropdown draws its rows in.
func topologyOptionsRect() image.Rectangle {
	return image.Rect(
		topologyOptionsLeft, topologyOptionsTop, topologyOptionsRight, topologyOptionsBottom)
}

// fileDialogRect is the modal panel the file explorer is drawn into.
func fileDialogRect() image.Rectangle {
	return image.Rect(
		fileDialogPanelLeft, fileDialogPanelTop,
		fileDialogPanelLeft+fileDialogPanelWidth, fileDialogPanelTop+fileDialogPanelHeight)
}

// fileStatusMask covers the toolbar's current-file path, which is per-run.
func fileStatusMask() image.Rectangle {
	return image.Rect(fileStatusLeft, 0, WindowWidth, fileStatusBottom)
}

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
