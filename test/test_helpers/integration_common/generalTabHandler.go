//go:build integration_test

package integration_common

import "gioui.org/f32"

// GeneralTabHandler drives the General tab. It embeds the pointer, not a copy,
// so the layout-shift state it records stays visible to every other handler.
type GeneralTabHandler struct {
	*BaseHandler
}

// The map size dropdown shows this many rows before and after the experimental
// sizes are allowed; see internal/common/mapSizes.go.
const (
	baseMapSizeOptionCount         = 11
	experimentalMapSizeOptionCount = 28
)

// ToggleExperimentalMapSizes flips "Allow non official larger map sizes (>240)"
// and records the resulting option count, which SelectMapSize addresses against.
func (this *GeneralTabHandler) ToggleExperimentalMapSizes() *GeneralTabHandler {
	this.runner.tb.Helper()
	this.runner.ClickAt(f32.Pt(experimentalMapSizesCheckboxX, experimentalMapSizesCheckboxY))
	this.isExperimentalMapSizes = !this.isExperimentalMapSizes
	this.commit()
	this.runner.VerifySnapshot()
	return this
}

// OpenMapSizeSelector drops the map size list open. The list renders inline, so
// this is the only state in which the experimental checkbox's layout shift is
// visible.
func (this *GeneralTabHandler) OpenMapSizeSelector() *GeneralTabHandler {
	this.runner.tb.Helper()
	this.openMapSizeSelector()
	this.runner.VerifySnapshot()
	return this
}

// SelectMapSize opens the map size dropdown and picks the row at optionIndex.
// The dropdown renders inline, so which rows exist depends on what the
// experimental checkbox is currently allowing.
func (this *GeneralTabHandler) SelectMapSize(optionIndex int) *GeneralTabHandler {
	this.runner.tb.Helper()
	if optionIndex < 0 || optionIndex >= this.mapSizeOptionCount() {
		this.runner.tb.Fatalf(
			"map size option %d is out of range: the dropdown shows %d rows (experimental sizes allowed: %v)",
			optionIndex, this.mapSizeOptionCount(), this.isExperimentalMapSizes)
	}

	this.openMapSizeSelector()
	this.runner.ClickAt(f32.Pt(
		mapSizeOptionX,
		mapSizeOptionFirstCenterY+float32(optionIndex*mapSizeOptionPitch)))
	this.commit()
	return this
}

func (this *GeneralTabHandler) openMapSizeSelector() {
	this.runner.tb.Helper()
	this.runner.ClickAt(f32.Pt(mapSizeSelectorTriggerX, mapSizeSelectorTriggerY))
}

// SelectGameMode picks Classic or SingleHero. SingleHero removes three slider
// rows, but only from the right column, and no handler coordinate sits below
// them - so unlike the map size checkbox this shift needs no bookkeeping.
func (this *GeneralTabHandler) SelectGameMode(singleHero bool) *GeneralTabHandler {
	this.runner.tb.Helper()
	buttonX := float32(gameModeClassicX)
	if singleHero {
		buttonX = gameModeSingleHeroX
	}
	this.runner.ClickAt(f32.Pt(buttonX, gameModeCenterY))
	this.commit()
	return this
}

func (this *GeneralTabHandler) mapSizeOptionCount() int {
	if this.isExperimentalMapSizes {
		return experimentalMapSizeOptionCount
	}
	return baseMapSizeOptionCount
}
