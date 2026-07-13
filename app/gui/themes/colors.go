//nolint:gochecknoglobals // centralized color pallet for UI
package themes

import "image/color"

// ColorsBase - Crimson Night themes colors.
var ColorsBase = struct {
	Background color.NRGBA // #161418 window background
	Panel      color.NRGBA // #1E1B20 panels, cards, dialogs
	Input      color.NRGBA // #28242A text inputs, dropdown bodies, wells
	Button     color.NRGBA // #2D2A31 standard button fill
	Border     color.NRGBA // #4A323A default border / outline

	Accent       color.NRGBA // #D9304C headers, focus, active borders
	AccentDim    color.NRGBA // #8E2237 de-emphasised accent
	AccentBright color.NRGBA // #FF5E76 active labels, hot highlights

	Text    color.NRGBA // #EEEAEE primary text
	TextDim color.NRGBA // #9E96A2 secondary text

	Error          color.NRGBA // #FF964B errors
	WarnBackground color.NRGBA // #2A211A warning fills
	WarnBorder     color.NRGBA // #6A5024 warning borders
	WarnText       color.NRGBA // #E0AC4E warnings

	PrimaryButton  color.NRGBA // #771A2C generate / confirm / big bright button fill
	ButtonDisabled color.NRGBA // #2A2326 disabled fill
	BorderDisabled color.NRGBA // #41373B disabled outline
	Selection      color.NRGBA // #5C2433 checked picker rows, dropdown selection

	Hover color.NRGBA // #8E2237 hovered control border highlight
	Scrim color.NRGBA // modal backdrop dimmer
}{
	Background: color.NRGBA{R: 0x16, G: 0x14, B: 0x18, A: 0xFF},
	Panel:      color.NRGBA{R: 0x1E, G: 0x1B, B: 0x20, A: 0xFF},
	Input:      color.NRGBA{R: 0x28, G: 0x24, B: 0x2A, A: 0xFF},
	Button:     color.NRGBA{R: 0x2D, G: 0x2A, B: 0x31, A: 0xFF},
	Border:     color.NRGBA{R: 0x4A, G: 0x32, B: 0x3A, A: 0xFF},

	Accent:       color.NRGBA{R: 0xD9, G: 0x30, B: 0x4C, A: 0xFF},
	AccentDim:    color.NRGBA{R: 0x8E, G: 0x22, B: 0x37, A: 0xFF},
	AccentBright: color.NRGBA{R: 0xFF, G: 0x5E, B: 0x76, A: 0xFF},

	Text:    color.NRGBA{R: 0xEE, G: 0xEA, B: 0xEE, A: 0xFF},
	TextDim: color.NRGBA{R: 0x9E, G: 0x96, B: 0xA2, A: 0xFF},

	Error:          color.NRGBA{R: 0xFF, G: 0x96, B: 0x4B, A: 0xFF},
	WarnBackground: color.NRGBA{R: 0x2A, G: 0x21, B: 0x1A, A: 0xFF},
	WarnBorder:     color.NRGBA{R: 0x6A, G: 0x50, B: 0x24, A: 0xFF},
	WarnText:       color.NRGBA{R: 0xE0, G: 0xAC, B: 0x4E, A: 0xFF},

	PrimaryButton:  color.NRGBA{R: 0x77, G: 0x1A, B: 0x2C, A: 0xFF},
	ButtonDisabled: color.NRGBA{R: 0x2A, G: 0x23, B: 0x26, A: 0xFF},
	BorderDisabled: color.NRGBA{R: 0x41, G: 0x37, B: 0x3B, A: 0xFF},
	Selection:      color.NRGBA{R: 0x5C, G: 0x24, B: 0x33, A: 0xFF},

	Hover: color.NRGBA{R: 0x8E, G: 0x22, B: 0x37, A: 0xFF},
	Scrim: color.NRGBA{A: 0xB0},
}

// ColorsPreview - Map-preview palette.
var ColorsPreview = struct {
	Background color.NRGBA // #19161A canvas
	Frame      color.NRGBA // #9C2E42 crimson frame

	BronzeFill color.NRGBA
	BronzeEdge color.NRGBA
	SilverFill color.NRGBA
	SilverEdge color.NRGBA
	GoldFill   color.NRGBA
	GoldEdge   color.NRGBA
	SpawnFill  color.NRGBA
	SpawnEdge  color.NRGBA
	HubFill    color.NRGBA
	HubEdge    color.NRGBA

	DirectLine color.NRGBA // crimson roads
	PortalLine color.NRGBA // portals stay cool blue

	ZoneLabel   color.NRGBA // zone value text
	CastleBadge color.NRGBA // castle-count badge
}{
	Background:  color.NRGBA{R: 0x19, G: 0x16, B: 0x1A, A: 0xFF},
	Frame:       color.NRGBA{R: 0x9C, G: 0x2E, B: 0x42, A: 0xFF},
	BronzeFill:  color.NRGBA{R: 0x5E, G: 0x40, B: 0x26, A: 0xFF},
	BronzeEdge:  color.NRGBA{R: 0xCD, G: 0x7F, B: 0x32, A: 0xFF},
	SilverFill:  color.NRGBA{R: 0x46, G: 0x4A, B: 0x52, A: 0xFF},
	SilverEdge:  color.NRGBA{R: 0xC4, G: 0xC6, B: 0xCC, A: 0xFF},
	GoldFill:    color.NRGBA{R: 0x75, G: 0x58, B: 0x16, A: 0xFF},
	GoldEdge:    color.NRGBA{R: 0xFF, G: 0xD2, B: 0x32, A: 0xFF},
	SpawnFill:   color.NRGBA{R: 0x28, G: 0x58, B: 0x34, A: 0xFF},
	SpawnEdge:   color.NRGBA{R: 0x64, G: 0xC8, B: 0x78, A: 0xFF},
	HubFill:     color.NRGBA{R: 0x35, G: 0x4E, B: 0x60, A: 0xFF},
	HubEdge:     color.NRGBA{R: 0x82, G: 0xB4, B: 0xC8, A: 0xFF},
	DirectLine:  color.NRGBA{R: 0xC4, G: 0x44, B: 0x58, A: 0xFF},
	PortalLine:  color.NRGBA{R: 0x5A, G: 0xAA, B: 0xD2, A: 0xB4},
	ZoneLabel:   color.NRGBA{R: 0xF2, G: 0xEE, B: 0xF2, A: 0xFF},
	CastleBadge: color.NRGBA{R: 0xFF, G: 0xD9, B: 0x9C, A: 0xFF},
}

// ColorsZoneEditor - Zone-editor canvas accents.
var ColorsZoneEditor = struct {
	EdgeSelected color.NRGBA
	UserAddedDot color.NRGBA
	GuardLabel   color.NRGBA
	GridLine     color.NRGBA
	SnapGuide    color.NRGBA
}{
	EdgeSelected: color.NRGBA{R: 0xFF, G: 0x9E, B: 0x2E, A: 0xFF}, // selected connection
	UserAddedDot: color.NRGBA{R: 0xFF, G: 0xB3, B: 0xC0, A: 0xFF}, // marker on user-added zones
	GuardLabel:   color.NRGBA{R: 0xEF, G: 0xE8, B: 0xD6, A: 0xFF}, // guard-value text on edges
	GridLine:     color.NRGBA{R: 0xF2, G: 0xEE, B: 0xF2, A: 0x0E},
	SnapGuide:    color.NRGBA{R: 0x50, G: 0xDC, B: 0x78, A: 0xB4}, // green alignment guide
}

// ColorsDotCategories - Bonus / ban category dots.
var ColorsDotCategories = struct {
	Movement color.NRGBA
	Combat   color.NRGBA
	Magic    color.NRGBA
	Set      color.NRGBA
	Resource color.NRGBA
	Default  color.NRGBA
}{
	Movement: color.NRGBA{R: 0x6E, G: 0x9C, B: 0xF0, A: 0xFF}, // movement - cornflower blue
	Combat:   color.NRGBA{R: 0xE8, G: 0x6A, B: 0x6A, A: 0xFF}, // combat - soft red
	Magic:    color.NRGBA{R: 0xA0, G: 0x7C, B: 0xE8, A: 0xFF}, // magic - violet
	Set:      color.NRGBA{R: 0xCC, G: 0x66, B: 0xDE, A: 0xFF}, // item sets - orchid
	Resource: color.NRGBA{R: 0xE0, G: 0xAC, B: 0x4E, A: 0xFF}, // resources - amber
	Default:  color.NRGBA{R: 0x8A, G: 0x86, B: 0x8E, A: 0xFF}, // fallback - neutral grey
}

// ColorsSpellSchools - Spell-school accents.
var ColorsSpellSchools = struct {
	HighNeutral color.NRGBA
	Daylight    color.NRGBA
	Nightshade  color.NRGBA
	Arcane      color.NRGBA
	Primal      color.NRGBA
}{
	HighNeutral: color.NRGBA{R: 0xA8, G: 0xA4, B: 0xAC, A: 0xFF},
	Daylight:    color.NRGBA{R: 0xD6, G: 0xC2, B: 0x7E, A: 0xFF},
	Nightshade:  color.NRGBA{R: 0x9B, G: 0x7E, B: 0xD4, A: 0xFF},
	Arcane:      color.NRGBA{R: 0x7E, G: 0xC9, B: 0xD0, A: 0xFF},
	Primal:      color.NRGBA{R: 0xD4, G: 0x7A, B: 0x48, A: 0xFF},
}
