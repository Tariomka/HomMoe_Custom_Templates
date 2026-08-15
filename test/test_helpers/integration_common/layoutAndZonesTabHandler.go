//go:build integration_test

package integration_common

import "gioui.org/f32"

// LayoutAndZonesTabHandler drives the Layout & Zones tab. It embeds the pointer,
// not a copy, so the layout-shift state it records stays visible to every other
// handler.
type LayoutAndZonesTabHandler struct {
	*BaseHandler
}

func (this *LayoutAndZonesTabHandler) OpenZoneEditor() *ZoneEditorHandler {
	this.runner.tb.Helper()
	this.runner.ClickAt(f32.Pt(zoneEditorButtonX, zoneEditorButtonY))
	return &ZoneEditorHandler{base: this.BaseHandler}
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
