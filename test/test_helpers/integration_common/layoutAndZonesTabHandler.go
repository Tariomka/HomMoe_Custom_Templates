//go:build integration_test

package integration_common

import "gioui.org/f32"

// randomTopologyOptionLabel is the one topology whose preview is regenerated
// with fresh randomness on every run, so it is the one SelectTopology must not
// lift the preview canvas mask for.
const randomTopologyOptionLabel = "Random"

// LayoutAndZonesTabHandler drives the Layout & Zones tab. It embeds the pointer,
// not a copy, so the layout-shift state it records stays visible to every other
// handler.
type LayoutAndZonesTabHandler struct {
	*BaseHandler
}

// SelectTopology opens the topology dropdown and picks the row labelled name.
// Every row is a labelled button, so the option is addressed by its label rather
// than by a row index that a change to the option list would invalidate.
//
// Picking anything but Random makes the generated layout deterministic, so this
// also lifts the preview canvas mask NewHandler registered: from here on the
// preview is compared against its golden like the rest of the window.
func (this *LayoutAndZonesTabHandler) SelectTopology(name string) *LayoutAndZonesTabHandler {
	this.runner.tb.Helper()
	// The tab click that got here is applied during the layout it is polled on,
	// so this panel has not been drawn yet and its trigger has no input area.
	this.runner.NextFrame()
	this.runner.ClickAt(f32.Pt(topologySelectorTriggerX, topologySelectorTriggerY))
	this.runner.ClickButtonIn(topologyOptionsRect(), name)
	this.commit()
	if this.isRandomTopology && name != randomTopologyOptionLabel {
		this.isRandomTopology = false
		this.runner.UnmaskRect(previewCanvasMask())
	}

	this.runner.VerifySnapshot()
	return this
}

func (this *LayoutAndZonesTabHandler) OpenZoneEditor() *ZoneEditorHandler {
	this.runner.tb.Helper()
	this.runner.ClickAt(f32.Pt(zoneEditorButtonX, zoneEditorButtonY))
	// The dialog is pushed during the layout the button is polled on, so it is
	// only drawn - and its canvas only measured - on the frame after that.
	this.runner.NextFrame()
	return newZoneEditorHandler(this.BaseHandler)
}

// ToggleAdvancedZoneControl splits the neutral zone counts into four tier rows.
// It is what gives this panel a scroll range worth testing: it overflows by only
// about 18px otherwise, against roughly 386px with the tiers shown.
func (this *LayoutAndZonesTabHandler) ToggleAdvancedZoneControl() *LayoutAndZonesTabHandler {
	this.runner.tb.Helper()
	this.runner.ClickAt(f32.Pt(advancedZoneControlCheckboxX, advancedZoneControlCheckboxY))
	this.commit()
	this.runner.VerifySnapshot()
	return this
}
