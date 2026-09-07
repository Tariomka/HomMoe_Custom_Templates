//go:build integration_test

package integration_common

// BonusesAndBansTabHandler drives the Bonuses & Bans tab. It embeds the pointer,
// not a copy, so the layout-shift state it records stays visible to every other
// handler.
type BonusesAndBansTabHandler struct {
	*BaseHandler
}
