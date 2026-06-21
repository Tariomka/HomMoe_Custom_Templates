package themes

import "image/color"

// Categorical accents — colors that encode a category rather than theme
// chrome. They are tuned to sit comfortably on Crimson Night surfaces.

// Bonus / ban category dots
var (
	ColorDotMovement = color.NRGBA{R: 0x6E, G: 0x9C, B: 0xF0, A: 0xFF} // movement — cornflower blue
	ColorDotCombat   = color.NRGBA{R: 0xE8, G: 0x6A, B: 0x6A, A: 0xFF} // combat — soft red
	ColorDotMagic    = color.NRGBA{R: 0xA0, G: 0x7C, B: 0xE8, A: 0xFF} // magic — violet
	ColorDotSet      = color.NRGBA{R: 0xCC, G: 0x66, B: 0xDE, A: 0xFF} // item sets — orchid
	ColorDotResource = color.NRGBA{R: 0xE0, G: 0xAC, B: 0x4E, A: 0xFF} // resources — amber
	ColorDotDefault  = color.NRGBA{R: 0x8A, G: 0x86, B: 0x8E, A: 0xFF} // fallback — neutral grey
)

// Spell-school accents
var (
	ColorSchoolHighNeutral = color.NRGBA{R: 0xA8, G: 0xA4, B: 0xAC, A: 0xFF}
	ColorSchoolDaylight    = color.NRGBA{R: 0xD6, G: 0xC2, B: 0x7E, A: 0xFF}
	ColorSchoolNightshade  = color.NRGBA{R: 0x9B, G: 0x7E, B: 0xD4, A: 0xFF}
	ColorSchoolArcane      = color.NRGBA{R: 0x7E, G: 0xC9, B: 0xD0, A: 0xFF}
	ColorSchoolPrimal      = color.NRGBA{R: 0xD4, G: 0x7A, B: 0x48, A: 0xFF}
)
