//go:build integration_test

package integration_common

// The labels the Bonuses & Bans tab's buttons are addressed by.
const (
	addBannedItemButtonLabel  = "+ Add banned item..."
	addBannedSpellButtonLabel = "+ Add banned spell..."
	addOverrideButtonLabel    = "+ Add override..."
	removeEntryButtonLabel    = "X"
)

// BonusesAndBansTabHandler drives the Bonuses & Bans tab. It embeds the pointer,
// not a copy, so the layout-shift state it records stays visible to every other
// handler.
type BonusesAndBansTabHandler struct {
	*BaseHandler
}

// BanItem opens the banned item picker, picks the item with the given sid and
// confirms it.
func (this *BonusesAndBansTabHandler) BanItem(sid string) *BonusesAndBansTabHandler {
	this.runner.tb.Helper()
	return this.pickThrough(addBannedItemButtonLabel, sid)
}

// BanSpell opens the banned spell picker, picks the spell with the given sid and
// confirms it.
func (this *BonusesAndBansTabHandler) BanSpell(sid string) *BonusesAndBansTabHandler {
	this.runner.tb.Helper()
	return this.pickThrough(addBannedSpellButtonLabel, sid)
}

// OverrideGuardValue opens the guard value override picker, picks the sid and
// confirms it; the picker gives the entry its default guard value.
func (this *BonusesAndBansTabHandler) OverrideGuardValue(sid string) *BonusesAndBansTabHandler {
	this.runner.tb.Helper()
	return this.pickThrough(addOverrideButtonLabel, sid)
}

// RemoveEntry presses the remove button of the index-th row on the tab, counted
// top to bottom and then left to right over both columns.
func (this *BonusesAndBansTabHandler) RemoveEntry(index int) *BonusesAndBansTabHandler {
	this.runner.tb.Helper()
	this.runner.ClickNthButtonIn(bonusesPanelRect(), removeEntryButtonLabel, index)
	this.commit()
	return this
}

// pickThrough presses the button that opens a picker, picks one entry in it and
// presses Add. A click is consumed on the frame after it is queued, so a frame
// follows each step.
func (this *BonusesAndBansTabHandler) pickThrough(openLabel string, id string) *BonusesAndBansTabHandler {
	this.runner.tb.Helper()
	this.runner.NextFrame()
	this.runner.ClickButtonIn(bonusesPanelRect(), openLabel)
	this.runner.NextFrame()
	this.runner.SelectInTopPicker(id)
	this.runner.NextFrame()
	this.runner.ConfirmTopPicker()
	this.runner.NextFrame()
	this.commit()
	return this
}
